package router

type Node struct {
	Part     string
	IsParam  bool
	Handlers Handlers
	Children []*Node
}

func (n *Node) matchChild(part string) *Node {
	for _, child := range n.Children {
		if child.Part == part || child.IsParam {
			return child
		}
	}
	return nil
}
