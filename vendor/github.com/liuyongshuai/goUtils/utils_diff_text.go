// 两段文本按行对比
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-22 18:43

package goUtils

import (
	"fmt"
	"regexp"
	"strings"
)

//打印两段文本的不同之处，逐行对比哟，必须以换行符
func PrintTextDiff(text1, text2 string) {
	reg := regexp.MustCompile(`(\n|\r|\n\r)`)
	t1 := reg.Split(strings.TrimSpace(text1), -1)
	t2 := reg.Split(strings.TrimSpace(text2), -1)
	PrintTextDiffByGroup([][]string{t1}, [][]string{t2})
}

//分组打印文本差别
func PrintTextDiffByGroup(leftText, rightText [][]string) {
	rightTextLen := len(rightText)
	leftTextLen := len(leftText)
	var maxLeftRowLen, maxRightRowLen int
	//左则最长的一行
	for _, rows := range leftText {
		for _, row := range rows {
			rowLen := RuneStringWidth(row)
			if rowLen > maxLeftRowLen {
				maxLeftRowLen = rowLen
			}
		}
	}
	//右则最长的一行
	for _, rows := range rightText {
		for _, row := range rows {
			rowLen := RuneStringWidth(row)
			if rowLen > maxRightRowLen {
				maxRightRowLen = rowLen
			}
		}
	}
	if maxLeftRowLen <= 0 {
		maxLeftRowLen = 8
	}
	if maxRightRowLen <= 0 {
		maxRightRowLen = 8
	}
	//分隔符
	delim := strings.Repeat("-", maxRightRowLen+maxLeftRowLen)
	//开始打印
	for leftIdx1, leftRows1 := range leftText {
		fmt.Println(delim)
		//如果右边没了，继续打印左边
		if leftIdx1 >= rightTextLen {
			for _, row := range leftRows1 {
				fmt.Println(fmt.Sprintf("%s |", Red(RuneFillRight(row, maxLeftRowLen))))
			}
			continue
		}
		rightRows1 := rightText[leftIdx1]
		//将两行填充到数量相等
		leftLen1 := len(leftRows1)
		rightLen1 := len(rightRows1)
		if leftLen1 > rightLen1 {
			d := leftLen1 - rightLen1
			for i := 0; i < d; i++ {
				rightRows1 = append(rightRows1, "")
			}
		} else {
			d := rightLen1 - leftLen1
			for i := 0; i < d; i++ {
				leftRows1 = append(leftRows1, "")
			}
		}
		for leftIdx2, leftRow2 := range leftRows1 {
			rightRow2 := rightRows1[leftIdx2]
			if rightRow2 != leftRow2 {
				rightRow2 = Red(rightRow2)
				leftRow2 = Red(RuneFillRight(leftRow2, maxLeftRowLen))
			} else {
				leftRow2 = RuneFillRight(leftRow2, maxLeftRowLen)
			}
			fmt.Println(fmt.Sprintf("%s | %s", leftRow2, rightRow2))
		}
	}
	//如果右边还有的话
	if rightTextLen > leftTextLen {
		for i := leftTextLen; i < rightTextLen; i++ {
			fmt.Println(delim)
			for _, row := range rightText[i] {
				fmt.Println(fmt.Sprintf("%s | %s", strings.Repeat(" ", maxLeftRowLen), Red(row)))
			}
		}
	}
	fmt.Println(delim)
}
