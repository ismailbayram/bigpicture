package browser

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"regexp"
	"strings"
)

const (
	LangGo   = "go"
	LangPy   = "py"
	LangJava = "java"
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
	case LangPy:
		return &PythonBrowser{
			ignoredPaths: ignoredPaths,
			tree:         tree,
		}
	case LangJava:
		return &JavaBrowser{
			ignoredPaths: ignoredPaths,
			tree:         tree,
		}
	}
	return nil
}

func isIgnored(ignoredPaths []string, entryPath string) bool {
	isIgnored := false
	for _, ignore := range ignoredPaths {
		regxp := ignore
		if strings.HasPrefix(ignore, "*") {
			regxp = fmt.Sprintf("^%s$", ignore)
		}
		re := regexp.MustCompile(regxp)
		if re.MatchString(entryPath) {
			isIgnored = true
			break
		}
	}

	return isIgnored
}
