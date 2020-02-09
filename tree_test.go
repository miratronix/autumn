package autumn

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type noop struct{}

type selfInject struct {
	S *selfInject `autumn:"selfInject"`
}

func (s *selfInject) GetLeafName() string {
	return "selfInject"
}

type brokenDependency struct {
	A string `autumn:"broken"`
}

type parent struct {
	C       *child `autumn:"child"`
	pcValue int
}

func (p *parent) PostConstruct() {
	p.pcValue = p.C.pcValue + 1
}

type child struct {
	pcValue int
	pdValue int
}

func (c *child) PostConstruct() {
	c.pcValue = 1
}

func (c *child) PreDestroy() {
	c.pdValue = 1
}

func (c *child) GetLeafName() string {
	return "child"
}

type circularFoo struct {
	Bar *circularBar `autumn:"circularBar"`
}

func (c *circularFoo) GetLeafName() string {
	return "circularFoo"
}

type circularBar struct {
	Foo *circularFoo `autumn:"circularFoo"`
}

func (c *circularBar) GetLeafName() string {
	return "circularBar"
}

type lifecycleCounter struct {
	pcCount int
	pdCount int
}

func (l *lifecycleCounter) PostConstruct() {
	l.pcCount++
}

func (l *lifecycleCounter) PreDestroy() {
	l.pdCount++
}

type aliasSelfInject struct {
	Self    *aliasSelfInject `autumn:"self"`
	This    *aliasSelfInject `autumn:"this"`
	pcCount int
}

func (a *aliasSelfInject) PostConstruct() {
	a.pcCount++
}

func TestChop(t *testing.T) {
	Convey("Calls PreDestroy in each leaf", t, func() {
		leaf := &child{}
		NewTree().AddLeaf(leaf).Chop()
		So(leaf.pdValue, ShouldEqual, 1)
	})

	Convey("Calls the aliased leaf's PreDestroy once", t, func() {
		leaf := &lifecycleCounter{}
		NewTree().AddNamedLeaf("a", leaf).AddAlias("a", "b").Chop()

		So(leaf.pdCount, ShouldEqual, 1)
	})
}

func TestAddLeaf(t *testing.T) {
	Convey("Adds a leaf", t, func() {

		Convey("Stores the leaf", func() {
			tree := NewTree()
			tree.AddLeaf(&noop{})
			So(tree.leaves, ShouldHaveLength, 1)
		})

		Convey("Panics if the supplied value is not a pointer", func() {
			So(func() {
				NewTree().AddLeaf(noop{})
			}, ShouldPanic)
		})

		Convey("Panics if the name is already taken", func() {
			So(func() {
				NewTree().AddLeaf(&noop{}).AddLeaf(&noop{})
			}, ShouldPanic)
		})
	})
}

func TestAddNamedLeaf(t *testing.T) {
	Convey("Adds a leaf by name", t, func() {

		Convey("Stores the leaf", func() {
			tree := NewTree()
			tree.AddNamedLeaf("a", &noop{})
			So(tree.leaves, ShouldHaveLength, 1)
		})

		Convey("Panics if the supplied value is not a pointer", func() {
			So(func() {
				NewTree().AddNamedLeaf("b", noop{})
			}, ShouldPanic)
		})

		Convey("Panics if the name is already taken", func() {
			So(func() {
				NewTree().AddNamedLeaf("a", &noop{}).AddNamedLeaf("a", &noop{})
			}, ShouldPanic)
		})
	})
}

func TestAddAlias(t *testing.T) {
	Convey("Adds a leaf alias", t, func() {

		Convey("Panics if the specified leaf doesn't exist", func() {
			So(func() {
				NewTree().AddAlias("a", "b")
			}, ShouldPanic)
		})

		Convey("Panics if no aliases are supplied", func() {
			So(func() {
				NewTree().AddNamedLeaf("a", &noop{}).AddAlias("a")
			}, ShouldPanic)
		})

		Convey("Panics if the alias is already taken", func() {
			So(func() {
				NewTree().AddNamedLeaf("a", &noop{}).AddNamedLeaf("b", &noop{}).AddAlias("a", "b")
			}, ShouldPanic)
		})

		Convey("Adds the alias", func() {
			tree := NewTree().AddNamedLeaf("a", &noop{}).AddAlias("a", "b")
			So(tree.addedLeaves, ShouldHaveLength, 1)
			So(tree.leaves, ShouldHaveLength, 2)
			So(tree.leaves["a"], ShouldEqual, tree.leaves["b"])
		})
	})
}

func TestAddSameLeaf(t *testing.T) {
	Convey("Adds same leaf twice", t, func() {
		Convey("Panics if the same leaf is added twice without a specified name", func() {
			leaf := &noop{}
			So(func() {
				NewTree().AddLeaf(leaf).AddLeaf(leaf)
			}, ShouldPanic)
		})

		Convey("Panics if the same leaf is added twice with the same name", func() {
			leaf := &noop{}
			So(func() {
				NewTree().AddNamedLeaf("test", leaf).AddNamedLeaf("test", leaf)
			}, ShouldPanic)
		})
	})
}

func TestGrow(t *testing.T) {
	Convey("Resolves dependencies on leaves", t, func() {

		Convey("Handles self-injection", func() {
			tree := NewTree()
			si := &selfInject{}
			tree.AddLeaf(si).Grow()
			So(si.S, ShouldEqual, si)
		})

		Convey("Resolves dependencies in the supplied order", func() {
			p1 := &parent{}
			c1 := &child{}
			p2 := &parent{}
			c2 := &child{}

			NewTree().AddLeaf(p1).AddLeaf(c1).Grow()
			NewTree().AddLeaf(c2).AddLeaf(p2).Grow()

			So(p1.pcValue, ShouldEqual, 1)
			So(c1.pcValue, ShouldEqual, 1)

			So(p2.pcValue, ShouldEqual, 2)
			So(c2.pcValue, ShouldEqual, 1)
		})

		Convey("Handles circular dependencies", func() {
			f := circularFoo{}
			b := circularBar{}

			NewTree().AddLeaf(&f).AddLeaf(&b).Grow()

			So(f.Bar, ShouldEqual, &b)
			So(b.Foo, ShouldEqual, &f)
		})

		Convey("Panics if a dependency can't be found", func() {
			So(func() {
				NewTree().AddLeaf(&brokenDependency{}).Grow()
			}, ShouldPanic)
		})

		Convey("Calls an aliased leaf's PostConstruct once", func() {
			leaf := &lifecycleCounter{}
			NewTree().AddNamedLeaf("a", leaf).AddAlias("a", "b").Grow()

			So(leaf.pcCount, ShouldEqual, 1)
		})

		Convey("Resolves aliased leaves", func() {
			leaf := &aliasSelfInject{}
			NewTree().AddNamedLeaf("selfInject", leaf).AddAlias("selfInject", "self", "this").Grow()

			So(leaf.pcCount, ShouldEqual, 1)
			So(leaf.Self, ShouldEqual, leaf)
			So(leaf.This, ShouldEqual, leaf)
		})
	})
}
