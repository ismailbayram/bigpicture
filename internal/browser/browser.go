package browser

import (
	"bytes"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"os"
	"regexp"
	"strings"
)

func Browse(path string, moduleName string, parentNode *graph.Node) error {
	entries, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() && !strings.Contains(e.Name(), ".") {
			dirPath := fmt.Sprintf("%s/%s", path, e.Name())
			node := graph.NewNode(e.Name(), dirPath, parentNode)
			parentNode.AddChild(node)
			if err := Browse(dirPath, moduleName, node); err != nil {
				return err
			}

		} else {
			fName := e.Name()
			if !strings.Contains(fName, ".go") {
				continue
			}

			filePath := fmt.Sprintf("%s/%s", path, fName)
			file, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			node := graph.NewNode(e.Name(), filePath, parentNode)
			firstLine := string(file[:bytes.IndexByte(file, '\n')])
			packageName := strings.Split(firstLine, " ")[1]
			parentNode.AddChild(node)
			parentNode.PackageName = packageName
			// TODO: make this dynamic
			re := regexp.MustCompile(`"github\.com/ismailbayram/bigpicture.*?"`)
			matches := re.FindAllString(string(file), -1)
			if len(matches) == 0 {
				continue
			}
		}
	}
	return nil
}

func browse(path string) {

}
