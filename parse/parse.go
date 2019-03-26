package link

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

// get the href links from the body of the wiki page
func parse(url string) ([]Link, error) {

	doc1, err := getDiv(url)
	if err != nil {
		fmt.Println("Missing div")
		return nil, err
	}

	nodes := linkNodes(doc1)
	var links []Link
	for _, node := range nodes {
		link := buildLink(node)
		if link.Href != "" {
			links = append(links, link)
		}
	}
	return links, nil
}

// gets just the body of the wiki page
func getDiv(url string) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)

	pageSource := read(url)
	s := strings.NewReader(pageSource)
	doc, err := html.Parse(s)
	if err != nil {
		return nil, err
	}

	// searches for the body of the wiki page
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "id" {
					if attr.Val == "mw-content-text" {
						b = n
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	body := renderNode(b)
	if err != nil {
		return nil, err
	}
	bodys := strings.NewReader(body)
	doc1, err := html.Parse(bodys)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return doc1, nil
}

func renderNode(n *html.Node) string {
	if n != nil {
		var buf bytes.Buffer
		w := io.Writer(&buf)
		html.Render(w, n)
		return buf.String()
	}
	return ""
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}
