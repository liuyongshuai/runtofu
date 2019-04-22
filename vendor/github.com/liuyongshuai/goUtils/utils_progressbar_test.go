// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-12-07 16:24

package goUtils

import (
	"fmt"
	"testing"
	"time"
)

func TestNewProgressBar(t *testing.T) {
	testStart()
	bar := NewProgressBar()
	bar.SetTotalNum(1111)
	for i := 0; i < 1111; i += 23 {
		time.Sleep(100 * time.Millisecond)
		bar.SetFinishNum(float64(i))
		bar.Render()
	}
	bar.ForceFinish()
	fmt.Println()
	testEnd()
}
