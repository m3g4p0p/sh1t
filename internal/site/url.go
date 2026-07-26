package site

import (
	"net/url"

	"m3g4p0p/sh1t/internal/urlpath"
)

type EmbedURL struct {
	*url.URL
	path *urlpath.Path
}

func ParseURL(rawURL string) (*EmbedURL, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	path := urlpath.Parse(url.Path)
	return &EmbedURL{url, path}, nil
}

func (u *EmbedURL) SetOption(k, v string) {
	u.path.Set(k, v)
	u.URL.Path = u.path.String()
}
