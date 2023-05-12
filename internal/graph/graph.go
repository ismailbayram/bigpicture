package graph

type Node struct {
	PackageName string
	Dir         string
	Parent      *Node
	Imports     []Edge
}

func NewNode(packageName string, dir string, parent *Node) *Node {
	return &Node{
		PackageName: packageName,
		Dir:         dir,
		Parent:      parent,
	}
}

func (n *Node) AddImport(imported *Node) {
	edge := NewEdge(n, imported)
	n.Imports = append(n.Imports, edge)

	node := n
	for node.Parent != nil {
		node = node.Parent
		edge := NewEdge(n, imported)
		node.Imports = append(node.Imports, edge)
	}
}

type Edge struct {
	From *Node
	To   *Node
}

func NewEdge(from *Node, to *Node) Edge {
	return Edge{From: from, To: to}
}
