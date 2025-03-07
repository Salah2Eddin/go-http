package router

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/uri"
)

const routeNotFoundID = -1
const routeAlreadyExistsID = -2

//	 hashingPrime
//		used in hashing for router tree
//		we need a hashingPrime bigger then character set size.
//		uri charset is ascii which has size 128 (0-127)
//		so we use first prime > 128
const hashingPrime = 131

type RoutesTree struct {
	root        routerTreeNode
	hasher      hasher
	idGenerator idGenerator
}

func NewRoutesTree() RoutesTree {
	return RoutesTree{
		root:        newRouterTreeNode(0),
		hasher:      getHasher(hashingPrime),
		idGenerator: getIDGenerator(),
	}
}

func (tree *RoutesTree) newTreeNode() *routerTreeNode {
	newChild := newRouterTreeNode(tree.idGenerator())
	return &newChild
}

func (tree *RoutesTree) find(uri uri.Uri, allowWildcard bool) (int, error) {
	path := uri.GetPath()
	current := &tree.root

	for _, uriPart := range path {
		partHash := tree.hasher(uriPart)
		next := current.find(partHash)
		if next == nil {
			if allowWildcard || isWildcard(uriPart) {
				next = current.wildcard()
			}
			if next == nil {
				return routeNotFoundID, errors.ErrRouteNotFound{}
			}
		}
		current = next
	}
	return current.id, nil
}

func (tree *RoutesTree) getOrCreateTreeNode(current *routerTreeNode, name string) (*routerTreeNode, error) {
	// Check if the next node already exists
	nameHash := tree.hasher(name)
	next := current.find(nameHash)
	if next == nil && isWildcard(name) {
		next = current.wildcard()
	}
	// If no matching child node exists, create a new one
	if next == nil {
		next = tree.newTreeNode()
		if isWildcard(name) {
			err := current.addWildcardChild(next)
			if err != nil {
				return nil, err
			}
		} else {
			err := current.addChild(nameHash, next)
			if err != nil {
				return nil, err
			}
		}
	}
	return next, nil
}

func (tree *RoutesTree) addRoute(uri uri.Uri) (int, error) {
	path := uri.GetPath()
	current := &tree.root

	for _, uriPart := range path {
		next, err := tree.getOrCreateTreeNode(current, uriPart)
		if err != nil {
			return routeAlreadyExistsID, err
		}
		current = next
	}
	return current.id, nil
}
