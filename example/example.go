package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/reddot-watch/feedfetcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run example.go <feed_url> [max_items]")
		fmt.Println("Example: go run example.go https://news.ycombinator.com/rss 10")
		os.Exit(1)
	}

	feedURL := os.Args[1]
	maxItems := 10 // Default value
	if len(os.Args) > 2 {
		_, err := fmt.Sscanf(os.Args[2], "%d", &maxItems)
		if err != nil {
			log.Fatal().Err(err).Msg("Invalid max_items parameter")
		}
	}

	// Create a customized feed fetcher
	fetcher := FeedFetcher.NewDefaultFeedFetcher().
		WithMaxItems(maxItems).
		WithMaxAge(7 * 24 * time.Hour). // Accept items up to a week old
		WithRequestTimeout(15 * time.Second)

	// Add custom logger
	fetcher = fetcher.WithLogger(log.With().Str("component", "example_app").Logger())

	log.Info().
		Str("url", feedURL).
		Int("max_items", maxItems).
		Dur("max_age", 7*24*time.Hour).
		Msg("Fetching feed")

	// Fetch and process the feed
	ctx := context.Background()
	items, err := fetcher.FetchAndProcess(ctx, feedURL)
	if err != nil {
		log.Fatal().Err(err).Str("url", feedURL).Msg("Failed to fetch feed")
	}

	log.Info().Int("count", len(items)).Msg("Feed items retrieved")

	// Display basic information about each item
	for i, item := range items {
		log.Info().
			Int("index", i+1).
			Str("headline", item.Headline).
			Str("url", item.URL).
			Time("published", item.PublishedAt).
			Msg("Feed item")
	}

	// Output the first item in JSON format if available
	if len(items) > 0 {
		jsonData, err := json.MarshalIndent(items[0], "", "  ")
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal item to JSON")
		} else {
			fmt.Println("\nFirst item as JSON:")
			fmt.Println(string(jsonData))
		}
	}
}
