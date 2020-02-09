package autumn

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type testStruct struct{}

func TestIsStructurePointer(t *testing.T) {
	Convey("Correctly identifies structure pointers", t, func() {

		Convey("Returns false when a primitive is supplied", func() {
			So(isStructurePointer(5), ShouldBeFalse)
		})

		Convey("Returns true when a structure is supplied by value", func() {
			So(isStructurePointer(testStruct{}), ShouldBeFalse)
		})

		Convey("Returns false when a structure is supplied by reference", func() {
			So(isStructurePointer(&testStruct{}), ShouldBeTrue)
		})
	})
}
