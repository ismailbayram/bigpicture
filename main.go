package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/browser"
	"github.com/ismailbayram/bigpicture/internal/config"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"net/http"
	"os"
	"strings"
)

//go:embed web/*
var staticFiles embed.FS
var staticDir = "web"

func rootPath(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")

		if r.URL.Path == "/" {
			r.URL.Path = fmt.Sprintf("/%s/", staticDir)
		} else {
			b := strings.Split(r.URL.Path, "/")[0]
			if b != staticDir {
				r.URL.Path = fmt.Sprintf("/%s%s", staticDir, r.URL.Path)
			}
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.Init()

	moduleName := getModuleName()
	rootNode := graph.NewNode(moduleName, ".", graph.Dir, nil)
	tree := graph.NewTree(rootNode)

	if err := browser.Browse(cfg, ".", moduleName, rootNode, tree); err != nil {
		panic(err)
	}
	tree.GenerateLinks()

	var staticFS = http.FS(staticFiles)
	fs := rootPath(http.FileServer(staticFS))

	http.Handle("/", fs)

	http.HandleFunc("/graph", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(tree.ToJSON()))
	})
	err := http.ListenAndServe(":44525", nil)
	if err != nil {
		panic(err)
	}
}

func getModuleName() string {
	file, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	firstLine := string(file[:bytes.IndexByte(file, '\n')])
	return strings.Split(firstLine, " ")[1]
}
