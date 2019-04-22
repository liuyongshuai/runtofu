// 简单的进度条操作
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-12-07 14:52

package goUtils

import (
	"fmt"
	"os"
	"strings"
)

//构造一个进度条
func NewProgressBar() *ProgressBar {
	ret := &ProgressBar{
		finishChar:  '=',
		unreachChar: '-',
	}
	w, _, e := GetTerminalSize()
	if e != nil {
		w = 200
	} else {
		w = int(float64(w) * 0.85)
	}
	ret.barWidth = float64(w)
	return ret
}

//进度条结构体
type ProgressBar struct {
	finishChar  byte    //已完成的进度条显示的字符
	unreachChar byte    //未完成的进度条显示字符
	total       float64 //总数量
	finish      float64 //已完成的数量
	barWidth    float64 //进度条宽度
}

//设置已完成的进度条
func (pb *ProgressBar) SetFinishChar(c byte) *ProgressBar {
	pb.finishChar = c
	return pb
}

//设置未完成的部分显示的字符
func (pb *ProgressBar) SetUnFinishChar(c byte) *ProgressBar {
	pb.unreachChar = c
	return pb
}

//设置总数
func (pb *ProgressBar) SetTotalNum(t float64) *ProgressBar {
	pb.total = t
	return pb
}

//设置进度条总宽度
func (pb *ProgressBar) SetBarWidth(w float64) *ProgressBar {
	pb.barWidth = w
	return pb
}

//设置已完成数量
func (pb *ProgressBar) SetFinishNum(t float64) *ProgressBar {
	pb.finish = t
	return pb
}

//显示完了,可选的方法
func (pb *ProgressBar) ForceFinish() {
	pb.finish = pb.total
	pb.Render()
}

//在终端渲染进度条
func (pb *ProgressBar) Render() {
	delim := ">"
	if pb.finish >= pb.total {
		pb.finish = pb.total
		delim = ""
	}
	ratio := pb.finish / pb.total
	percent := ratio * 100
	finishNum := ratio * pb.barWidth
	unreachNum := pb.barWidth - finishNum
	h1 := Green(strings.Repeat(string(pb.finishChar), int(finishNum)))
	h2 := Red(strings.Repeat(string(pb.unreachChar), int(unreachNum)))
	percentStr := fmt.Sprintf("%.02f%%", percent)
	fmt.Fprintf(os.Stdout, "\r[%s%s%s]%s", h1, delim, h2, Yellow(percentStr))
	os.Stdout.Sync()
}
