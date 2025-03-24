package feedfetcher

import (
	"context"
	"errors"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/reddot-watch/feedfetcher/internal/feedparser"
)

type MockFeedParser struct {
	MockParseFn func(url string, ctx context.Context) (*gofeed.Feed, error)
}

func (m *MockFeedParser) ParseURLWithContext(url string, ctx context.Context) (*gofeed.Feed, error) {
	return m.MockParseFn(url, ctx)
}

// NewFeedFetcherWithParser allows injecting a custom feedparser implementation
func NewFeedFetcherWithParser(config Config, parser feedparser.Parser) *FeedFetcher {
	return &FeedFetcher{
		config: config,
		parser: parser,
	}
}

func TestFeedFetcher_FetchFeed(t *testing.T) {
	// Create a mock feed
	mockFeed := &gofeed.Feed{
		Title: "Test feed",
		Items: []*gofeed.Item{
			{
				Title:           "Test Item",
				Link:            "https://example.com/item",
				Description:     "Test description",
				Published:       "Mon, 02 Jan 2023 15:04:05 GMT",
				PublishedParsed: timePtr(time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC)),
			},
		},
	}

	// Create mock feedparser
	mockParser := &MockFeedParser{
		MockParseFn: func(url string, ctx context.Context) (*gofeed.Feed, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			if url == "https://example.com/valid" {
				return mockFeed, nil
			}
			if url == "https://example.com/timeout" {
				return nil, context.DeadlineExceeded
			}
			return nil, errors.New("invalid url")
		},
	}

	// Create fetcher with mock feedparser
	config := DefaultConfig
	fetcher := NewFeedFetcherWithParser(config, mockParser)

	// Test successful fetch
	t.Run("successful fetch", func(t *testing.T) {
		ctx := context.Background()
		f := &feed{url: "https://example.com/valid"}
		err := fetcher.download(ctx, f)
		assert.NoError(t, err)
		assert.Equal(t, mockFeed, f.data)
	})

	// Test timeout
	t.Run("timeout", func(t *testing.T) {
		ctx := context.Background()
		f := &feed{url: "https://example.com/timeout"}
		err := fetcher.download(ctx, f)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.DeadlineExceeded))
	})

	// Test canceled context
	t.Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		f := &feed{url: "https://example.com/valid"}
		err := fetcher.download(ctx, f)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled))
	})
}

// Helper function to create a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
