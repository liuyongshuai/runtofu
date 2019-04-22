// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-12-07 19:43

package goUtils

import (
	"fmt"
	"runtime"
	"strings"
)

//开始测试
func testStart() {
	s, _ := getTestDelim()
	fmt.Print(s)
}

func testEnd() {
	_, e := getTestDelim()
	fmt.Print(e)
}

//内部使用的获取测试方法的输出分隔符
func getTestDelim() (start, end string) {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	msg := f.Name()
	w, _, e := GetTerminalSize()
	if e != nil {
		start = fmt.Sprintf("\n\n\n========START %s========\n", Yellow(msg))
		end = fmt.Sprintf("========END %s========\n", Yellow(msg))
		return
	}
	w -= 10
	d := w - len(msg)
	if d <= 0 {
		d = 6
	}
	start = fmt.Sprintf("\n\n\n%sSTART %s%s\n", strings.Repeat("=", d/2), Yellow(msg), strings.Repeat("=", d/2))
	end = fmt.Sprintf("%sEND   %s%s\n", strings.Repeat("=", d/2), Yellow(msg), strings.Repeat("=", d/2))
	return
}
