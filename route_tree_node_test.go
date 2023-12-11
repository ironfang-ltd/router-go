package router

import (
	"testing"
)

func TestRouteTreeNode_FindNode(t *testing.T) {

	rt := newRouteTreeNode()

	rt.GetOrCreateNode("/")

	node, _ := rt.Find("/")

	if node == nil {
		t.Error("node is nil")
	}
}

func TestRouteTreeNode_FindNodeNested(t *testing.T) {

	rt := newRouteTreeNode()

	rt.GetOrCreateNode("/users/create")

	node, _ := rt.Find("/users/create")

	if node == nil {
		t.Fatal("node is nil")
	}

	if node.parent != nil && node.parent.segment != "users" {
		t.Fatal("node has no parent or parent segment is not 'users'")
	}

	if node.segment != "create" {
		t.Fatal("node segment is not 'create'")
	}
}

func BenchmarkFindRoot(b *testing.B) {

	rt := newRouteTreeNode()

	rt.GetOrCreateNode("/")

	for n := 0; n < b.N; n++ {
		rt.Find("/")
	}
}

func BenchmarkFindWithSegment(b *testing.B) {

	rt := newRouteTreeNode()

	rt.GetOrCreateNode("/users")

	for n := 0; n < b.N; n++ {
		rt.Find("/users")
	}
}

func BenchmarkFindWithTwoSegments(b *testing.B) {

	rt := newRouteTreeNode()

	rt.GetOrCreateNode("/users/create")

	for n := 0; n < b.N; n++ {
		rt.Find("/users/create")
	}
}
