package router

import (
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
)

func TestNewRouterTreeNode(t *testing.T) {
	id := 42
	node := newRouterTreeNode(id)

	if node.id != id {
		t.Errorf("expected id %d, got %d", id, node.id)
	}

	if node.children == nil {
		t.Fatal("expected children map to be initialized")
	}

	if len(node.children) != 0 {
		t.Errorf("expected empty children map, got %d children", len(node.children))
	}
}

func TestAddChild(t *testing.T) {
	node := newRouterTreeNode(1)
	child := newRouterTreeNode(2)
	hash := 100

	err := node.addChild(hash, &child)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if node.children[hash] != &child {
		t.Error("expected child to be added at hash")
	}
}

func TestAddChildAlreadyExists(t *testing.T) {
	node := newRouterTreeNode(1)
	child1 := newRouterTreeNode(2)
	child2 := newRouterTreeNode(3)
	hash := 100

	err := node.addChild(hash, &child1)
	if err != nil {
		t.Fatalf("expected no error on first add, got %v", err)
	}

	err = node.addChild(hash, &child2)
	if err == nil {
		t.Error("expected error when adding child with existing hash")
	}

	if _, ok := err.(pkgerrors.ErrRouteExists); !ok {
		t.Errorf("expected ErrRouteExists, got %T", err)
	}

	if node.children[hash] != &child1 {
		t.Error("expected first child to remain unchanged")
	}
}

func TestAddWildcardChild(t *testing.T) {
	node := newRouterTreeNode(1)
	child := newRouterTreeNode(2)

	err := node.addWildcardChild(&child)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if node.children[wildcardHash] != &child {
		t.Error("expected wildcard child to be added")
	}
}

func TestAddWildcardChildAlreadyExists(t *testing.T) {
	node := newRouterTreeNode(1)
	child1 := newRouterTreeNode(2)
	child2 := newRouterTreeNode(3)

	err := node.addWildcardChild(&child1)
	if err != nil {
		t.Fatalf("expected no error on first add, got %v", err)
	}

	err = node.addWildcardChild(&child2)
	if err == nil {
		t.Error("expected error when adding wildcard child twice")
	}

	if _, ok := err.(pkgerrors.ErrRouteExists); !ok {
		t.Errorf("expected ErrRouteExists, got %T", err)
	}

	if node.children[wildcardHash] != &child1 {
		t.Error("expected first wildcard child to remain unchanged")
	}
}

func TestFind(t *testing.T) {
	testCases := []struct {
		name        string
		hash        int
		expectFound bool
	}{
		{name: "existing child", hash: 100, expectFound: true},
		{name: "non-existing child", hash: 200, expectFound: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := newRouterTreeNode(1)

			if tc.expectFound {
				child := newRouterTreeNode(2)
				err := node.addChild(tc.hash, &child)
				if err != nil {
					t.Fatalf("failed to add child: %v", err)
				}
			}

			found := node.find(tc.hash)

			if tc.expectFound && found == nil {
				t.Error("expected to find child, got nil")
			}

			if !tc.expectFound && found != nil {
				t.Errorf("expected nil, got %v", found)
			}
		})
	}
}

func TestWildcard(t *testing.T) {
	testCases := []struct {
		name        string
		addWildcard bool
	}{
		{name: "existing wildcard", addWildcard: true},
		{name: "non-existing wildcard", addWildcard: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := newRouterTreeNode(1)

			if tc.addWildcard {
				child := newRouterTreeNode(2)
				err := node.addWildcardChild(&child)
				if err != nil {
					t.Fatalf("failed to add wildcard: %v", err)
				}
			}

			found := node.wildcard()

			if tc.addWildcard && found == nil {
				t.Error("expected to find wildcard, got nil")
			}

			if !tc.addWildcard && found != nil {
				t.Errorf("expected nil, got %v", found)
			}
		})
	}
}
