package crawler

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/meowalien/RabbitGather-proto/go/interest"
	"github.com/meowalien/RabbitGather-proto/go/share"
	"github.com/meowalien/go-meowalien-lib/grpcs"
)

type A struct {
	Url           string
	QuerySelector string
}

func TestRPCService(t *testing.T) {
	grpcConnect := grpcs.GRPCClientConstructor{
		ClientParameters: nil,
		ServerHost:       "127.0.0.1:50051",
	}.New()
	client := interest.NewInterestCrawlerClient(grpcConnect)

	bff := bytes.Buffer{}

	err := gob.NewEncoder(&bff).Encode(A{
		Url:           "http://go-colly.org/",
		QuerySelector: "body > div.content > div.ui.inverted.vertical.masthead.center.aligned.segment > div > h2",
	})
	if err != nil {
		return
	}

	resp, err := client.Crawl(context.Background(), &interest.CrawlRequest{
		Type: interest.CrawlerType_SimpleCrawler,
		Message: &share.EncodedMessage{
			Encoding: share.Encoding_GOB,
			Data:     bff.Bytes(),
		},
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(resp)
}
