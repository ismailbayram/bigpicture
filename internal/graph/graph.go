package graph

type Node struct {
	PackageName string
	Path        string
	Parent      *Node
	Children    []*Node
	Imports     []Edge
}

func NewNode(packageName string, path string, parent *Node) *Node {
	return &Node{
		PackageName: packageName,
		Path:        path,
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

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

type Edge struct {
	From *Node
	To   *Node
}

func NewEdge(from *Node, to *Node) Edge {
	return Edge{From: from, To: to}
}
