package lib

// Tree defines a set of leaves
type Tree struct {
	unresolved map[string][]string
	leaves     []*Leaf
}

// NewTree constructs a new tree
func NewTree() *Tree {
	return &Tree{
		unresolved: map[string][]string{},
		leaves:     []*Leaf{},
	}
}

// AddLeaf adds a leaf to the tree
func (t *Tree) AddLeaf(value interface{}) *Tree {
	t.checkType(value)
	return t.add(NewLeaf(value))
}

// AddNamedLead adds a named leaf to the tree
func (t *Tree) AddNamedLeaf(name string, value interface{}) *Tree {
	t.checkType(value)
	return t.add(NewNamedLeaf(name, value))
}

// Resolves loops over the leaves in the tree, setting all dependencies
func (t *Tree) Resolve() {
	for _, leaf := range t.leaves {

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
func (t *Tree) GetLeaf(name string) *Leaf {
	for _, leaf := range t.leaves {
		if leaf.name == name {
			return leaf
		}
	}
	return nil
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
func (t *Tree) add(leaf *Leaf) *Tree {
	t.checkName(leaf.name)
	t.leaves = append(t.leaves, leaf)
	return t
}
