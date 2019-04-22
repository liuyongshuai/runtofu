// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2019-03-02 15:26

package goUtils

import (
	"fmt"
	"testing"
)

func TestRuneWrap(t *testing.T) {
	str := "擘画强军蓝图，指引奋进征程。2013年到2018年全国两会期间，中共中央总书记、国家主席、中央军委主席习近平连续出席解放军和武警部队代表团全体会议并发表重要讲话，提出一系列新思想、新论断、新要求。6年来，全军部队认真贯彻习主席重要讲话精神，牢固确立习近平强军思想的指导地位，重振政治纲纪，重塑组织形态，重整斗争格局，重构建设布局，重树作风形象，在中国特色强军之路上迈出坚定步伐。"
	ret, lineNum := RuneWrap(str, ScreenWidth)
	fmt.Println(ret, ScreenWidth, lineNum)
	str = "擘画强军蓝图，\n指引奋进征程。"
	ret, lineNum = RuneWrap(str, 10)
	fmt.Println(ret, ScreenWidth, lineNum)
	str = ""
	ret, lineNum = RuneWrap(str, ScreenWidth)
	fmt.Println(ret, ScreenWidth, lineNum)
}
