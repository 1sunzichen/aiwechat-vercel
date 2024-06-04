package Videourl

import (
	"context"
	"crypto/tls"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
)

type Video struct {
	Link string
	Text string
}

func NewClient() *redis.Client {
	options, err := redis.ParseURL("redis://default:5e27d347869141eeb77127ccbe40b5ad@in-maggot-33849.upstash.io:33849")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	client := redis.NewClient(options)
	return client
}
func VideoConvert(videoname string) {
	client := NewClient()
	data := ""
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("so.iqiyi.com", "www.iqiyi.com", "iqiyi.com", "v.qq.com"),
	)
	Url := videoname
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		match, _ := regexp.MatchString(Url, e.Text)
		if match && (strings.Contains(link, "www.iqiyi.com")) {

			data = "https://mj.mailseason.com/vip?url=http:" + link
			client.Set(context.Background(), videoname, data, 0)
			fmt.Println(data, "data")
			Url = `^(100|[1-9][0-9]?|)$`
			c.Visit(e.Request.AbsoluteURL(link))
		}

	})

	c.Visit("https://so.iqiyi.com/so/q_" + videoname)

}
