package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"m3g4p0p/sh1t/internal/caching"
	"m3g4p0p/sh1t/internal/pipeline"
	"m3g4p0p/sh1t/internal/site"
	"m3g4p0p/sh1t/internal/writer"
)

var config struct {
	cacheDir string
	urls     []string
}

type deployment struct {
	startTime time.Time
	players   []*site.Player
	client    *http.Client
}

func (d *deployment) start(context.Context) (pipeline.Step, error) {
	d.startTime = time.Now()
	d.players = make([]*site.Player, len(config.urls))

	if config.cacheDir != "" {
		d.client = &http.Client{
			Transport: caching.NewCacheTransport(config.cacheDir, nil),
		}
	} else {
		d.client = http.DefaultClient
	}

	return pipeline.Sequence(d.collectPlayers, d.generatePage), nil
}

func (d *deployment) collectPlayers(context.Context) (pipeline.Step, error) {
	var steps []pipeline.Step
	for i, url := range config.urls {
		steps = append(steps, d.extractPlayer(i, url))
	}
	return pipeline.ParallelN(3, steps...), nil
}

func (d *deployment) extractPlayer(i int, url string) pipeline.Step {
	return func(ctx context.Context) (pipeline.Step, error) {
		slog.Info("extracting", slog.String("url", url))

		player, err := site.ExtractPlayer(ctx, d.client, url)
		if err != nil {
			return nil, err
		}

		slog.Info("extracted", slog.Any("player", player))
		player.EmbedURL.SetOption("tracklist", "true")
		player.EmbedURL.SetOption("bgcolor", "333333")
		d.players[i] = player

		return nil, nil
	}
}

func (d *deployment) generatePage(context.Context) (pipeline.Step, error) {
	return d.finish, site.Generate(d.players)
}

func (d *deployment) finish(context.Context) (pipeline.Step, error) {
	elapsed := time.Since(d.startTime)

	slog.Info("Build finished.", slog.GroupAttrs(
		"metrics",
		slog.Float64("duration", elapsed.Seconds()),
	))

	return nil, nil
}

func runBuild() error {
	flag.StringVar(&config.cacheDir, "cache-dir", "", "")
	flag.Func("url", "", func(s string) error {
		config.urls = append(config.urls, s)
		return nil
	})
	flag.Parse()

	ctx := context.Background()
	dep := &deployment{}

	return pipeline.Run(ctx, dep.start)
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
