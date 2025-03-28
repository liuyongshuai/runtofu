// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-12-07 11:12

package negoutils

import (
	"fmt"
	"testing"
)

func TestNewRoutineTool(t *testing.T) {
	testStart()
	rt := NewRoutineTool(100, 100000, testForRoutineFunc)
	for i := 0; i < 99; i++ {
		rt.AddArg(fmt.Sprintf("第 %d 个数据", i))
	}
	rt.AddFileLines([]string{"./snowflake.txt"})
	rt.Wait()
	testEnd()
}

func testForRoutineFunc(arg interface{}, retryTimes int, commonArg interface{}) error {
	fmt.Println("retryTimes: ", retryTimes, commonArg)
	str, err := TryBestToString(arg)
	fmt.Println(str, err)
	return nil
}
