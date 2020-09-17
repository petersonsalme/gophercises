package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an HTML Link element
type Link struct {
	Href string
	Text string
}

func Parse(reader io.Reader) ([]Link, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		panic(err)
	}

	var links []Link
	for _, node := range linkNodes(doc) {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func buildLink(node *html.Node) (ret Link) {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			ret.Text = text(node)
			break
		}
	}
	return
}

func text(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.ElementNode {
		return ""
	}

	var ret string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		ret += text(child)
	}

	return strings.Join(strings.Fields(ret), " ")
}

func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}

	var ret []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		ret = append(ret, linkNodes(child)...)
	}

	return ret
}
