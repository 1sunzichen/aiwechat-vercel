package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Video struct {
	Link string
	Text string
}

func VideoConvert() []Video {

	data := []Video{}
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("so.iqiyi.com", "www.iqiyi.com", "iqiyi.com", "v.qq.com"),
	)
	Url := "春色寄情人"
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		match, _ := regexp.MatchString(Url, e.Text)
		fmt.Println("link", match)
		if match && (strings.Contains(link, "www.iqiyi.com")) {

			fmt.Printf("Link found城中之城: %q -> %s\n", e.Text, link)

			if Url == `^(100|[1-9][0-9]?|)$` {

				data = append(data, Video{Link: link, Text: e.Text})

			}
			Url = `^(100|[1-9][0-9]?|)$`
			c.Visit(e.Request.AbsoluteURL(link))
		}

	})

	c.Visit("https://so.iqiyi.com/so/q_" + "春色寄情人")
	newdata := []Video{}
	for i, v := range data {
		it, _ := strconv.Atoi(v.Text)
		if i <= it {
			newdata = append(newdata, v)
		}
	}
	for _, v := range newdata {
		fmt.Println(v)
	}
	return newdata
}
func main() {
	VideoConvert()
}
