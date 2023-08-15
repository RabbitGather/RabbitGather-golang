package crawler

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"

	"github.com/meowalien/go-meowalien-lib/errs"
)

type CrawlerType uint32

// Crawler is the interface that abstracts the crawler.
type Crawler interface {
	// Crawl starts the crawler, the result or error will be sent to the channel.
	// The result channel will be closed when the crawler is done.
	// The crawler will be terminated when the context is done.
	Crawl(ctx context.Context) (result <-chan []byte, errChannel <-chan error)
}

var ErrUnknownCrawlerType = errors.New("unknown crawler type")

// NewCrawler creates a new crawler from the given type and the gob encoded struct.
func NewCrawler(t CrawlerType, gobEncodingStruct []byte) (cw Crawler, err error) {
	switch t {
	case SimpleCrawler:
		var setting SimpleCrawlerSetting
		err = decodeInto(gobEncodingStruct, &setting)
		if err != nil {
			err = errs.New(err)
			return
		}
		return setting.New(), nil
	default:
		err = ErrUnknownCrawlerType
		return
	}
}

func decodeInto[T any](encodingStruct []byte, s *T) (err error) {
	err = gob.NewDecoder(bytes.NewReader(encodingStruct)).Decode(s)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
