# Autumn

Autumn is a basic, spring-inspired dependency injection framework for Go. It's a work in progress, but some baseline functionality is there:

* Structure tag name-based wiring
* Singleton leaves (analogous to Spring Beans)
* Circular dependency resolution
* Self-injection of leaves
* `PostConstruct` functionality

Naturally, there's lots to do:

* Function based construction (similar to Springs `autowired` constructors)
* `Prototype` scope leaves
* `PostConstruct` ordering
* Intelligent type-based wiring

## Usage

Before jumping into usage, let's define some terms:

* `Leaf` - A leaf is a singleton structure pointer. You can think of it as a Spring `Bean`. It has 2 properties:
    * a name, used to wire it into other leaves. This can be set with `GetLeafName`, or by assigning a name when adding the leaf to a tree.
    * a `PostConstruct` function, which is called when dependencies have been resolved.
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
```

You can now wire it them together:
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
tree.Resolve()

// You can also set the leaf name while adding it, which overrides the leaf name defined in the structure. This is
// useful when you want to add multiple copies of the same leaf with different names
tree.AddNamedLeaf("AnotherFirst", first)
```
