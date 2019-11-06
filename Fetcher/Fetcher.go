package Fetcher

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

var client = &http.Client{}

func FetchRelativeUrlsFromPage(baseURL *string, targetURL string) ([]string, bool) {
	resp, err := httpGet(targetURL)
	var urls []string
	if err == nil {
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		_ = resp.Body.Close()
		urls = extractRelativeUrls(baseURL, doc)
	}
	return urls, err == nil
}

func httpGet(url string) (*http.Response, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func extractRelativeUrls(baseUrl *string, doc *goquery.Document) []string {
	var relativeUrls []string
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			if link, ok := s.Attr("href"); ok {
				if strings.HasPrefix(link, *baseUrl) {
					relativeUrls = append(relativeUrls, strings.TrimSuffix(link, "/"))
				} else if strings.HasPrefix(link, "/") {
					resolvedURL := fmt.Sprintf("%s%s", *baseUrl, link)
					relativeUrls = append(relativeUrls, strings.TrimSuffix(resolvedURL, "/"))
				}
			}
		})
	}
	return relativeUrls
}
