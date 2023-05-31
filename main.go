package main

import (
	"embed"
	"github.com/ismailbayram/bigpicture/internal/browser"
	"github.com/ismailbayram/bigpicture/internal/config"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/ismailbayram/bigpicture/internal/server"
)

//go:embed web/*
var staticFiles embed.FS

type BigPicture struct {
	cfg  *config.Configuration
	tree *graph.Tree
}

func main() {
	bp := BigPicture{
		cfg:  config.Init(),
		tree: graph.NewTree("root"),
	}

	brow := browser.NewBrowser(browser.LangGo, bp.tree, bp.cfg.IgnoredPaths)
	brow.Browse(".")
	bp.tree.GenerateLinks()

	server.RunServer(staticFiles, bp.cfg.Port, bp.tree.ToJSON())
}
