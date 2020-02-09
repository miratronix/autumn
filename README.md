# Autumn
Autumn is a basic, spring-inspired dependency injection framework for Go. It's a work in progress, but some baseline functionality is there:

* Structure tag name-based wiring
* Singleton leaves (analogous to Spring Beans)
* Circular dependency resolution
* Self-injection of leaves
* `PostConstruct` functionality
* `PreDestroy` functionality

Naturally, there's lots to do:

* Function based construction (similar to Springs `autowired` constructors)
* `Prototype` scope leaves
* `PostConstruct` ordering
* Intelligent type-based wiring

## Usage

Before jumping into usage, let's define some terms:

* `Leaf` - A leaf is a singleton structure pointer. You can think of it as a Spring `Bean`. It has 3 properties:
    * a name, used to wire it into other leaves. This can be set with `GetLeafName`, or by assigning a name when adding the leaf to a tree.
    * an optional `PostConstruct` function, which is called when dependencies have been resolved.
    * an optional `PreDestroy` function, which is called when the tree is "chopped" (stopped).
* `Tree` - A tree contains a list of leaves, and does the heavy lifting when resolving dependencies.

So, let's say you define a leaf like so:
```go
package leaves

type FirstLeaf struct {
	SecondLeaf *SecondLeaf `autumn:"SecondLeaf"`
}

func (f *FirstLeaf) GetLeafName() string {
	return "FirstLeaf"
}

func (f *FirstLeaf) PostConstruct() {
	fmt.Println("First constructed, f.SecondLeaf is not nil here")
}

func (f *FirstLeaf) PreDestroy() {
	fmt.Println("First destroyed")
}
```

And a second one:
```go
package leaves

type SecondLeaf struct {
	FirstLeaf *FirstLeaf `autumn:"FirstLeaf"`
}

func (s *SecondLeaf) GetLeafName() string {
	return "SecondLeaf"
}

func (s *SecondLeaf) PostConstruct() {
	fmt.Println("Second constructed, s.FirstLeaf is not nil here")
}

func (s *SecondLeaf) PreDestroy() {
	fmt.Println("Second destroyed")
}
```

You can now wire them together:
```go
package leaves

first := &FirstLeaf{}
second := &SecondLeaf{}

tree := autumn.NewTree()
tree.AddLeaf(first)
tree.AddLeaf(second)

// You can now resolve the dependencies. Once this operation completes, first.SecondLeaf will point to second. and 
// second.FirstLeaf will point to first. Because of the order in which these were added, "First constructed" will be 
// printed first, followed by "Second constructed"
tree.Grow()

// You can also set the leaf name while adding it, which overrides the leaf name defined in the structure. Note that if 
// you add the leaf twice this way, its PostConstruct() function will be called twice. To avoid this, use an alias as
// described below
tree.AddNamedLeaf("AnotherFirst", first)

// To kill all the leaves in the tree, call Chop(). This is useful when gracefully shutting down an application, and
// gives each leaf a chance to clean up after itself. Post destruct for each leaf will be called once, in reverse 
// resolve order
tree.Chop()
```

### Aliasing
You can also add aliases to leaves, which are alternate names for the same leaf object. For example, lets say you define
your leaves like so:
```go
package leaves

type FirstLeaf struct {
	SecondLeaf *SecondLeaf `autumn:"second"`
}

type SecondLeaf struct {
	FirstLeaf *FirstLeaf `autumn:"someOtherName"`
}
```

and your tree like so:
```go
package leaves

// Construct the instances
first := &FirstLeaf{}
second := &SecondLeaf{}

// Add them to the tree
tree := autumn.NewTree()
tree.AddNamedLeaf("first", first)
tree.AddNamedLeaf("second", second)

// Add an alias to the first leaf. The first leaf will now be accessible as "first" or "someOtherName"
tree.AddAlias("first", "someOtherName")
```

The dependencies will be correctly resolved when the tree is grown, and the `FistLeaf.PostConstruct()` will only be called
once (if present).

### Configuration
To configure a tree, use the `Configure` function:
```go
package leaves

// Construct a new configuration object
config := autumn.NewConfig().
    TagName("autumn").                      // The tag name to use
    LeafNameMethod("GetLeafName").          // The name of the function to call to get the leaf name - must be public
    PostConstructMethod("PostConstruct").   // The name of the function to call when dependencies are resolved - must be public
    PreDestroyMethod("PreDestroy")          // The name of the function to call when the tree is chopped - must be public

// And apply it to the tree
tree := autumn.NewTree().Configure(config)
```
