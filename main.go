package main

import (
	"embed"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/browser"
	"github.com/ismailbayram/bigpicture/internal/config"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/ismailbayram/bigpicture/internal/server"
	"github.com/ismailbayram/bigpicture/internal/validators"
	"os"
)

//go:embed web/*
var staticFiles embed.FS

type BigPicture struct {
	cfg  *config.Configuration
	tree *graph.Tree
}

func (bp *BigPicture) Validate() error {
	for _, validatorConf := range bp.cfg.Validators {
		validator, err := validators.NewValidator(validatorConf.Type, validatorConf.Args, bp.tree)
		if err != nil {
			return err
		}
		if err := validator.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	bp := BigPicture{
		cfg:  config.Init(),
		tree: graph.NewTree("root"),
	}

	brow := browser.NewBrowser(browser.LangGo, bp.tree, bp.cfg.IgnoredPaths)
	brow.Browse(".")
	bp.tree.GenerateLinks()

	if os.Args[1] == "server" {
		server.RunServer(staticFiles, bp.cfg.Port, bp.tree.ToJSON())
	} else if os.Args[1] == "validate" {
		if err := bp.Validate(); err != nil {
			fmt.Println("validation failed")
			fmt.Println(err)
		}
	} else {
		fmt.Println("invalid command")
		os.Exit(1)
	}
}
