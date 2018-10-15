package autumn

import "github.com/miratronix/autumn/lib"

// Tree defines a dependency tree
type Tree = *lib.Tree

// NewTree constructs a new tree
func NewTree() Tree {
	return lib.NewTree()
}
