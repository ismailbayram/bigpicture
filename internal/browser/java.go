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
	rootDir      string
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

		importItem = b.rootDir + strings.Replace(importItem, ".", "/", -1) + ".java"
		if strings.HasSuffix(importItem, "*.java") {
			importItem = importItem[:len(importItem)-7]
		}

		imports = append(imports, importItem)
	}

	return imports
}

func (b *JavaBrowser) findFunctions(fileContent string) map[string]string {
	functions := make(map[string]string)

	functionPattern := regexp.MustCompile(`(?:public|private|protected)?\s+(?:static\s+)?\w+\s+([\w<>]+)\s+(\w+)\s*\([^)]*\)\s*(?:throws\s+\w+(?:\s*,\s*\w+)*)?\s*\{`)
	matches := functionPattern.FindAllStringSubmatch(fileContent, -1)

	for _, match := range matches {
		startIndex := strings.Index(fileContent, match[0])
		endIndex := findFunctionEndIndex(fileContent, startIndex)
		functionContent := fileContent[startIndex+len(match[0]) : endIndex]
		functions[match[2]] = functionContent
	}

	return functions
}

func findFunctionEndIndex(javaCode string, startIndex int) int {
	openBraces := 0
	for i := startIndex + 1; i < len(javaCode); i++ {
		if javaCode[i] == '{' {
			openBraces++
		} else if javaCode[i] == '}' {
			openBraces--
			if openBraces == 0 {
				return i - 2
			}
		}
	}
	return -1
}
