package router

func Make() *Tree {
	return &Tree{root: &Node{}}
}
