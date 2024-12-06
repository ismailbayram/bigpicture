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
	rootDir      string
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

func (b *PythonBrowser) browse(parentPath string, parentNode *graph.Node) {
	entries, err := os.ReadDir(parentPath)

	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		fName := e.Name()
		path := fmt.Sprintf("%s/%s", parentPath, fName)
		if isIgnored(b.ignoredPaths, path) {
			continue
		}

		if e.IsDir() && !strings.Contains(fName, ".") {
			node := graph.NewNode(fName, path, graph.Dir, nil)
			b.tree.Nodes[node.Path] = node
			b.browse(path, node)
			parentNode.LineCount += node.LineCount
		} else if strings.HasSuffix(fName, ".py") {
			node := b.parseFile(path, parentNode)
			b.tree.Nodes[node.Path] = node
			parentNode.LineCount += node.LineCount
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
	functionsInfo := b.findFunctions(fileContent)
	for name, functionBody := range functionsInfo {
		functions = append(functions, graph.Function{
			Name:      name,
			LineCount: len(strings.Split(functionBody, "\n")) - 1,
		})
	}

	fName := path[strings.LastIndex(path, "/")+1:]
	node := graph.NewNode(fName, path, graph.File, b.findImports(fileContent))
	node.Functions = functions
	node.LineCount = strings.Count(fileContent, "\n")

	return node
}

func (b *PythonBrowser) findImports(pythonCode string) []string {
	var imports []string

	lines := strings.Split(pythonCode, "\n")
	importRegex := regexp.MustCompile(`^\s*import\s+([^\s#]+)`)
	fromImportRegex := regexp.MustCompile(`^\s*from\s+([^\s]+)\s+import`)

	for _, line := range lines {
		var importItem string

		if matches := importRegex.FindStringSubmatch(line); len(matches) > 1 {
			importItem = b.rootDir + strings.Replace(matches[1], ".", "/", -1) + ".py"
		} else if matches := fromImportRegex.FindStringSubmatch(line); len(matches) > 1 {
			importItem = b.rootDir + strings.Replace(matches[1], ".", "/", -1) + ".py"
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

func (b *PythonBrowser) findFunctions(fileContent string) map[string]string {
	functions := make(map[string]string)
	lines := strings.Split(fileContent, "\n")
	functionRegex := regexp.MustCompile(`^\s*def\s+([a-zA-Z_]\w*)\s*\(([^)]*)\)\s*(?:->\s*[^:]+)?:\s*$`)
	classRegex := regexp.MustCompile(`^\s*class\s+`)

	var currentFunctionName string
	var currentFunctionContent strings.Builder
	var baseIndentLevel int = -1
	var inClass bool
	var classIndentLevel int = -1

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			if currentFunctionName != "" {
				currentFunctionContent.WriteString(line + "\n")
			}
			continue
		}

		// Calculate current line's indent level
		currentIndentLevel := len(line) - len(strings.TrimLeft(line, " \t"))

		// Check if this is a class definition
		if classRegex.MatchString(line) {
			inClass = true
			classIndentLevel = currentIndentLevel
			continue
		}

		// Check if we're exiting a class
		if inClass && currentIndentLevel <= classIndentLevel {
			inClass = false
			classIndentLevel = -1
		}

		// Check if this is a function definition
		if matches := functionRegex.FindStringSubmatch(line); len(matches) > 1 {
			// If we were processing a previous function, save it
			if currentFunctionName != "" {
				functions[currentFunctionName] = strings.TrimSpace(currentFunctionContent.String())
			}

			// Only process class methods (skip standalone functions)
			if inClass {
				currentFunctionName = matches[1]
				currentFunctionContent.Reset()
				currentFunctionContent.WriteString(line + "\n")
				baseIndentLevel = currentIndentLevel
			}
			continue
		}

		// If we're inside a function
		if currentFunctionName != "" {
			// Check if we're still in the function's scope
			if baseIndentLevel >= 0 && currentIndentLevel <= baseIndentLevel && trimmedLine != "" {
				// We've exited the function
				functions[currentFunctionName] = strings.TrimSpace(currentFunctionContent.String())
				currentFunctionName = ""
				currentFunctionContent.Reset()
				baseIndentLevel = -1
			} else {
				currentFunctionContent.WriteString(line + "\n")
			}
		}
	}

	// Handle the last function if exists
	if currentFunctionName != "" {
		functions[currentFunctionName] = strings.TrimSpace(currentFunctionContent.String())
	}

	return functions
}
