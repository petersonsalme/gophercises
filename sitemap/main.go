package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	link "github.com/petersonsalme/gophercises/html-link-parser"
)

var (
	urlFlag  = flag.String("url", "https://gophercises.com", "the url that we want to build this sitemap for.")
	maxDepth = flag.Int("depth", 3, "the maximum number of links deep to traverse")
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type location struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	Urls  []location `xml:"url"`
	Xmlns string     `xml:"xmlns,attr"`
}

func init() {
	flag.Parse()
}

func main() {
	pages := breadthFirstSearch(*urlFlag, *maxDepth)
	toXML := transformToURLSet(pages)

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	if err := enc.Encode(&toXML); err != nil {
		panic(err)
	}

	fmt.Println()
}

func transformToURLSet(pages []string) urlSet {
	ret := urlSet{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		ret.Urls = append(ret.Urls, location{page})
	}
	return ret
}

func breadthFirstSearch(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})

	var queue map[string]struct{}
	nextQueue := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		queue, nextQueue = nextQueue, make(map[string]struct{})
		if len(queue) == 0 {
			break
		}

		for url := range queue {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range getLinksOn(url) {
				if _, ok := seen[link]; !ok {
					nextQueue[link] = struct{}{}
				}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func getLinksOn(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}

	defer resp.Body.Close()

	baseURL := transformToBaseURL(resp.Request.URL)
	base := baseURL.String()
	hrefs := hypertextReferences(resp.Body, base)

	return filter(hrefs, withPrefix(base))
}

func transformToBaseURL(reqURL *url.URL) *url.URL {
	return &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
}

func hypertextReferences(r io.Reader, baseURL string) (ret []string) {
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			ret = append(ret, baseURL+link.Href)
		case strings.HasPrefix(link.Href, "http"):
			ret = append(ret, link.Href)
		}
	}

	return
}

type filterFn func(string) bool

func filter(links []string, keepFn filterFn) (ret []string) {
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return
}

func withPrefix(prefix string) filterFn {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
