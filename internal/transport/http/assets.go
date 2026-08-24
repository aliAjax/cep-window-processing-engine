package httptransport

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var consoleAssets embed.FS

func consoleHandler() http.Handler {
	content, err := fs.Sub(consoleAssets, "static")
	if err != nil {
		panic(err)
	}
	return http.StripPrefix("/console/", http.FileServer(http.FS(content)))
}
