package Videourl

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/pwh-pwh/aiwechat-vercel/db"
)

type Video struct {
	Link string
	Text string
}

func VideoConvert(videoname string) {

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
			Url = `^(100|[1-9][0-9]?|)$`
			c.Visit(e.Request.AbsoluteURL(link))
		}

	})

	c.Visit("https://so.iqiyi.com/so/q_" + videoname)

	db.ChatDbInstance.SetVideoValue(videoname, data)
}
