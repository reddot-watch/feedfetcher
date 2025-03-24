package limiter

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter interface for domain-based rate limiting
type RateLimiter interface {
	Wait(ctx context.Context) error
	Allow() bool
}

// DomainRateLimiter limits requests by domain
type DomainRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

// NewDomainRateLimiter creates a rate limiter that limits by domain
// r is requests per second, b is burst size
func NewDomainRateLimiter(r rate.Limit, b int) *DomainRateLimiter {
	return &DomainRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
		mu:       sync.RWMutex{},
	}
}

// getLimiter gets or creates a limiter for a domain
func (l *DomainRateLimiter) getLimiter(domain string) *rate.Limiter {
	l.mu.RLock()
	limiter, exists := l.limiters[domain]
	l.mu.RUnlock()

	if !exists {
		l.mu.Lock()
		// Double-check to avoid race conditions
		if limiter, exists = l.limiters[domain]; !exists {
			limiter = rate.NewLimiter(l.r, l.b)
			l.limiters[domain] = limiter
		}
		l.mu.Unlock()
	}

	return limiter
}

// WaitForDomain waits until a request is allowed for the domain
func (l *DomainRateLimiter) WaitForDomain(ctx context.Context, urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("failed to parse URL for rate limiting: %w", err)
	}

	// Normalize domain by trimming any potential "www." prefix
	host := strings.TrimPrefix(u.Host, "www.")
	if host == "" {
		return fmt.Errorf("empty host in URL: %s", urlStr)
	}

	return l.getLimiter(host).Wait(ctx)
}
