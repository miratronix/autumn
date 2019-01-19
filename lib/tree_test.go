package lib

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

type lifecylceCounter struct {
	pcCount int
	pdCount int
}

func (l *lifecylceCounter) PostConstruct() {
	l.pcCount++
}

func (l *lifecylceCounter) PreDestroy() {
	l.pdCount++
}

func TestChop(t *testing.T) {
	Convey("Calls PreDestroy in each leaf", t, func() {
		leaf := &child{}
		NewTree().AddLeaf(leaf).Chop()
		So(leaf.pdValue, ShouldEqual, 1)
	})

	Convey("Calls aliased leaf PreDestroy once", t, func() {
		leaf := &lifecylceCounter{}
		NewTree().AddNamedLeaf("a", leaf).AddNamedLeaf("b", leaf).Chop()

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

func TestAddSameLeaf(t *testing.T) {
	Convey("Adds same leaf twice", t, func() {
		Convey("Panics if adding same leaf with same name via AddLeaf", func() {
			leaf := &noop{}
			So(func() {
				NewTree().AddLeaf(leaf).AddLeaf(leaf)
			}, ShouldPanic)
		})

		Convey("Panics if adding same leaf with same name via AddNamedLeaf", func() {
			leaf := &noop{}
			So(func() {
				NewTree().AddNamedLeaf("test", leaf).AddNamedLeaf("test", leaf)
			}, ShouldPanic)
		})

		Convey("Adds alias to existing leaf if adding same leaf added with new name via AddLeaf", func() {
			leaf := &noop{}
			tree := NewTree().AddNamedLeaf("test", leaf)
			So(tree.leaves, ShouldHaveLength, 1)

			tree.AddLeaf(leaf)
			So(tree.leaves, ShouldHaveLength, 1)
		})

		Convey("Adds alias to existing leaf if adding same leaf added with new name via AddNamedLeaf", func() {
			leaf := &noop{}
			tree := NewTree().AddLeaf(leaf)
			So(tree.leaves, ShouldHaveLength, 1)

			tree.AddNamedLeaf("test", leaf)
			So(tree.leaves, ShouldHaveLength, 1)
		})
	})
}

func TestResolve(t *testing.T) {
	Convey("Resolves dependencies on leaves", t, func() {

		Convey("Handles self-injection", func() {
			tree := NewTree()
			si := &selfInject{}
			tree.AddLeaf(si).Resolve()
			So(tree.unresolved, ShouldHaveLength, 0)
			So(si.S, ShouldEqual, si)
		})

		Convey("Resolves dependencies in the supplied order", func() {
			p1 := &parent{}
			c1 := &child{}
			p2 := &parent{}
			c2 := &child{}

			NewTree().AddLeaf(p1).AddLeaf(c1).Resolve()
			NewTree().AddLeaf(c2).AddLeaf(p2).Resolve()

			So(p1.pcValue, ShouldEqual, 1)
			So(c1.pcValue, ShouldEqual, 1)

			So(p2.pcValue, ShouldEqual, 2)
			So(c2.pcValue, ShouldEqual, 1)
		})

		Convey("Handles circular dependencies", func() {
			f := circularFoo{}
			b := circularBar{}

			NewTree().AddLeaf(&f).AddLeaf(&b).Resolve()

			So(f.Bar, ShouldEqual, &b)
			So(b.Foo, ShouldEqual, &f)
		})

		Convey("Panics if a dependency can't be found", func() {
			So(func() {
				NewTree().AddLeaf(&brokenDependency{}).Resolve()
			}, ShouldPanic)
		})

		Convey("Calls aliased leaf PostConstruct once", func() {
			leaf := &lifecylceCounter{}
			NewTree().AddNamedLeaf("a", leaf).AddNamedLeaf("b", leaf).Resolve()

			So(leaf.pcCount, ShouldEqual, 1)
		})
	})
}
