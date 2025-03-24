# FeedFetcher

FeedFetcher is a simple Go wrapper around the [gofeed](https://github.com/mmcdole/gofeed) library that adds some useful features for handling RSS/Atom feeds.

## Overview

This package uses [gofeed](https://github.com/mmcdole/gofeed) for the core feed parsing functionality and adds some practical utilities to help with common feed processing tasks.

## Additional Features

- Domain-based rate limiting to manage request frequency
- Basic content validation and sanitization
- Publication date validation
- Headline length enforcement
- Integration with zerolog for logging

## Installation

```bash
go get github.com/reddot-watch/feedfetcher
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/reddot-watch/feedfetcher"
)

func main() {
	fetcher := FeedFetcher.NewFeedFetcher(FeedFetcher.Config{
		UserAgent:            "MyFeedReader/1.0",
		RequestTimeout:       15 * time.Second,
		MaxItems:             50,
		MaxHeadingLength:     200,
		MaxAge:               48 * time.Hour,
		FutureDriftTolerance: 12 * time.Hour,
	})
	
	items, err := fetcher.FetchAndProcess(context.Background(), "https://example.com/feed.xml")
	if err != nil {
		log.Fatalf("Error fetching feed: %v", err)
	}

	for _, item := range items {
		fmt.Printf("Headline: %s\nURL: %s\nPublished: %s\n\n",
			item.Headline, item.URL, item.PublishedAt)
	}
}
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| UserAgent | HTTP User-Agent header | "Mozilla/5.0 (compatible; ReddotWatchBot/1.0; +https://reddot.watch/bot)" |
| RequestTimeout | Timeout for feed requests | 10 seconds |
| MaxItems | Maximum number of items to process (0 = unlimited) | 1000 |
| MaxHeadingLength | Maximum allowed headline length | 250 characters |
| MaxAge | Maximum age of feed items to consider valid | 24 hours |
| FutureDriftTolerance | Tolerance for items with future timestamps | 24 hours |

## Method Chaining

FeedFetcher supports method chaining for configuration:

```go
fetcher := FeedFetcher.NewDefaultFeedFetcher().
    WithUserAgent("MyFeedReader/1.0").
    WithMaxItems(50).
    WithRequestTimeout(15 * time.Second)
```

## Use Cases

- When you need feed parsing with rate limiting
- When you want basic validation of feed content
- When you need to manage time drift age restrictions
- As a simple extension to gofeed with some helpful utilities

## License

This project is licensed under Apache 2.0 - See [LICENSE](LICENSE)