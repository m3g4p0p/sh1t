package main

import (
	"context"
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

	ctx := context.Background()

	var players []*site.Player
	for _, url := range config.urls {
		slog.Info("extracting", slog.String("url", url))

		player, err := site.ExtractPlayer(ctx, url)
		if err != nil {
			return err
		}

		slog.Info("extracted", slog.Any("player", player))
		player.EmbedURL.SetOption("tracklist", "true")
		player.EmbedURL.SetOption("bgcolor", "333333")
		players = append(players, player)
	}

	return site.Generate(players)
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
