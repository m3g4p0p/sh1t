package site

import (
	"fmt"
	"net/http"
	"net/url"

	"m3g4p0p/sh1t/internal/urlpath"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Player struct {
	Title    string
	AlbumURL string
	EmbedURL string
}

func ExtractPlayer(url string) (*Player, error) {
	doc, err := parsePage(url)
	if err != nil {
		return nil, err
	}

	title := findMetaProperty(doc, "og:title")
	if title == "" {
		return nil, fmt.Errorf("no title found parsing %s", url)
	}

	embedURL := findMetaProperty(doc, "og:video")
	if embedURL == "" {
		return nil, fmt.Errorf("no player found parsing %s", url)
	}

	processed, err := processEmbedURL(embedURL)
	if err != nil {
		return nil, err
	}

	return &Player{
		Title:    title,
		AlbumURL: url,
		EmbedURL: processed,
	}, nil
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

func findMetaProperty(doc *html.Node, name string) string {
	for node := range doc.Descendants() {
		if node.DataAtom != atom.Meta {
			continue
		}

		attrs := attrMap(node)
		if attrs["property"] == name {
			return attrs["content"]
		}
	}

	return ""
}

func processEmbedURL(rawURL string) (string, error) {
	embedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	path := urlpath.Parse(embedURL.Path)
	path.Set("tracklist", "true")
	path.Set("bgcol", "333333")
	embedURL.Path = path.String()

	return embedURL.String(), nil
}
