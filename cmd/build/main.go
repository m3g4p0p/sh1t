package main

import (
	"flag"
	"log/slog"
	"os"

	"m3g4p0p/sh1t/internal/site"
	"m3g4p0p/sh1t/internal/writer"
)

var config struct {
	urls []string
}

func runBuild() error {
	flag.Func("url", "", func(s string) error {
		config.urls = append(config.urls, s)
		return nil
	})
	flag.Parse()

	var players []*site.Player
	for _, url := range config.urls {
		slog.Info("extracting", slog.String("url", url))
		player, err := site.ExtractPlayer(url)
		if err != nil {
			return err
		}
		slog.Info("extracted", slog.Any("player", player))
		players = append(players, player)
	}

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
