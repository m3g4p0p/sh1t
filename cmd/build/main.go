package main

import (
	"log/slog"
	"os"

	"m3g4p0p/sh1t/internal/site"
)

func main() {
	if err := site.Generate(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
