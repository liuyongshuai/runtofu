// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-27 19:02

package goUtils

import (
	"bytes"
	"fmt"
	"strings"
)

type TerminalTable struct {
	headerData     []string   //表头数据
	rowData        [][]string //行的数据
	useSeparator   bool       //是否需要每行间的分隔线
	columnNum      int        //列的数量，以最多的一行的列为准
	maxColumnWidth []int      //每列的最大宽度，对齐用的
	separatorLine  string     //行分隔符
}

func NewTerminalTable() *TerminalTable {
	t := &TerminalTable{
		useSeparator: true,
	}
	return t
}

//是否使用行的分隔符
func (t *TerminalTable) IsUseRowSeparator(b bool) *TerminalTable {
	t.useSeparator = b
	return t
}

//添加表头数据
func (t *TerminalTable) SetHeader(header []string) *TerminalTable {
	if len(header) > t.columnNum {
		t.columnNum = len(header)
	}
	for _, h := range header {
		t.headerData = append(t.headerData, h)
	}
	return t
}

//添加一行数据
func (t *TerminalTable) AddRow(row []string) *TerminalTable {
	t.rowData = append(t.rowData, row)
	if len(row) > t.columnNum {
		t.columnNum = len(row)
	}
	return t
}

//添加许多行数据
func (t *TerminalTable) AddRows(rows [][]string) *TerminalTable {
	for _, row := range rows {
		if len(row) > t.columnNum {
			t.columnNum = len(row)
		}
		t.rowData = append(t.rowData, row)
	}
	return t
}

//计算属性
func (t *TerminalTable) prepareSomething() {
	if len(t.headerData) <= 0 && len(t.rowData) <= 0 {
		return
	}
	var rows [][]string
	t.maxColumnWidth = make([]int, t.columnNum)
	headerLen := len(t.headerData)
	//把所有的行的列数补齐到一致，方便输出
	if headerLen > 0 {
		if headerLen < t.columnNum {
			for i := headerLen; i < t.columnNum; i++ {
				t.headerData = append(t.headerData, "")
			}
		}
	}
	for idx, row := range t.rowData {
		if len(row) < t.columnNum {
			for i := len(row); i < t.columnNum; i++ {
				t.rowData[idx] = append(t.rowData[idx], "")
			}
		}
	}
	//计算每行最大的列宽
	for cellIdx, cellStr := range t.headerData {
		cellStr = strings.TrimSpace(cellStr)
		cellStr = fmt.Sprintf(" %s ", cellStr)
		t.maxColumnWidth[cellIdx] = RuneStringWidth(cellStr)
	}
	if headerLen > 0 {
		rows = append(rows, t.headerData)
	}
	for _, row := range t.rowData {
		for cellIdx, cellStr := range row {
			//每个表格单元的数据，两边加一个空格，显的好看
			cellStr = strings.TrimSpace(cellStr)
			cellStr = fmt.Sprintf(" %s ", cellStr)
			//再计算每列的最大宽度
			l := RuneStringWidth(cellStr)
			if t.maxColumnWidth[cellIdx] < l {
				t.maxColumnWidth[cellIdx] = l
			}
			row[cellIdx] = cellStr
		}
		rows = append(rows, row)
	}
	//行分隔符，根据每列的最大列宽来决定
	buf := bytes.Buffer{}
	buf.WriteString("+")
	for _, w := range t.maxColumnWidth {
		buf.WriteString(strings.Repeat("-", w))
		buf.WriteString("+")
	}
	t.separatorLine = buf.String()
	//将每列的数据补整齐：填充到相同的宽度，给表头加黄色
	for rowIdx, row := range rows {
		for cellIdx, cellStr := range row {
			cellStr = RuneFillRight(cellStr, t.maxColumnWidth[cellIdx])
			//表头数据标黄
			if headerLen > 0 && rowIdx == 0 {
				cellStr = Yellow(cellStr)
			}
			rows[rowIdx][cellIdx] = cellStr
		}
	}
	t.rowData = rows
}

//开始返回表格数据
func (t *TerminalTable) Render() string {
	headerLen := len(t.headerData)
	rowLen := len(t.rowData)
	if headerLen <= 0 && rowLen <= 0 {
		return ""
	}
	t.prepareSomething()
	buf := bytes.Buffer{}

	buf.WriteString(t.separatorLine)
	buf.WriteString("\n")
	for idx, row := range t.rowData {
		buf.WriteString("|" + strings.Join(row, "|") + "|")
		buf.WriteString("\n")
		if t.useSeparator {
			buf.WriteString(t.separatorLine)
			buf.WriteString("\n")
		} else if idx == 0 && headerLen > 0 && rowLen > 0 {
			buf.WriteString(t.separatorLine)
			buf.WriteString("\n")
		}
	}
	if !t.useSeparator {
		buf.WriteString(t.separatorLine)
		buf.WriteString("\n")
	}

	return buf.String()
}
