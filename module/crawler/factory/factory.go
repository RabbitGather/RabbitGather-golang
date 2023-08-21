package factory

import (
	"errors"

	"github.com/meowalien/RabbitGather-proto/go/interest"
	"github.com/meowalien/RabbitGather-proto/go/share"

	"github.com/meowalien/RabbitGather-interest-crawler.git/lib"
	"github.com/meowalien/RabbitGather-interest-crawler.git/module/crawler"
	"github.com/meowalien/RabbitGather-interest-crawler.git/module/crawler/simple_crawler"
)

type Factory interface {
	NewCrawler(t interest.CrawlerType, msg *share.EncodedMessage) (cw crawler.Crawler, err error)
}
type Constructor struct {
}

func (c Constructor) New() Factory {
	return &crawlerFactory{}
}

type crawlerFactory struct {
}

// NewCrawler creates a new crawler from the given type and the gob encoded struct.
func (c *crawlerFactory) NewCrawler(t interest.CrawlerType, msg *share.EncodedMessage) (cw crawler.Crawler, err error) {
	switch t {
	case interest.CrawlerType_SimpleCrawler:
		var setting simple_crawler.SimpleCrawlerConstructor
		err = lib.DecodeMessage(msg.Encoding, msg.Data, &setting)
		if err != nil {
			return nil, err
		}
		return setting.New(), nil
	default:
		return nil, ErrUnknownCrawlerType
	}
}

var ErrUnknownCrawlerType = errors.New("unknown crawler type")
