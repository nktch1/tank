package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseSearchResults(ctx context.Context, response []byte) *SearchResults {
	res := &SearchResults{
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

func getRootDomain(domain string) string {
	domain = strings.ToLower(domain)

	parts := strings.Split(domain, ".")
	if len(parts) < 3 {
		return domain
	}

	if _, ok := tlds[strings.Join(parts[len(parts)-2:], ".")]; ok {
		return strings.Join(parts[len(parts)-3:], ".")
	}

	return strings.Join(parts[len(parts)-2:], ".")
}