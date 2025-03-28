/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     es
 * @date        2018-06-25 16:03
 */
package negoutils

import (
	"fmt"
	"github.com/kr/pretty"
	"testing"
)

func TestMakeElemType(t *testing.T) {
	testStart()

	data := "435"
	ref := MakeElemType(data)
	if !ref.IsString() {
		t.Errorf("call IsString failed")
	}
	idata, err := ref.ToInt()
	fmt.Println(idata, err)
	m := map[interface{}]interface{}{
		"k1":    "k1val",
		2:       "k2val",
		"k3":    "k3val",
		4:       "k4val",
		3.44444: 5.88888,
	}
	md := MakeElemType(m)
	sm, err := md.ToSlice()
	fmt.Println(err)
	fmt.Printf("%# v\n", pretty.Formatter(sm))
	fmt.Println(ref.IsSimpleType(), md.IsComplexType())

	testEnd()
}
