package validation

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/reddot-watch/FeedFetcher/internal/dateparser"
)

var (
	ErrInvalidURL                = errors.New("invalid url")
	ErrEmptyHeadline             = errors.New("empty headline")
	ErrFeedPublicationDateFormat = errors.New("publication date format is invalid")
	ErrHeadlineTooLong           = errors.New("headline exceeds maximum length")
	ErrPublicationTooOld         = errors.New("publication date exceeds maximum age")
	ErrFuturePublication         = errors.New("publication date is in the future beyond allowed tolerance")
	ErrMissingPublishDate        = errors.New("missing publication date")
)

// ValidateAndResolveURL validates and resolves a relative url against the feed url.
// Extracted as a package function for better testability.
func ValidateAndResolveURL(feedURL *url.URL, rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", ErrInvalidURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrInvalidURL, err)
	}

	return feedURL.ResolveReference(parsed).String(), nil
}

var spaceRegexp = regexp.MustCompile(`\s+`)

// ValidateAndSanitizeHeadline cleans and validates the item headline.
// Extracted as a package function for better testability.
func ValidateAndSanitizeHeadline(rawHeadline string, maxLength int) (string, error) {
	if rawHeadline == "" {
		return "", ErrEmptyHeadline
	}

	// Normalize whitespace and trim
	headline := spaceRegexp.ReplaceAllString(rawHeadline, " ")
	headline = strings.TrimSpace(headline)

	// Check length (using rune count for proper Unicode handling)
	if len([]rune(headline)) > maxLength {
		return "", ErrHeadlineTooLong
	}

	return headline, nil
}

// ValidatePublicationDate validates that the item's publication date is within acceptable bounds.
// Extracted as a package function for better testability.
func ValidatePublicationDate(item *gofeed.Item, maxAge, futureTolerance time.Duration) (time.Time, error) {
	if item.PublishedParsed == nil {
		if pubDate := item.Published; pubDate == "" {
			return time.Time{}, ErrMissingPublishDate
		} else if t, err := dateparser.ParseDateWithDefaultTZ(pubDate); err == nil {
			item.PublishedParsed = &t
		} else {
			return time.Time{}, fmt.Errorf("%w: %s", ErrFeedPublicationDateFormat, err)
		}
	}

	now := time.Now().UTC()
	pubDate := item.PublishedParsed.UTC()

	// Check if the publication date is too far in the future
	if pubDate.Sub(now) > futureTolerance {
		return time.Time{}, ErrFuturePublication
	}

	// Check if the publication date is too old
	if now.Sub(pubDate) > maxAge {
		return time.Time{}, ErrPublicationTooOld
	}

	return pubDate, nil
}

// ExtractContent gets the best available content from an item.
// Extracted as a package function for better testability.
func ExtractContent(item *gofeed.Item) string {
	// First try description, then fallback to content
	if content := strings.TrimSpace(item.Description); content != "" {
		return content
	}
	return strings.TrimSpace(item.Content)
}
