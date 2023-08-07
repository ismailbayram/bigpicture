package browser

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"os"
	"regexp"
	"strings"
)

type JavaBrowser struct {
	ignoredPaths []string
	moduleName   string
	tree         *graph.Tree
}

func (b *JavaBrowser) Browse(parentPath string) {
	b.moduleName = b.getModuleName()
	b.browse(parentPath, b.tree.Root)
	b.clearNonProjectImports()
}

func (b *JavaBrowser) getModuleName() string {
	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return strings.Split(directory, "/")[len(strings.Split(directory, "/"))-1]
}

func (b *JavaBrowser) clearNonProjectImports() {
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

func (b *JavaBrowser) browse(parentPath string, parentNode *graph.Node) {
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
		} else if strings.HasSuffix(fName, ".java") {
			node := b.parseFile(path, parentNode)
			b.tree.Nodes[node.Path] = node
		}
	}
}

func (b *JavaBrowser) parseFile(path string, parentNode *graph.Node) *graph.Node {
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

func (b *JavaBrowser) findImports(javaCode string) []string {
	var imports []string

	importPattern := regexp.MustCompile(`import\s+([^;\n]+);`)

	matches := importPattern.FindAllSubmatch([]byte(javaCode), -1)
	for _, match := range matches {
		importItem := string(match[1])
		if strings.HasPrefix(importItem, "java.") {
			continue
		}

		importItem = "/" + strings.Replace(importItem, ".", "/", -1) + ".java"
		if strings.HasSuffix(importItem, "*.java") {
			importItem = importItem[:len(importItem)-7]
		}

		imports = append(imports, importItem)
	}
	fmt.Println(imports)

	return imports
}

func (b *JavaBrowser) findFunctions(fileContent string) map[string]string {
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
