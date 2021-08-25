package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseYandexURL = "https://yandex.ru/search/touch/" +
	"?service=www.yandex&ui=webmobileapp.yandex&numdoc=50&lr=213&p=0&text=%s"

func ParseSearchResponse(ctx context.Context, response []byte) Response {
	res := Response{
		Items: make([]responseItem, 0),
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response))
	if err != nil {
		res.Error = fmt.Errorf("can't create parser for body: %v", err)
		return res
	}

	items := doc.Find("div.serp-item")

	items.Each(func(i int, selection *goquery.Selection) {
		_, aExists := selection.Attr("data-fast-name")
		_, cidExists := selection.Attr("data-cid")

		if !selection.Is("div.Label") &&
			!aExists && !selection.Is("span.organic__advLabel") && cidExists {

			link := selection.Find("a.Link").First()

			if link == nil {
				return
			}

			urlStr, _ := link.Attr("href")

			dcStr, _ := link.Attr("data-counter")
			if strings.HasPrefix(urlStr, "https://yandex.ru/turbo/") ||
				strings.Contains(urlStr, "turbopages.org") && dcStr != "" {
				var dc []string
				if err = json.Unmarshal([]byte(dcStr), &dc); err != nil || len(dc) < 2 {
					return
				}

				urlStr = dc[1]
			}

			u, err := url.Parse(urlStr)
			if err != nil {
				return
			}

			if u.Host == "" || u.Host == "yabs.yandex.ru" {
				return
			}

			res.Items = append(res.Items, responseItem{
				Host: getRootDomain(u.Host),
				Url:  urlStr,
			})
		}
	})

	return res
}