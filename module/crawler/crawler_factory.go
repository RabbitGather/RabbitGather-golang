package crawler

import (
	"github.com/meowalien/RabbitGather-proto/go/interest"
	"github.com/meowalien/RabbitGather-proto/go/share"
)

type Factory interface {
	NewCrawler(t interest.CrawlerType, msg *share.EncodedMessage) (cw Crawler, err error)
}

type CrawlerFactoryConstructor struct {
}

func (c CrawlerFactoryConstructor) New() Factory {
	return &crawlerFactory{}
}

type crawlerFactory struct {
}

// NewCrawler creates a new crawler from the given type and the gob encoded struct.
func (c *crawlerFactory) NewCrawler(t interest.CrawlerType, msg *share.EncodedMessage) (cw Crawler, err error) {
	switch t {
	case interest.CrawlerType_SimpleCrawler:
		var setting SimpleCrawlerConstructor
		err = decodeInto(msg, &setting)
		if err != nil {
			return nil, err
		}
		return setting.New(), nil
	default:
		return nil, ErrUnknownCrawlerType
	}
}
