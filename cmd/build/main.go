package main

import (
	"flag"
	"log/slog"
	"os"

	"m3g4p0p/sh1t/internal/site"
)

var config struct {
	url string
}

func runBuild() error {
	flag.StringVar(&config.url, "url", "", "")
	flag.Parse()

	return site.ExtractPlayer(config.url)
}

func main() {
	if err := runBuild(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
