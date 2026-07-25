package site

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var assigmentRegexp = regexp.MustCompile(`/(\w+)=(\w+)|/[^/]*`)

func ExtractPlayer(url string) error {
	doc, err := parsePage(url)
	if err != nil {
		return err
	}

	embedURL := findEmbedURL(doc)
	if embedURL == "" {
		return fmt.Errorf("no player found parsing %s", url)
	}

	slog.Info(embedURL)
	return nil
}

func parsePage(url string) (*html.Node, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return html.Parse(res.Body)
}

func attrMap(n *html.Node) map[string]string {
	m := make(map[string]string, len(n.Attr))
	for _, a := range n.Attr {
		m[a.Key] = a.Val
	}
	return m
}

func findEmbedURL(doc *html.Node) string {
	for node := range doc.Descendants() {
		if node.DataAtom != atom.Meta {
			continue
		}

		attrs := attrMap(node)
		if attrs["property"] == "og:video" {
			return attrs["content"]
		}
	}
	return ""
}
