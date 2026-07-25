package site

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed templates/*
var fs embed.FS

var indexTpl = template.Must(template.ParseFS(fs, "templates/index.html"))

func Generate() error {
	f, err := os.Create(filepath.Join("docs", "index.html"))
	if err != nil {
		return err
	}
	return indexTpl.Execute(f, nil)
}
