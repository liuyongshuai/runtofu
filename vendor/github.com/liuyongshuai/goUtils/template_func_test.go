// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-29 16:16

package goUtils

import (
	"fmt"
	"html/template"
	"testing"
	"time"
)

func TestNewTplBuilder(t *testing.T) {
	testStart()

	d := CommonTplFuncs["date"].(func(timestamp int64, format string) template.HTML)
	fmt.Println(d(time.Now().Unix(), "Y/m/d H:i:s"))

	testEnd()
}
