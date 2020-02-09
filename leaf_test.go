package autumn

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type namedByPointer struct{}

func (a *namedByPointer) GetLeafName() string {
	return "namedByPointer"
}

type namedByValue struct{}

func (a *namedByValue) GetLeafName() string {
	return "namedByValue"
}

type withDep struct {
	A string `autumn:"a"`
}

type foo struct {
	Bar      *bar `autumn:"bar"`
	pcCalled bool
}

func (f *foo) PostConstruct() {
	f.pcCalled = true
}

type bar struct{}

func (b *bar) GetLeafName() string {
	return "bar"
}

func TestNewLeaf(t *testing.T) {
	Convey("Constructs a new leaf", t, func() {

		Convey("Sets the leaf name", func() {

			Convey("With a pointer receiver", func() {
				So(newLeaf(NewConfig(), &namedByPointer{}).name, ShouldEqual, "namedByPointer")
			})

			Convey("With a value receiver", func() {
				So(newLeaf(NewConfig(), &namedByValue{}).name, ShouldEqual, "namedByValue")
			})
		})

		Convey("Sets the unresolved leaf dependencies", func() {
			So(newLeaf(NewConfig(), &withDep{}).unresolvedDependencies, ShouldHaveLength, 1)
		})
	})
}

func TestNewNamedLeaf(t *testing.T) {
	Convey("Constructs a new named leaf", t, func() {

		Convey("Sets the leaf name", func() {

			Convey("With a pointer receiver", func() {
				So(newNamedLeaf(NewConfig(), "test", &namedByPointer{}).name, ShouldEqual, "test")
			})

			Convey("With a value receiver", func() {
				So(newNamedLeaf(NewConfig(), "test", &namedByValue{}).name, ShouldEqual, "test")
			})
		})

		Convey("Sets the unresolved leaf dependencies", func() {
			So(newNamedLeaf(NewConfig(), "test", &withDep{}).unresolvedDependencies, ShouldHaveLength, 1)
		})
	})
}

func TestResolveDependencies(t *testing.T) {
	Convey("Resolves the leaf dependencies", t, func() {

		f := &foo{}
		b := &bar{}

		fLeaf := newLeaf(NewConfig(), f)
		bLeaf := newLeaf(NewConfig(), b)

		fLeaf.resolveDependencies(NewTree().add(fLeaf).add(bLeaf))
		So(f.Bar, ShouldEqual, b)
		So(fLeaf.resolvedDependencies, ShouldHaveLength, 1)
		So(fLeaf.unresolvedDependencies, ShouldHaveLength, 0)
	})
}
