package site

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed templates/*
var fs embed.FS

var baseTpl = template.Must(template.ParseFS(
	fs,
	"templates/base.html",
	"templates/partials/*.html",
))

var indexTpl = template.Must(
	template.Must(baseTpl.Clone()).
		ParseFS(fs, "templates/pages/index.html"),
)

type tplData struct {
	Players []*Player
}

func Generate(players []*Player) error {
	f, err := os.Create(filepath.Join("docs", "index.html"))
	if err != nil {
		return err
	}
	return indexTpl.Execute(f, tplData{players})
}
