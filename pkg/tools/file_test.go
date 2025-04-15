package tools

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsDir(t *testing.T) {
	Convey("Test IsDir()", t, func() {
		Convey("Test IsDir(File)", func() {
			f, err := os.CreateTemp("", "file.temp")
			So(err, ShouldBeNil)
			defer os.Remove(f.Name())
			So(IsDir(f.Name()), ShouldBeFalse)
		})
		Convey("Test IsDir(Dir)", func() {
			dir, err := os.MkdirTemp("", "dir.temp")
			So(err, ShouldBeNil)
			defer os.Remove(dir)
			So(IsDir(dir), ShouldBeTrue)
		})
		Convey("Test IsDir Not Exist", func() {
			So(IsDir("dir.null"), ShouldBeFalse)
		})
	})
}
