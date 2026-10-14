package router

import (
	"errors"
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/uri"
)

func TestNewRoutesTree(t *testing.T) {
	tree := NewRoutesTree()

	if tree.root.id != 0 {
		t.Errorf("expected root id to be 0, got %d", tree.root.id)
	}

	if tree.hasher == nil {
		t.Fatal("expected hasher to be initialized")
	}

	if tree.idGenerator == nil {
		t.Fatal("expected idGenerator to be initialized")
	}
}

func TestNewTreeNode(t *testing.T) {
	tree := NewRoutesTree()

	node1 := tree.newTreeNode()
	if node1 == nil {
		t.Fatal("newTreeNode() returned nil, expected non-nil node")
	}

	node2 := tree.newTreeNode()
	if node2 == nil {
		t.Fatal("newTreeNode() returned nil, expected non-nil node")
	}

	if node1.id == node2.id {
		t.Errorf("expected different IDs, got %d for both nodes", node1.id)
	}
}

func TestGetOrCreateTreeNodeNew(t *testing.T) {
	testCases := []struct {
		name    string
		segment string
	}{
		{name: "regular segment", segment: "users"},
		{name: "wildcard segment", segment: "*"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := NewRoutesTree()
			parent := &tree.root

			node, err := tree.getOrCreateTreeNode(parent, tc.segment)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if node == nil {
				t.Fatal("getOrCreateTreeNode() returned nil, expected non-nil node")
			}
		})
	}
}

func TestGetOrCreateTreeNodeExisting(t *testing.T) {
	testCases := []struct {
		name    string
		segment string
	}{
		{name: "regular segment", segment: "users"},
		{name: "wildcard segment", segment: "*"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := NewRoutesTree()
			parent := &tree.root

			node1, err := tree.getOrCreateTreeNode(parent, tc.segment)
			if err != nil {
				t.Fatalf("first call failed: %v", err)
			}

			node2, err := tree.getOrCreateTreeNode(parent, tc.segment)
			if err != nil {
				t.Fatalf("second call failed: %v", err)
			}

			if node1 != node2 {
				t.Error("expected same node instance")
			}
		})
	}
}

func TestAddRoute(t *testing.T) {
	tree := NewRoutesTree()
	testURI := uri.NewUri("/users/profile")

	id, err := tree.addRoute(testURI)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}
}

func TestAddRouteSameTwice(t *testing.T) {
	tree := NewRoutesTree()
	testURI := uri.NewUri("/users/profile")

	id1, err := tree.addRoute(testURI)
	if err != nil {
		t.Fatalf("first add failed: %v", err)
	}

	id2, err := tree.addRoute(testURI)
	if err != nil {
		t.Errorf("second add failed: %v", err)
	}

	if id1 != id2 {
		t.Errorf("expected same id, got %d and %d", id1, id2)
	}
}

func TestFindExisting(t *testing.T) {
	tree := NewRoutesTree()
	testURI := uri.NewUri("/users/profile")

	expectedID, err := tree.addRoute(testURI)
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}

	foundID, err := tree.find(testURI, false)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if foundID != expectedID {
		t.Errorf("expected id %d, got %d", expectedID, foundID)
	}
}

func TestFindNotFound(t *testing.T) {
	testCases := []struct {
		name       string
		setupRoute string
		findRoute  string
	}{
		{name: "completely missing", setupRoute: "", findRoute: "/users/profile"},
		{name: "partial exists", setupRoute: "/users", findRoute: "/users/profile"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := NewRoutesTree()

			if tc.setupRoute != "" {
				setupURI := uri.NewUri(tc.setupRoute)
				_, err := tree.addRoute(setupURI)
				if err != nil {
					t.Fatalf("failed to add setup route: %v", err)
				}
			}

			findURI := uri.NewUri(tc.findRoute)

			foundID, err := tree.find(findURI, false)
			if err == nil {
				t.Error("expected error for non-existing route")
			}

			if foundID != routeNotFoundID {
				t.Errorf("expected routeNotFoundID (%d), got %d", routeNotFoundID, foundID)
			}

			var errRouteNotFound pkgerrors.ErrRouteNotFound
			if !errors.As(err, &errRouteNotFound) {
				t.Errorf("expected ErrRouteNotFound, got %T", err)
			}
		})
	}
}

func TestFindWildcardMatch(t *testing.T) {
	tree := NewRoutesTree()
	wildcardURI := uri.NewUri("/users/*")

	expectedID, err := tree.addRoute(wildcardURI)
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}

	specificURI := uri.NewUri("/users/123")

	foundID, err := tree.find(specificURI, true)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if foundID != expectedID {
		t.Errorf("expected id %d, got %d", expectedID, foundID)
	}
}
