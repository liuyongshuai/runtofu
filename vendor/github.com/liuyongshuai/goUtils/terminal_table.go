// 支持超长行，会对其进行拆行处理
//
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-27 19:02

package goUtils

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
)

func NewTerminalTable() *TerminalTable {
	t := &TerminalTable{
		headerFontColorFunc: Yellow,
		rowFontColorFunc:    nil,
		borderColorFunc:     nil,
	}
	return t
}

type TerminalTable struct {
	//表格的原始数据，目前只有表头和行内容
	rawHeaderData []string   //原始的表头数据
	rawRowData    [][]string //原始的行的数据

	//控制字体颜色
	headerFontColorFunc ColorFunc //表头的字体颜色，默认Yellow
	rowFontColorFunc    ColorFunc //表格内容的字体颜色
	borderColorFunc     ColorFunc //边框的颜色

	//以下是根据添加的数据自动生成的
	maxColumnNum   int                 //列的数量，以最多的一行的列为准
	maxColumnWidth []int               //每列的最大宽度，对齐用的
	rowData        []*terminalTableRow //所有行

	//是根据表格的列数、屏幕的宽度来自动生成的
	allTableAllowWidth int //表格允许的最大宽度
}

//添加表头数据
func (t *TerminalTable) SetHeader(header []string) *TerminalTable {
	if len(header) > t.maxColumnNum {
		t.maxColumnNum = len(header)
	}
	for _, h := range header {
		t.rawHeaderData = append(t.rawHeaderData, h)
	}
	return t
}

//添加表头字体颜色
func (t *TerminalTable) SetHeaderFontColor(color ColorType) *TerminalTable {
	colorFunc, ok := GetColorFunc(color)
	if ok {
		t.headerFontColorFunc = colorFunc
	}
	return t
}

//添加一行数据
func (t *TerminalTable) AddRow(row []string) *TerminalTable {
	t.rawRowData = append(t.rawRowData, row)
	if len(row) > t.maxColumnNum {
		t.maxColumnNum = len(row)
	}
	return t
}

//添加许多行数据
func (t *TerminalTable) AddRows(rows [][]string) *TerminalTable {
	for _, row := range rows {
		if len(row) > t.maxColumnNum {
			t.maxColumnNum = len(row)
		}
		t.rawRowData = append(t.rawRowData, row)
	}
	return t
}

//添加行内容字体颜色
func (t *TerminalTable) SetRowFontColor(color ColorType) *TerminalTable {
	colorFunc, ok := GetColorFunc(color)
	if ok {
		t.rowFontColorFunc = colorFunc
	}
	return t
}

//添加边框颜色
func (t *TerminalTable) SetBorderFontColor(color ColorType) *TerminalTable {
	colorFunc, ok := GetColorFunc(color)
	if ok {
		t.borderColorFunc = colorFunc
	}
	return t
}

//开始返回表格数据
func (t *TerminalTable) Render() string {
	headerLen := len(t.rawHeaderData)
	rowLen := len(t.rawRowData)
	if headerLen <= 0 && rowLen <= 0 {
		return ""
	}
	t.prepareSomething()

	//行分隔符，根据每列的最大列宽来决定
	sepBuf := bytes.Buffer{}
	joinStr := t.borderStr("+")
	sepBuf.WriteString(joinStr)
	for _, w := range t.maxColumnWidth {
		sepBuf.WriteString(t.borderStr(strings.Repeat("-", w)))
		sepBuf.WriteString(joinStr)
	}
	horizontalLine := sepBuf.String()

	dataBuf := bytes.Buffer{}

	//第一个行分隔符，先写进去
	dataBuf.WriteString(horizontalLine)
	dataBuf.WriteString("\n")
	for idx := range t.rowData {
		row := t.rowData[idx]
		rowStr := t.renderSingleRow(row)
		if rowStr == "" {
			continue
		}
		dataBuf.WriteString(rowStr)
		dataBuf.WriteString(horizontalLine)
		dataBuf.WriteString("\n")
	}

	return dataBuf.String()
}

type rowType byte

const (
	rowTypeHeader rowType = iota
	rowTypeData
)

//一行数据，包含多个列
type terminalTableRow struct {
	lineNum  int     //本行各个小格子的数据行数，本行所有列的行数都是一样的
	rowType  rowType //本行数据的类型，是表头还是数据内容
	cellList []*terminalTableCell
}

//一个小格子里的数据
type terminalTableCell struct {
	columnNo    int      //第几列
	maxWidth    int      //宽度
	cellStrList []string //本小格子的数据，可能要分多行
}

var (
	splitReg, _ = regexp.Compile(`\n`)
	replaceList = strings.NewReplacer(
		"“", "\"",
		"”", "\"",
		"\t", "    ",
		"…", "...",
	)
)

//计算属性
func (t *TerminalTable) prepareSomething() {
	if len(t.rawHeaderData) <= 0 && len(t.rawRowData) <= 0 {
		return
	}
	//每一行中都有一些多余的字符，要将屏幕宽度减去这部分
	t.getMaxColumnWidths()
	headerLen := len(t.rawHeaderData)

	//把所有的行的列数补齐到一致，方便输出
	if headerLen > 0 {
		if headerLen < t.maxColumnNum {
			for i := headerLen; i < t.maxColumnNum; i++ {
				t.rawHeaderData = append(t.rawHeaderData, " ")
			}
		}
	}
	for idx, row := range t.rawRowData {
		if len(row) < t.maxColumnNum {
			for i := len(row); i < t.maxColumnNum; i++ {
				t.rawRowData[idx] = append(t.rawRowData[idx], " ")
			}
		}
	}

	//统一给各行折一下行
	if headerLen > 0 {
		tmp := t.wrapTableRows(t.rawHeaderData)
		if tmp != nil {
			tmp.rowType = rowTypeHeader
			t.rowData = append(t.rowData, tmp)
		}
	}
	for idx := range t.rawRowData {
		tmp := t.wrapTableRows(t.rawRowData[idx])
		if tmp != nil {
			tmp.rowType = rowTypeData
			t.rowData = append(t.rowData, tmp)
		}
	}

	//将每列的数据补整齐
	for rowIdx := range t.rowData {
		row := t.rowData[rowIdx]
		for cellIdx, cellUnit := range row.cellList {
			maxWidth := t.maxColumnWidth[cellIdx]
			for subCellIdx, subCellStr := range cellUnit.cellStrList {
				subCellStr = RuneFillRight(subCellStr, maxWidth)
				cellUnit.cellStrList[subCellIdx] = subCellStr
			}
		}
	}
}

//计算各列的最大长度
func (t *TerminalTable) getMaxColumnWidths() {
	t.allTableAllowWidth = ScreenWidth - t.maxColumnNum*4
	t.maxColumnWidth = make([]int, t.maxColumnNum)

	//算上表头，计算各列宽度最大值
	for idx, row := range t.rawHeaderData {
		t.maxColumnWidth[idx] = RuneStringWidth(row)
	}
	for _, rows := range t.rawRowData {
		for idx, row := range rows {
			l := RuneStringWidth(row)
			if l > t.maxColumnWidth[idx] {
				t.maxColumnWidth[idx] = l
			}
		}
	}

	//再次校验一下总宽度
	allWidth := int(0)
	var maxColWidthList []*terminalTableCell
	for idx, w := range t.maxColumnWidth {
		allWidth += w
		maxColWidthList = append(maxColWidthList, &terminalTableCell{columnNo: idx, maxWidth: w})
	}

	//如果各小格子的宽度和大于屏幕宽度，每轮循环都将最宽的列折行，一直折到适合屏幕大小为止
	if allWidth > t.allTableAllowWidth {
		for {
			diff := allWidth - t.allTableAllowWidth
			sort.Slice(maxColWidthList, func(i, j int) bool {
				return maxColWidthList[i].maxWidth > maxColWidthList[j].maxWidth
			})
			//本次总宽度消除量
			reduce := maxColWidthList[0].maxWidth / 3
			if reduce > diff {
				reduce = diff
			}
			maxColWidthList[0].maxWidth -= reduce
			allWidth -= reduce
			if allWidth <= t.allTableAllowWidth {
				break
			}
		}
	}

	//重新设置各列宽度最大值
	for _, colW := range maxColWidthList {
		t.maxColumnWidth[colW.columnNo] = colW.maxWidth + 2
	}
}

//表头字体获取
func (t *TerminalTable) headerStr(str string) string {
	if t.headerFontColorFunc != nil {
		return t.headerFontColorFunc(str)
	}
	return Yellow(str)
}

//边框字体获取
func (t *TerminalTable) borderStr(str string) string {
	if t.borderColorFunc != nil {
		return t.borderColorFunc(str)
	}
	return str
}

//行内容字体获取
func (t *TerminalTable) rowStr(str string) string {
	if t.rowFontColorFunc != nil {
		return t.rowFontColorFunc(str)
	}
	return str
}

//生成一行数据的格式
func (t *TerminalTable) renderSingleRow(row *terminalTableRow) string {
	if row == nil || len(row.cellList) <= 0 {
		return ""
	}
	buf := bytes.Buffer{}

	//竖线分隔符
	verSepLine := t.borderStr("|")

	//列的数量
	colNum := len(row.cellList)

	//小格子里的内容被拆分成了多少行，所有小格子的行都是一样的
	srowNum := len(row.cellList[0].cellStrList)
	for i := 0; i < srowNum; i++ {
		for j := 0; j < colNum; j++ {
			buf.WriteString(verSepLine)
			str := row.cellList[j].cellStrList[i]
			switch row.rowType {
			case rowTypeHeader:
				str = t.headerStr(str)
			case rowTypeData:
				str = t.rowStr(str)
			}
			buf.WriteString(str)
		}
		//最后一个分隔符
		buf.WriteString(verSepLine)
		buf.WriteString("\n")
	}
	return buf.String()
}

//将一行数据折行，并返回最大行数，主要策略是每次都将最长的行折半，一直折到所有行的长度小于屏幕长度
func (t *TerminalTable) wrapTableRows(rawRow []string) (retRow *terminalTableRow) {
	if len(rawRow) <= 0 {
		return nil
	}
	retRow = &terminalTableRow{}
	cellNoMap := make(map[int]*terminalTableCell)

	//统计各小格子宽度
	for idx, row := range rawRow {
		rawRow[idx] = replaceList.Replace(row)
		cell := &terminalTableCell{
			columnNo: idx,
			maxWidth: t.maxColumnWidth[idx],
		}
		cellNoMap[idx] = cell
	}

	//各小格子的最大行数
	maxLineNum := 0

	//开始对各小格子进行拆行处理
	for idx, cellStr := range rawRow {
		cellUnit := cellNoMap[idx]
		cellWidth := cellUnit.maxWidth - 2
		lineNum := 1
		if RuneStringWidth(cellStr) > cellWidth {
			cellStr, lineNum = RuneWrap(cellStr, cellWidth)
		}
		if lineNum > maxLineNum {
			maxLineNum = lineNum
		}
		tmp := splitReg.Split(cellStr, -1)
		for _, t := range tmp {
			t = strings.Trim(t, "\n")
			cellUnit.cellStrList = append(cellUnit.cellStrList, " "+t+" ")
		}
	}

	//如果每行数据不够最大的行数，用空行补齐
	for idx := range rawRow {
		cellUnit := cellNoMap[idx]
		tmpLen := len(cellUnit.cellStrList)
		for i := 0; i < maxLineNum-tmpLen; i++ {
			cellUnit.cellStrList = append(cellUnit.cellStrList, " ")
		}
		retRow.cellList = append(retRow.cellList, cellUnit)
	}

	return retRow
}
