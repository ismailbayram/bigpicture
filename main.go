package main

import (
	"bytes"
	"embed"
	"github.com/ismailbayram/bigpicture/internal/browser"
	"github.com/ismailbayram/bigpicture/internal/config"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/ismailbayram/bigpicture/internal/server"
	"os"
	"strings"
)

//go:embed web/*
var staticFiles embed.FS

type BigPicture struct {
	cfg  *config.Configuration
	tree *graph.Tree
}

func main() {
	moduleName := getModuleName()

	bp := BigPicture{
		cfg:  config.Init(),
		tree: graph.NewTree(moduleName),
	}

	brow := browser.NewBrowser(moduleName, bp.tree, bp.cfg.IgnoredPaths)

	brow.Browse(".")
	bp.tree.GenerateLinks()

	server.RunServer(staticFiles, bp.cfg.Port, bp.tree.ToJSON())
}

func getModuleName() string {
	file, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	firstLine := string(file[:bytes.IndexByte(file, '\n')])
	return strings.Split(firstLine, " ")[1]
}
