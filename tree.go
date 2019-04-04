package autumn

// Tree defines a set of leaves
type Tree struct {
	unresolved    map[string][]string
	leaves        map[uintptr]*leaf
	orderedLeaves []*leaf
	stopChannel   chan struct{}
}

// NewTree constructs a new tree
func NewTree() *Tree {
	return &Tree{
		unresolved:    make(map[string][]string),
		leaves:        make(map[uintptr]*leaf),
		orderedLeaves: []*leaf{},
		stopChannel:   nil,
	}
}

// AddLeaf adds a leaf to the tree
func (t *Tree) AddLeaf(value interface{}) *Tree {
	t.checkType(value)
	return t.add(newLeaf(value))
}

// AddNamedLeaf adds a named leaf to the tree
func (t *Tree) AddNamedLeaf(name string, value interface{}) *Tree {
	t.checkType(value)
	return t.add(newNamedLeaf(name, value))
}

// Grow loops over the leaves in the tree, setting all dependencies
func (t *Tree) Grow() {
	for _, leaf := range t.orderedLeaves {

		// Resolve the dependencies for the leaf
		leaf.resolveDependencies(t)

		// If the leaf has some outstanding dependencies, store those so we can print a nice error
		if !leaf.dependenciesResolved() {
			t.unresolved[leaf.name] = []string{}
			for name := range leaf.unresolvedDependencies {
				t.unresolved[leaf.name] = append(t.unresolved[leaf.name], name)
			}
		}
	}
	t.checkUnresolved()
}

// GetLeaf gets a leaf in the tree by name
func (t *Tree) GetLeaf(name string) *leaf {
	for _, leaf := range t.leaves {
		if leaf.hasAlias(name) {
			return leaf
		}
	}
	return nil
}

// Chop chops down the tree, calling pre-destroy on all the leaves that have it in reverse ordder
func (t *Tree) Chop() {
	for i := len(t.orderedLeaves) - 1; i >= 0; i-- {
		t.orderedLeaves[i].callPreDestroy()
	}
}

// CheckType checks the type of the supplied interface
func (t *Tree) checkType(value interface{}) {
	if !isStructurePointer(value) {
		panic("Please only supply structure pointers to AddNamedLeaf")
	}
}

// checkName checks if the leaf name already exists
func (t *Tree) checkName(name string) {
	for _, leaf := range t.leaves {
		if name == leaf.name {
			panic("Duplicate leaf name found: " + name)
		}
	}
}

// checkUnresolved checks the unresolved map of dependencies
func (t *Tree) checkUnresolved() {
	if len(t.unresolved) != 0 {
		err := "Failed to wire the following dependencies: \n"
		for leaf, deps := range t.unresolved {
			err += "- " + leaf + " \n"
			for _, dep := range deps {
				err += "    - " + dep + "\n"
			}
		}
		panic(err)
	}
}

// add adds a leaf to the tree
func (t *Tree) add(leaf *leaf) *Tree {
	t.checkName(leaf.name)

	address := leaf.structureAddress

	// If the leaf has been added before, just add a alias
	_, ok := t.leaves[address]
	if ok {
		t.leaves[address].addAlias(leaf.name)
		return t
	}

	// First time adding it
	t.leaves[address] = leaf
	t.orderedLeaves = append(t.orderedLeaves, leaf)
	return t
}
