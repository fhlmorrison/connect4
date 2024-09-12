//go:build !dev
// +build !dev

package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed public
var publicFS embed.FS

//go:embed templates
var templateFS embed.FS

func public() http.Handler {
	return http.FileServerFS(publicFS)
}

func loadTemplates() (*template.Template, error) {
	return template.ParseFS(templateFS, "**/*.html")
}
