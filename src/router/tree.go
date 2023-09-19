package router

import (
	"net/http"
	"strings"
)

type Tree struct {
	root *Node
}

func (router *Tree) Add(path, method string, handler http.HandlerFunc) {
	parts := strings.Split(path, "/")[1:]
	current := router.root
	for _, part := range parts {
		isParam := strings.HasPrefix(part, ":")
		node := current.matchChild(part)
		if node == nil {
			node = &Node{Part: part, IsParam: isParam}
			current.Children = append(current.Children, node)
		}
		current = node
	}
	if current.Handlers == nil {
		current.Handlers = make(Handlers)
	}
	current.Handlers[method] = handler
}

func (router *Tree) matchNode(node *Node, parts []string, index int) *Node {
	if index == len(parts) {
		return node
	}
	for _, child := range node.Children {
		if child.Part == parts[index] || child.IsParam {
			matched := router.matchNode(child, parts, index+1)
			if matched != nil {
				return matched
			}
		}
	}
	return nil
}

func (router *Tree) GetHandler(path, method string) (http.HandlerFunc, bool) {
	parts := strings.Split(path, "/")[1:]
	node := router.matchNode(router.root, parts, 0)
	if node == nil {
		return nil, false
	}
	handler, exists := node.Handlers[method]
	return handler, exists
}
