package crawler

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"

	"github.com/meowalien/RabbitGather-proto/go/interest"
	"github.com/meowalien/RabbitGather-proto/go/share"
	"github.com/meowalien/go-meowalien-lib/errs"
)

type Constructor struct {
	Type              interest.CrawlerType
	GobEncodingStruct []byte
}

// NewCrawler creates a new crawler from the given type and the gob encoded struct.
func NewCrawler(t interest.CrawlerType, msg *share.EncodedMessage) (cw Crawler) {
	switch t {
	case interest.CrawlerType_SimpleCrawler:
		var setting SimpleCrawlerConstructor
		err := decodeInto(msg, &setting)
		if err != nil {
			panic(err)
		}
		return setting.New()
	default:
		panic(ErrUnknownCrawlerType)
		return
	}
}

// Crawler is the interface that abstracts the crawler.
type Crawler interface {
	// Crawl starts the crawler, the result or error will be sent to the channel.
	// The result channel will be closed when the crawler is done.
	// The crawler will be terminated when the context is done.
	Crawl(ctx context.Context) (result <-chan []byte, errChannel <-chan error)
}

var ErrUnknownCrawlerType = errors.New("unknown crawler type")
var ErrUnknownEncoding = errors.New("unknown encoding")

func decodeInto[T any](msg *share.EncodedMessage, s *T) (err error) {
	switch msg.Encoding {
	case share.Encoding_GOB:
		err = gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(s)
		if err != nil {
			err = errs.New(err)
			return
		}
	default:
		err = errs.New(ErrUnknownEncoding)
	}

	return
}
