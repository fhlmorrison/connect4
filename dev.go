//go:build dev
// +build dev

package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func public() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("public")))
}

func loadTemplates() (*template.Template, error) {
	var templateFS = os.DirFS("templates")

	// Walk through the file system and print each file and directory name
	err := fs.WalkDir(templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path) // Print the file or directory name
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return template.ParseFS(templateFS, "**.html")
	// return template.ParseFiles("templates/index.html")
}
