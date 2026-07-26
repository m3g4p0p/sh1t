package caching

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type CacheTransport struct {
	cacheDir string
	wrapped  http.RoundTripper
}

func NewCacheTransport(cacheDir string, wrapped http.RoundTripper) http.RoundTripper {
	if wrapped == nil {
		wrapped = http.DefaultTransport
	}
	return &CacheTransport{cacheDir: cacheDir, wrapped: wrapped}
}

func (t *CacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet {
		return t.wrapped.RoundTrip(req)
	}

	name := t.cacheFilename(req.URL.String())
	if res, err := t.readFromCache(name, req); err != nil {
		slog.Error(err.Error())
	} else {
		return res, err
	}

	res, err := t.wrapped.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 200 && res.StatusCode < 400 {
		if err := t.writeToCache(name, res); err != nil {
			slog.Error(err.Error())
		}
	}
	return res, nil
}

func (t *CacheTransport) readFromCache(name string, req *http.Request) (*http.Response, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)
	return http.ReadResponse(bufio.NewReader(r), req)
}

func (t *CacheTransport) writeToCache(name string, res *http.Response) error {
	defer res.Body.Close()

	var buf bytes.Buffer
	if err := res.Write(&buf); err != nil {
		return err
	}
	if err := os.WriteFile(name, buf.Bytes(), os.ModePerm); err != nil {
		return err
	}

	res.Body = io.NopCloser(&buf)
	return nil
}

func (t *CacheTransport) cacheFilename(url string) string {
	hash := sha256.Sum256([]byte(url))
	name := base64.RawURLEncoding.EncodeToString(hash[:])
	return filepath.Join(t.cacheDir, name)
}
