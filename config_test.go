package autumn

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConfig(t *testing.T) {
	Convey("Constructs a new default config object", t, func() {
		c := NewConfig()
		So(c.tagName, ShouldEqual, "autumn")
		So(c.leafNameMethod, ShouldEqual, "GetLeafName")
		So(c.postConstructMethod, ShouldEqual, "PostConstruct")
		So(c.preDestroyMethod, ShouldEqual, "PreDestroy")
	})
}

func TestTagName(t *testing.T) {
	Convey("Sets the tag name", t, func() {

		c := NewConfig().TagName("test")
		So(c.tagName, ShouldEqual, "test")

		Convey("Panics if the supplied tag name is empty", func() {
			So(func() {
				NewConfig().TagName("")
			}, ShouldPanic)
		})
	})
}

func TestLeafNameMethod(t *testing.T) {
	Convey("Sets the leaf name method", t, func() {

		c := NewConfig().LeafNameMethod("Test")
		So(c.leafNameMethod, ShouldEqual, "Test")

		Convey("Panics if the supplied method name is empty", func() {
			So(func() {
				NewConfig().LeafNameMethod("")
			}, ShouldPanic)
		})

		Convey("Panics if the supplied method name isn't public", func() {
			So(func() {
				NewConfig().LeafNameMethod("getLeafName")
			}, ShouldPanic)
		})
	})
}

func TestPostConstructMethod(t *testing.T) {
	Convey("Sets the post construct method", t, func() {

		c := NewConfig().PostConstructMethod("Test")
		So(c.postConstructMethod, ShouldEqual, "Test")

		Convey("Panics if the supplied method name is empty", func() {
			So(func() {
				NewConfig().PostConstructMethod("")
			}, ShouldPanic)
		})

		Convey("Panics if the supplied method name isn't public", func() {
			So(func() {
				NewConfig().PostConstructMethod("postConstruct")
			}, ShouldPanic)
		})
	})
}

func TestPreDestroyMethod(t *testing.T) {
	Convey("Sets the pre destroy method", t, func() {

		c := NewConfig().PreDestroyMethod("Test")
		So(c.preDestroyMethod, ShouldEqual, "Test")

		Convey("Panics if the supplied method name is empty", func() {
			So(func() {
				NewConfig().PreDestroyMethod("")
			}, ShouldPanic)
		})

		Convey("Panics if the supplied method name isn't public", func() {
			So(func() {
				NewConfig().PreDestroyMethod("postConstruct")
			}, ShouldPanic)
		})
	})
}
