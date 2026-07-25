package main

import (
	"flag"
	"log/slog"
	"os"

	"m3g4p0p/sh1t/internal/site"
	"m3g4p0p/sh1t/internal/writer"
)

var config struct {
	url string
}

func runBuild() error {
	flag.StringVar(&config.url, "url", "", "")
	flag.Parse()

	player, err := site.ExtractPlayer(config.url)
	if err != nil {
		return err
	}
	slog.Info("extract", slog.Any("player", player))
	return nil
}

func init() {
	h := slog.NewJSONHandler(writer.New(os.Stdout), nil)
	slog.SetDefault(slog.New(h))
}

func main() {
	if err := runBuild(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
