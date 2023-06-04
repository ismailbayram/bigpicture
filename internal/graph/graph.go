package graph

import (
	"encoding/json"
	"strings"
)

type Type uint

const (
	Dir Type = 1 << iota
	File
)

type Tree struct {
	Root  *Node            `json:"root"`
	Nodes map[string]*Node `json:"nodes"`
	Links []*Link          `json:"links"`
}

func NewTree(moduleName string) *Tree {
	root := NewNode(moduleName, ".", Dir, nil)

	tree := &Tree{
		Root:  root,
		Nodes: make(map[string]*Node),
	}
	tree.Nodes[root.Path] = root
	return tree
}

func (t *Tree) GenerateLinks() {
	for _, node := range t.Nodes {
		for _, imported := range node.ImportRaw {
			if node != nil && t.Nodes[imported] != nil {
				t.Links = append(t.Links, NewLink(node, t.Nodes[imported]))
			}

			//for node.Parent != "" && node.Parent != t.Root.Path {
			//	node = t.Nodes[node.Parent]
			//	t.Links = append(t.Links, NewLink(node, t.Nodes[imported]))
			//}
		}
	}
}

func (t *Tree) ToJSON() string {
	data, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type Node struct {
	PackageName string   `json:"package_name"`
	Path        string   `json:"path"`
	Parent      string   `json:"parent"`
	Type        Type     `json:"type"`
	ImportRaw   []string `json:"-"`
}

func NewNode(packageName string, path string, _type Type, importRaw []string) *Node {
	if strings.HasPrefix(path, "./") {
		path = path[1:]
	}

	node := &Node{
		PackageName: packageName,
		Path:        path,
		Type:        _type,
		ImportRaw:   importRaw,
	}
	// TODO: remove this and accept parent as parameter
	if strings.Contains(path, "/") {
		dirs := strings.Split(path, "/")
		node.Parent = strings.Join(dirs[:len(dirs)-1], "/")
		if node.Parent == "" {
			node.Parent = "."
		}
	}

	return node
}

func (n *Node) ToJSON() string {
	data, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type Link struct {
	From *Node `json:"from"`
	To   *Node `json:"to"`
}

func NewLink(from *Node, to *Node) *Link {
	return &Link{
		From: from,
		To:   to,
	}
}

func (l *Link) ToJSON() string {
	data, err := json.Marshal(struct {
		From string `json:"from"`
		To   string `json:"to"`
	}{
		From: l.From.Path,
		To:   l.To.Path,
	})
	if err != nil {
		panic(err)
	}
	return string(data)
}
