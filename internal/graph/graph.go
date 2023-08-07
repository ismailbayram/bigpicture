package graph

import (
	"encoding/json"
	"math"
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
				t.Links = append(t.Links, NewLink(node, t.Nodes[imported], true))
			}

			parentNode := node
			for parentNode.Parent != "." && parentNode.Parent != t.Root.Path {
				parentNode = t.Nodes[parentNode.Parent]
				t.Links = append(t.Links, NewLink(parentNode, t.Nodes[imported], false))
			}
		}
	}
}

func (t *Tree) CalculateInstability() {
	for _, node := range t.Nodes {
		importCount := 0   // from node to other modules
		importedCount := 0 // from other modules to node
		for _, link := range t.Links {
			if link.From.Path == node.Path && link.To.Parent == node.Parent {
				importCount += 1
			}
			if strings.HasPrefix(link.To.Path, node.Path) && link.From.Parent == node.Parent {
				importedCount += 1
			}
		}

		instability := float64(importCount) / float64(importedCount+importCount)
		if !math.IsNaN(instability) {
			node.Instability = new(float64)
			*node.Instability = instability
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

type Function struct {
	Name      string `json:"name"`
	LineCount int    `json:"line_count"`
}

type Node struct {
	PackageName string     `json:"package_name"`
	Path        string     `json:"path"`
	Parent      string     `json:"parent"`
	Type        Type       `json:"type"`
	ImportRaw   []string   `json:"-"`
	LineCount   int        `json:"line_count"`
	Instability *float64   `json:"instability"`
	Functions   []Function `json:"functions"`
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
		Instability: nil,
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

func (n *Node) FileName() string {
	parts := strings.Split(n.Path, "/")
	return strings.Split(parts[len(parts)-1], ".")[0]
}

func (n *Node) ToJSON() string {
	data, err := json.Marshal(&struct {
		Node
		Instability *float64 `json:"instability"`
	}{
		Node:        *n,
		Instability: nil,
	})

	//data, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type Link struct {
	From      *Node `json:"from"`
	To        *Node `json:"to"`
	IsVisible bool  `json:"is_visible"`
}

func NewLink(from *Node, to *Node, isVisible bool) *Link {
	return &Link{
		From:      from,
		To:        to,
		IsVisible: isVisible,
	}
}

func (l *Link) ToJSON() string {
	data, err := json.Marshal(struct {
		From      string `json:"from"`
		To        string `json:"to"`
		IsVisible bool   `json:"is_visible"`
	}{
		From:      l.From.Path,
		To:        l.To.Path,
		IsVisible: l.IsVisible,
	})
	if err != nil {
		panic(err)
	}
	return string(data)
}
