package browser

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
)

const (
	LangGo = "go"
	LangJS = "js"
	LangPy = "py"
)

type Browser interface {
	Browse(parentPath string)
}

func NewBrowser(lang string, tree *graph.Tree, ignoredPaths []string) Browser {
	switch lang {
	case LangGo:
		return &GoBrowser{
			ignoredPaths: ignoredPaths,
			tree:         tree,
		}
	case LangJS:
		return nil
	case LangPy:
		return nil
	}
	return nil
}
