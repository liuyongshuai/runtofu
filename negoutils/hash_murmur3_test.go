// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        hash_murmur3_test.go
// @date        2025-03-25 10:58

package negoutils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_ranklib(t *testing.T) {
	s := "2000000000059236943"
	hv := MurmurHash3([]byte(s))
	Convey("Hash Value", t, func() {
		So(hv, ShouldEqual, 1498636499)
	})
}
