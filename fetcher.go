package FeedFetcher

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"golang.org/x/time/rate"

	"github.com/mmcdole/gofeed"
	"github.com/reddot-watch/FeedFetcher/internal/feedparser"
	"github.com/reddot-watch/FeedFetcher/internal/limiter"
	"github.com/reddot-watch/FeedFetcher/internal/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var DefaultConfig = Config{
	UserAgent:            "Mozilla/5.0 (compatible; ReddotWatchBot/1.0; +https://reddot.watch/bot)",
	RequestTimeout:       10 * time.Second,
	MaxItems:             1000,
	MaxHeadingLength:     250,
	MaxAge:               24 * time.Hour,
	FutureDriftTolerance: 24 * time.Hour,
}

// Config holds the configuration for the feed fetcher.
type Config struct {
	UserAgent            string
	RequestTimeout       time.Duration
	MaxItems             int // Use 0 or negative value for no limit
	MaxHeadingLength     int
	MaxAge               time.Duration
	FutureDriftTolerance time.Duration
}

// FeedItem represents a single item from a feed.
type FeedItem struct {
	ID          int64
	FeedURL     string
	URL         string
	Headline    string
	Content     string
	PublishedAt time.Time
}

// FeedFetcher handles retrieving and processing feed data.
type FeedFetcher struct {
	config      Config
	parser      feedparser.Parser
	rateLimiter *limiter.DomainRateLimiter
	logger      zerolog.Logger
}

// NewFeedFetcher creates a new FeedFetcher with the provided configuration.
func NewFeedFetcher(config Config) *FeedFetcher {
	parser := feedparser.NewGoFeedParser(config.UserAgent)

	// Default limit 1 req/sec per domain with burst of 3
	rateLimiter := limiter.NewDomainRateLimiter(rate.Limit(1), 3)

	return &FeedFetcher{
		config:      config,
		parser:      parser,
		rateLimiter: rateLimiter,
		logger:      log.With().Str("component", "feed_fetcher").Logger(),
	}
}

// NewDefaultFeedFetcher creates a new FeedFetcher with default configuration.
func NewDefaultFeedFetcher() *FeedFetcher {
	return NewFeedFetcher(DefaultConfig)
}

// WithLogger returns a new FeedFetcher with a custom logger
func (f *FeedFetcher) WithLogger(logger zerolog.Logger) *FeedFetcher {
	newFetcher := *f
	newFetcher.logger = logger
	return &newFetcher
}

// WithMaxItems returns a new FeedFetcher with an updated MaxItems setting.
// Use maxItems <= 0 to fetch all available items.
func (f *FeedFetcher) WithMaxItems(maxItems int) *FeedFetcher {
	newFetcher := *f
	newConfig := f.config
	newConfig.MaxItems = maxItems
	newFetcher.config = newConfig
	return &newFetcher
}

// WithMaxAge returns a new FeedFetcher with an updated MaxAge setting.
func (f *FeedFetcher) WithMaxAge(duration time.Duration) *FeedFetcher {
	newFetcher := *f
	newConfig := f.config
	newConfig.MaxAge = duration
	newFetcher.config = newConfig
	return &newFetcher
}

// WithUserAgent returns a new FeedFetcher with an updated UserAgent setting.
func (f *FeedFetcher) WithUserAgent(userAgent string) *FeedFetcher {
	newFetcher := *f
	newConfig := f.config
	newConfig.UserAgent = userAgent
	newFetcher.config = newConfig

	// Special case: also need to update the parser
	newParser := feedparser.NewGoFeedParser(userAgent)
	newFetcher.parser = newParser

	return &newFetcher
}

// WithRequestTimeout returns a new FeedFetcher with an updated RequestTimeout setting.
func (f *FeedFetcher) WithRequestTimeout(timeout time.Duration) *FeedFetcher {
	newFetcher := *f
	newConfig := f.config
	newConfig.RequestTimeout = timeout
	newFetcher.config = newConfig
	return &newFetcher
}

// WithMaxHeadingLength returns a new FeedFetcher with an updated MaxHeadingLength setting.
func (f *FeedFetcher) WithMaxHeadingLength(length int) *FeedFetcher {
	newFetcher := *f
	newConfig := f.config
	newConfig.MaxHeadingLength = length
	newFetcher.config = newConfig
	return &newFetcher
}

// WithFutureDriftTolerance returns a new FeedFetcher with an updated FutureDriftTolerance setting.
func (f *FeedFetcher) WithFutureDriftTolerance(duration time.Duration) *FeedFetcher {
	newFetcher := *f
	newConfig := f.config
	newConfig.FutureDriftTolerance = duration
	newFetcher.config = newConfig
	return &newFetcher
}

// FetchAndProcess fetches and processes a feed, returning the parsed items.
func (f *FeedFetcher) FetchAndProcess(ctx context.Context, feedURL string) ([]*FeedItem, error) {
	if err := f.rateLimiter.WaitForDomain(ctx, feedURL); err != nil {
		return nil, err
	}

	ff, err := f.newFeed(feedURL)
	if err != nil {
		return nil, err
	}

	if err := f.download(ctx, ff); err != nil {
		return nil, err
	}

	return f.extractItems(ff)
}

type feed struct {
	url       string
	parsedURL *url.URL
	data      *gofeed.Feed
}

func (f *FeedFetcher) newFeed(feedURL string) (*feed, error) {
	parsedURL, err := url.Parse(feedURL)
	if err != nil {
		return nil, fmt.Errorf("invalid feed url %s: %w", feedURL, err)
	}

	return &feed{
		url:       feedURL,
		parsedURL: parsedURL,
	}, nil
}

// download retrieves a feed from the given url.
func (f *FeedFetcher) download(ctx context.Context, feed *feed) error {
	ctx, cancel := context.WithTimeout(ctx, f.config.RequestTimeout)
	defer cancel()

	f.logger.Debug().Str("url", feed.url).Msg("downloading feed")
	startTime := time.Now()

	result, err := f.parser.ParseURLWithContext(feed.url, ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			f.logger.Error().Str("url", feed.url).Err(err).Msg("deadline exceeded")
			return fmt.Errorf("timed out after %v fetching feed %s: %w",
				f.config.RequestTimeout, feed.url, err)
		}

		if errors.Is(err, context.Canceled) {
			f.logger.Warn().Str("url", feed.url).Msg("request was canceled")
			return fmt.Errorf("feed fetch canceled for %s: %w", feed.url, err)
		}

		f.logger.Error().Str("url", feed.url).Err(err).Msg("failed to parse feed")
		return fmt.Errorf("failed to parse feed url %s: %w", feed.url, err)
	}

	feed.data = result

	f.logger.Debug().
		Str("url", feed.url).
		Dur("duration", time.Since(startTime)).
		Int("items", len(feed.data.Items)).
		Msg("feed downloaded successfully")

	return nil
}

func (f *FeedFetcher) extractItems(feed *feed) ([]*FeedItem, error) {
	if feed == nil {
		return nil, errors.New("feed cannot be nil")
	}

	// Determine how many items to process
	itemCount := len(feed.data.Items)
	if f.config.MaxItems > 0 && f.config.MaxItems < itemCount {
		itemCount = f.config.MaxItems
	}

	result := make([]*FeedItem, 0, itemCount)

	for i := 0; i < itemCount; i++ {
		item := feed.data.Items[i]
		if item == nil {
			continue
		}

		parsed, err := f.validateAndConvertItem(feed.parsedURL, item)
		if err != nil {
			if errors.Is(err, validation.ErrFeedPublicationDateFormat) {
				// Do not process other items as they will all have the same error
				return nil, validation.ErrFeedPublicationDateFormat
			}
			// Continue processing other items
			continue
		}

		if parsed != nil {
			result = append(result, parsed)
		}
	}

	return result, nil
}

func (f *FeedFetcher) validateAndConvertItem(feedURL *url.URL, item *gofeed.Item) (*FeedItem, error) {
	if feedURL == nil || item == nil {
		return nil, errors.New("feedURL and item cannot be nil")
	}

	itemURL, err := validation.ValidateAndResolveURL(feedURL, item.Link)
	if err != nil {
		return nil, err
	}

	publishedAt, err := validation.ValidatePublicationDate(item, f.config.MaxAge, f.config.FutureDriftTolerance)
	if err != nil {
		return nil, err
	}

	headline, err := validation.ValidateAndSanitizeHeadline(item.Title, f.config.MaxHeadingLength)
	if err != nil {
		return nil, err
	}

	content := validation.ExtractContent(item)

	return &FeedItem{
		FeedURL:     feedURL.String(),
		URL:         itemURL,
		PublishedAt: publishedAt,
		Headline:    headline,
		Content:     content,
	}, nil
}
