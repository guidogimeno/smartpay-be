package parser

import "golang.org/x/net/html"

func HtmlIterator(n *html.Node, callback func(*html.Node)) {
	callback(n)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		HtmlIterator(c, callback)
	}
}
