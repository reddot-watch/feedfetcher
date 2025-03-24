package feedparser

import (
	"context"

	"github.com/mmcdole/gofeed"
)

type Parser interface {
	ParseURLWithContext(url string, ctx context.Context) (*gofeed.Feed, error)
}

type GoFeedParser struct {
	parser *gofeed.Parser
}

func NewGoFeedParser(userAgent string) *GoFeedParser {
	parser := gofeed.NewParser()
	parser.UserAgent = userAgent
	return &GoFeedParser{parser: parser}
}

func (p *GoFeedParser) ParseURLWithContext(url string, ctx context.Context) (*gofeed.Feed, error) {
	return p.parser.ParseURLWithContext(url, ctx)
}
