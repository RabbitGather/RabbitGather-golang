package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.CacheDir("./coursera_cache"),
		colly.Async(true),
	)
	c.OnHTML("body > div.content > div.ui.inverted.vertical.masthead.center.aligned.segment > div > h2", func(e *colly.HTMLElement) {
		fmt.Println("Text: ", e.Text)
	})
	// Find and visit all links
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	// print the element "body > div.content > div.ui.vertical.stripe.segment.padding > div > div > div:nth-child(1) > ul > li:nth-child(3)"
	//
	//	e.Request.Visit(e.Attr("href"))
	//})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}
