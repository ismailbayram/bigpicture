package browser

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"os"
	"regexp"
	"strings"
)

type PythonBrowser struct {
	ignoredPaths []string
	moduleName   string
	tree         *graph.Tree
}

func (b *PythonBrowser) Browse(parentPath string) {
	b.moduleName = b.getModuleName()
	b.browse(parentPath, b.tree.Root)
	b.clearNonProjectImports()
}

func (b *PythonBrowser) getModuleName() string {
	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return strings.Split(directory, "/")[len(strings.Split(directory, "/"))-1]
}

func (b *PythonBrowser) clearNonProjectImports() {
	for _, node := range b.tree.Nodes {
		var clearedImports []string
		for _, imp := range node.ImportRaw {
			if _, ok := b.tree.Nodes[imp]; ok {
				clearedImports = append(clearedImports, imp)
			}
		}
		node.ImportRaw = clearedImports
	}
}

func (b *PythonBrowser) isIgnored(entryPath string) bool {
	isIgnored := false
	for _, ignore := range b.ignoredPaths {
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

func (b *PythonBrowser) browse(parentPath string, parentNode *graph.Node) {
	entries, err := os.ReadDir(parentPath)

	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		fName := e.Name()
		path := fmt.Sprintf("%s/%s", parentPath, fName)
		if b.isIgnored(path) {
			continue
		}

		if e.IsDir() && !strings.Contains(fName, ".") {
			node := graph.NewNode(fName, path, graph.Dir, nil)
			b.tree.Nodes[node.Path] = node
			b.browse(path, node)
		} else if strings.HasSuffix(fName, ".py") {
			node := b.parseFile(path, parentNode)
			b.tree.Nodes[node.Path] = node
		}
	}
}

func (b *PythonBrowser) parseFile(path string, parentNode *graph.Node) *graph.Node {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fileContent := string(file)

	// extract functions
	functions := make([]graph.Function, 0)
	functionsInfo := findFunctions(fileContent)
	for name, functionBody := range functionsInfo {
		functions = append(functions, graph.Function{
			Name:      name,
			LineCount: len(strings.Split(functionBody, "\n")) - 1,
		})
	}

	fName := path[strings.LastIndex(path, "/")+1:]
	node := graph.NewNode(fName, path, graph.File, findImports(fileContent))
	node.Functions = functions
	node.LineCount = strings.Count(fileContent, "\n")

	return node
}

func findImports(pythonCode string) []string {
	var imports []string

	lines := strings.Split(pythonCode, "\n")
	importRegex := regexp.MustCompile(`^\s*import\s+([^\s#]+)`)
	fromImportRegex := regexp.MustCompile(`^\s*from\s+([^\s]+)\s+import`)

	for _, line := range lines {
		var importItem string

		if matches := importRegex.FindStringSubmatch(line); len(matches) > 1 {
			importItem = "/" + strings.Replace(matches[1], ".", "/", -1) + ".py"
		} else if matches := fromImportRegex.FindStringSubmatch(line); len(matches) > 1 {
			importItem = "/" + strings.Replace(matches[1], ".", "/", -1) + ".py"
		}
		if importItem == "" {
			continue
		}

		if strings.HasSuffix(importItem, "*.py") {
			importItem = importItem[:len(importItem)-5]
		}
		imports = append(imports, importItem)
	}

	return imports
}

func findFunctions(fileContent string) map[string]string {
	functions := make(map[string]string)

	lines := strings.Split(fileContent, "\n")
	functionRegex := regexp.MustCompile(`^\s*def\s+(\w+)\s*\((.*?)\):`)

	var currentFunctionName string
	var currentFunctionContent string

	for _, line := range lines {
		if matches := functionRegex.FindStringSubmatch(line); len(matches) > 1 {
			if currentFunctionName != "" {
				functions[currentFunctionName] = strings.TrimSpace(currentFunctionContent)
				currentFunctionContent = ""
			}

			currentFunctionName = matches[1]
		}

		if currentFunctionName != "" && (strings.Contains(line, "@") || strings.Contains(line, "class ")) {
			functions[currentFunctionName] = strings.TrimSpace(currentFunctionContent)
			currentFunctionContent = ""
			currentFunctionName = ""
			continue
		}

		if currentFunctionName != "" {
			currentFunctionContent += line
			currentFunctionContent += "\n"
		}
	}

	if currentFunctionName != "" {
		functions[currentFunctionName] = strings.TrimSpace(currentFunctionContent)
	}

	return functions
}
