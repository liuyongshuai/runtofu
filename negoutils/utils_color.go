/*
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-01-25 19:19
 */
package negoutils

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type ColorType int

const (
	ColorType_Invalid ColorType = iota
	ColorType_Green
	ColorType_LightGreen
	ColorType_Cyan
	ColorType_LightCyan
	ColorType_Red
	ColorType_LightRed
	ColorType_Yellow
	ColorType_Black
	ColorType_DarkGray
	ColorType_LightGray
	ColorType_White
	ColorType_Blue
	ColorType_LightBlue
	ColorType_Purple
	ColorType_LightPurple
	ColorType_Brown
)

var (
	ColorFuncMap = map[ColorType]ColorFunc{
		ColorType_Green:       Green,
		ColorType_LightGreen:  LightGreen,
		ColorType_Cyan:        Cyan,
		ColorType_LightCyan:   LightCyan,
		ColorType_Red:         Red,
		ColorType_LightRed:    LightRed,
		ColorType_Yellow:      Yellow,
		ColorType_Black:       Black,
		ColorType_DarkGray:    DarkGray,
		ColorType_LightGray:   LightGray,
		ColorType_White:       White,
		ColorType_Blue:        Blue,
		ColorType_LightBlue:   LightBlue,
		ColorType_Purple:      Purple,
		ColorType_LightPurple: LightPurple,
		ColorType_Brown:       Brown,
	}
	colorFns = []ColorFunc{Green, LightGreen, Cyan, LightCyan, Red, LightRed, Yellow, Black, DarkGray, LightGray, White, Blue, LightBlue, Purple, LightPurple, Brown}
)

// 所有的颜色函数
type ColorFunc func(string, ...interface{}) string

// 根据颜色获取颜色函数
func GetColorFunc(t ColorType) (ColorFunc, bool) {
	if fn, ok := ColorFuncMap[t]; ok {
		return fn, ok
	}
	return nil, false
}

// 绿色字体，modifier里，第一个控制闪烁，第二个控制下划线
func Green(str string, modifier ...interface{}) string {
	return CliColorRender(str, 32, 0, modifier...)
}

// 淡绿
func LightGreen(str string, modifier ...interface{}) string {
	return CliColorRender(str, 32, 1, modifier...)
}

// 青色/蓝绿色
func Cyan(str string, modifier ...interface{}) string {
	return CliColorRender(str, 36, 0, modifier...)
}

// 淡青色
func LightCyan(str string, modifier ...interface{}) string {
	return CliColorRender(str, 36, 1, modifier...)
}

// 红字体
func Red(str string, modifier ...interface{}) string {
	return CliColorRender(str, 31, 0, modifier...)
}

// 淡红色
func LightRed(str string, modifier ...interface{}) string {
	return CliColorRender(str, 31, 1, modifier...)
}

// 黄色字体
func Yellow(str string, modifier ...interface{}) string {
	return CliColorRender(str, 33, 0, modifier...)
}

// 黑色
func Black(str string, modifier ...interface{}) string {
	return CliColorRender(str, 30, 0, modifier...)
}

// 深灰色
func DarkGray(str string, modifier ...interface{}) string {
	return CliColorRender(str, 30, 1, modifier...)
}

// 浅灰色
func LightGray(str string, modifier ...interface{}) string {
	return CliColorRender(str, 37, 0, modifier...)
}

// 白色
func White(str string, modifier ...interface{}) string {
	return CliColorRender(str, 37, 1, modifier...)
}

// 蓝色
func Blue(str string, modifier ...interface{}) string {
	return CliColorRender(str, 34, 0, modifier...)
}

// 淡蓝
func LightBlue(str string, modifier ...interface{}) string {
	return CliColorRender(str, 34, 1, modifier...)
}

// 紫色
func Purple(str string, modifier ...interface{}) string {
	return CliColorRender(str, 35, 0, modifier...)
}

// 淡紫色
func LightPurple(str string, modifier ...interface{}) string {
	return CliColorRender(str, 35, 1, modifier...)
}

// 棕色
func Brown(str string, modifier ...interface{}) string {
	return CliColorRender(str, 33, 0, modifier...)
}

func CliColorRender(str string, color int, weight int, extraArgs ...interface{}) string {
	//闪烁效果
	var isBlink int64 = 0
	if len(extraArgs) > 0 {
		isBlink = reflect.ValueOf(extraArgs[0]).Int()
	}
	//下划线效果
	var isUnderLine int64 = 0
	if len(extraArgs) > 1 {
		isUnderLine = reflect.ValueOf(extraArgs[1]).Int()
	}
	var mo []string
	if isBlink > 0 {
		mo = append(mo, "05")
	}
	if isUnderLine > 0 {
		mo = append(mo, "04")
	}
	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	buf := bytes.Buffer{}
	buf.WriteString("\033[")
	buf.WriteString(strings.Join(mo, ";"))
	buf.WriteString(";")
	buf.WriteString(fmt.Sprintf("%d", color))
	buf.WriteString("m")
	buf.WriteString(str)
	buf.WriteString("\033[0m")
	//fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), color)
	return buf.String()
}
