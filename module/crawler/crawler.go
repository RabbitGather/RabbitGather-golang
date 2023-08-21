package crawler

import "context"

// Crawler is the interface that abstracts the crawler.
type Crawler interface {
	// Crawl starts the crawler, the result or error will be sent to the channel.
	// The result channel will be closed when the crawler is done.
	// The crawler will be terminated when the context is done.
	Crawl(ctx context.Context) (result <-chan []byte, errChannel <-chan error)
}
