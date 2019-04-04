package autumn

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsStructurePointer(t *testing.T) {
	Convey("Correctly identifies structure pointers", t, func() {

		Convey("Returns false when a primitive is supplied", func() {
			val := isStructurePointer(5)
			So(val, ShouldBeFalse)
		})

		Convey("Returns true when a structure is supplied by value", func() {
			type testStruct struct{}

			val := isStructurePointer(testStruct{})
			So(val, ShouldBeFalse)
		})

		Convey("Returns false when a structure is supplied by reference", func() {
			type testStruct struct{}

			val := isStructurePointer(&testStruct{})
			So(val, ShouldBeTrue)
		})
	})
}
