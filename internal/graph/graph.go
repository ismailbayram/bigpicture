package graph

import "strings"

type Type uint

const (
	Dir Type = 1 << iota
	File
)

type Tree struct {
	Root  *Node
	Nodes map[string]*Node
}

func NewTree(root *Node) *Tree {
	tree := &Tree{
		Root:  root,
		Nodes: make(map[string]*Node),
	}
	tree.Nodes[root.Path] = root
	return tree
}

func (t *Tree) ConvertImportRaw() {
	for _, node := range t.Nodes {
		for _, imported := range node.importRaw {
			node.Imports = append(node.Imports, t.Nodes[imported])
			for node.Parent != nil && node.PackageName == node.Parent.PackageName {
				node = node.Parent
				node.Imports = append(node.Imports, t.Nodes[imported])
			}
		}
	}
}

type Node struct {
	PackageName string
	Path        string
	Parent      *Node
	Children    []*Node
	Imports     []*Node
	Type        Type

	importRaw []string
}

func NewNode(packageName string, path string, parent *Node, _type Type, importRaw []string) *Node {
	if strings.HasPrefix(path, "./") {
		path = path[1:]
	}

	return &Node{
		PackageName: packageName,
		Path:        path,
		Parent:      parent,
		Type:        _type,
		importRaw:   importRaw,
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}
