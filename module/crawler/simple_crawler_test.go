package crawler

import (
	"context"
	"fmt"
	"testing"

	"github.com/meowalien/RabbitGather-interest-crawler.git/lib"
)

func TestSimpleCrawler(t *testing.T) {
	cw := SimpleCrawlerConstructor{
		Url:           "http://go-colly.org/",
		QuerySelector: "body > div.content > div.ui.inverted.vertical.masthead.center.aligned.segment > div > h2",
	}.New()

	resultChan, errChan := cw.Crawl(context.Background())
	err := lib.Pipe(resultChan, errChan, func(result []byte) bool {
		if result == nil {
			return false
		}
		fmt.Println("result: ", string(result))
		return true
	})
	if err != nil {
		t.Fatal(err)
		return
	}
}
