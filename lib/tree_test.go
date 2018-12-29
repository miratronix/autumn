package lib

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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

func TestChop(t *testing.T) {
	Convey("Calls PreDestroy in each leaf", t, func() {
		leaf := &child{}
		NewTree().AddLeaf(leaf).Chop()
		So(leaf.pdValue, ShouldEqual, 1)
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
	})
}
