package router

import (
	"ducky/http/pkg/pkgerrors"
)

const wildcardHash = -1

func newRouterTreeNode(id int) routerTreeNode {
	return routerTreeNode{
		id:       id,
		children: make(map[int]*routerTreeNode),
	}
}

type routerTreeNode struct {
	id       int
	children map[int]*routerTreeNode
}

func (node *routerTreeNode) addChild(hash int, child *routerTreeNode) error {
	if _, exists := node.children[hash]; exists {
		return &pkgerrors.ErrRouteExists{}
	}
	node.children[hash] = child
	return nil
}

func (node *routerTreeNode) addWildcardChild(child *routerTreeNode) error {
	if _, exists := node.children[wildcardHash]; exists {
		return &pkgerrors.ErrRouteExists{Route: "*"}
	}
	node.children[wildcardHash] = child
	return nil
}

func (node *routerTreeNode) find(hash int) *routerTreeNode {
	if child, exists := node.children[hash]; exists {
		return child
	}
	return nil
}

func (node *routerTreeNode) wildcard() *routerTreeNode {
	if child, exists := node.children[wildcardHash]; exists {
		return child
	}
	return nil
}
