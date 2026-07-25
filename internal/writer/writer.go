package writer

import (
	"encoding/json"
	"io"
)

type jsonWriter struct {
	io.Writer
}

func New(w io.Writer) io.Writer {
	return &jsonWriter{w}
}

func (w *jsonWriter) Write(p []byte) (int, error) {
	enc := json.NewEncoder(w.Writer)
	enc.SetIndent("", "  ")
	err := enc.Encode(json.RawMessage(p))
	return len(p), err
}
