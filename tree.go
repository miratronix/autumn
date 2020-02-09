package autumn

// Tree defines a set of leaves
type Tree struct {
	config      *config
	leaves      map[string]*leaf
	addedLeaves []string
}

// NewTree constructs a new tree
func NewTree() *Tree {
	return &Tree{
		config:      NewConfig(),
		leaves:      make(map[string]*leaf),
		addedLeaves: make([]string, 0),
	}
}

// Configure configures the tree
func (t *Tree) Configure(config *config) *Tree {
	t.config = config
	return t
}

// AddLeaf adds a leaf to the tree
func (t *Tree) AddLeaf(value interface{}) *Tree {
	t.checkType(value)
	return t.add(newLeaf(t.config, value))
}

// AddNamedLeaf adds a named leaf to the tree
func (t *Tree) AddNamedLeaf(name string, value interface{}) *Tree {
	t.checkType(value)
	return t.add(newNamedLeaf(t.config, name, value))
}

// AddAlias adds an alias to a leaf that's already been added
func (t *Tree) AddAlias(name string, alias ...string) *Tree {

	// Make sure the source leaf exists
	leaf := t.GetLeaf(name)
	if leaf == nil {
		panic("Leaf " + name + " does not exist")
	}

	// Make sure some alternate names were supplied
	if len(alias) == 0 {
		panic("Please supply one or more aliases")
	}

	// Add each alias
	for _, a := range alias {

		// Make sure a leaf doesn't already exist with the name
		t.checkName(a)

		// Add the alias
		t.leaves[a] = leaf
	}

	return t
}

// Grow loops over the leaves in the tree, setting all dependencies
func (t *Tree) Grow() *Tree {

	// Prepare a list of unresolved leaves so we can print it if required
	unresolved := make(map[string][]string)

	// Loop over the leaves and resolve their dependencies
	for _, leafName := range t.addedLeaves {

		// Grab the actual leaf
		leaf := t.GetLeaf(leafName)

		// Resolve the dependencies for the leaf
		leaf.resolveDependencies(t)

		// If the leaf has some outstanding dependencies, store those so we can print a nice error
		if !leaf.dependenciesResolved() {
			unresolved[leaf.name] = []string{}
			for name := range leaf.unresolvedDependencies {
				unresolved[leaf.name] = append(unresolved[leaf.name], name)
			}
		}
	}

	// If we have some unresolved dependencies, print a failure
	if len(unresolved) != 0 {
		err := "Failed to wire the following dependencies: \n"
		for leaf, deps := range unresolved {
			err += "- " + leaf + " \n"
			for _, dep := range deps {
				err += "    - " + dep + "\n"
			}
		}
		panic(err)
	}

	// Loop over the leaves again and call PostConstruct
	for _, leafName := range t.addedLeaves {
		t.GetLeaf(leafName).callPostConstruct()
	}

	return t
}

// GetLeaf gets a leaf in the tree by name
func (t *Tree) GetLeaf(name string) *leaf {
	leaf, ok := t.leaves[name]
	if !ok {
		return nil
	}
	return leaf
}

// Chop chops down the tree, calling pre-destroy on all the leaves that have it in reverse order
func (t *Tree) Chop() *Tree {
	for i := len(t.addedLeaves) - 1; i >= 0; i-- {
		t.leaves[t.addedLeaves[i]].callPreDestroy()
	}
	return t
}

// CheckType checks the type of the supplied interface
func (t *Tree) checkType(value interface{}) {
	if !isStructurePointer(value) {
		panic("Please only supply structure pointers to AddLeaf/AddNamedLeaf")
	}
}

// checkName checks if the leaf name already exists
func (t *Tree) checkName(name string) {
	_, exists := t.leaves[name]
	if exists {
		panic("A leaf with name " + name + " already exists")
	}
}

// add adds a leaf to the tree
func (t *Tree) add(leaf *leaf) *Tree {

	// Make sure the name's not in use
	t.checkName(leaf.name)

	// Add the leaf to the leaf map and the ordered list
	t.leaves[leaf.name] = leaf
	t.addedLeaves = append(t.addedLeaves, leaf.name)

	return t
}
