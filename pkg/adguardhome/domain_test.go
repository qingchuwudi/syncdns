package adguardhome

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHostRecord_IpDiff(t *testing.T) {
	data := HostRecord{
		oldIPs: []string{"1.1.1.1", "2.2.2.2"},
		IPs:    []string{"2.2.2.2", "3.3.3.3"},
	}
	Convey("TestFindIPsDiff", t, func() {
		del, add := data.IpDiff()
		ShouldEqual(del, []string{"1.1.1.1"})
		ShouldEqual(add, []string{"3.3.3.3"})
	})
}

func TestHostRecord_RemoveDuplicate(t *testing.T) {
	data := HostRecord{
		IPs: []string{"a", "a", "b", "c", "c", "c", "d"},
	}
	Convey("TestHostRecord_RemoveDuplicate", t, func() {
		data.RemoveDuplicate()
		ShouldEqual(data.IPs, []string{"a", "b", "c", "d"})
		data.IPs = nil
		data.RemoveDuplicate()
		ShouldBeNil(data.IPs)
	})
}
