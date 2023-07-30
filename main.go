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

	brow := browser.NewBrowser(detectLanguage(), bp.tree, bp.cfg.IgnoredPaths)
	brow.Browse(".")
	bp.tree.GenerateLinks()
	bp.tree.CalculateInstability()

	switch os.Args[1] {
	case "server":
		server.RunServer(staticFiles, bp.cfg.Port, bp.tree.ToJSON())
	case "validate":
		if err := bp.Validate(); err != nil {
			fmt.Println("validation failed")
			fmt.Println(err)
			os.Exit(1)
		}
	case "help":
		printHelp()
	default:
		fmt.Println("Invalid command. Use 'bigpicture help' to see available commands.")
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("server")
	fmt.Println("\tRuns the web server on port which is defined in .big.picture.json file. Default port is 44525.")

	fmt.Println("\nvalidate")
	fmt.Println("\tValidates the project structure. It checks if the project structure is valid according to the validators which are defined in .big.picture.json file.")

	fmt.Println("\nhelp")
	fmt.Println("\tPrints this help message.")
}

func detectLanguage() string {
	if _, err := os.Stat("go.mod"); !os.IsNotExist(err) {
		return browser.LangGo
	}
	if _, err := os.Stat("requirements.txt"); !os.IsNotExist(err) {
		return browser.LangPy
	}
	fmt.Println("Could not detect the project language. Please run the command in the root directory of the project.")
	os.Exit(1)
	return ""
}
