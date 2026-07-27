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
	return &EmbedURL{URL: url}, nil
}

func (u *EmbedURL) SetOption(k, v string) {
	if u.path == nil {
		u.path = urlpath.Parse(u.Path)
	}
	u.path.Set(k, v)
	u.URL.Path = u.path.String()
}

func (u *EmbedURL) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

func (u *EmbedURL) UnmarshalText(b []byte) error {
	url, err := url.Parse(string(b))
	if err != nil {
		return err
	}
	u.URL = url
	return nil
}
