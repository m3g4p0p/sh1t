package site

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Player struct {
	Title    string
	AlbumURL string
	EmbedURL *EmbedURL
}

func ExtractPlayer(ctx context.Context, url string) (*Player, error) {
	doc, err := parsePage(ctx, url)
	if err != nil {
		return nil, err
	}

	title := findTitle(doc)
	if title == "" {
		return nil, fmt.Errorf("no title found parsing %s", url)
	}

	rawEmbedURL := findEmbedURL(doc)
	if rawEmbedURL == "" {
		return nil, fmt.Errorf("no player found parsing %s", url)
	}

	embedURL, err := ParseURL(rawEmbedURL)
	if err != nil {
		return nil, err
	}

	return &Player{
		Title:    title,
		AlbumURL: url,
		EmbedURL: embedURL,
	}, nil
}

func parsePage(ctx context.Context, url string) (*html.Node, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
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

func findTitle(doc *html.Node) string {
	for node := range doc.Descendants() {
		if node.DataAtom != atom.Title {
			continue
		}

		var b strings.Builder
		for child := range node.ChildNodes() {
			if child.Type == html.TextNode {
				b.WriteString(child.Data)
			}
		}
		return b.String()
	}

	return ""
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
