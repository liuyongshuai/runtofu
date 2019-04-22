// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-28 20:52

package goUtils

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	//如果字段名以数字开头，替换一下
	intToWordMap = []string{
		"zero",
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
	}
	//常用的缩写，一般全大写
	commonInitialisms = map[string]bool{
		"API":   true,
		"ASCII": true,
		"CPU":   true,
		"CSS":   true,
		"DNS":   true,
		"EOF":   true,
		"GUID":  true,
		"HTML":  true,
		"HTTP":  true,
		"HTTPS": true,
		"ID":    true,
		"IP":    true,
		"JSON":  true,
		"LHS":   true,
		"QPS":   true,
		"RAM":   true,
		"RHS":   true,
		"RPC":   true,
		"SLA":   true,
		"SMTP":  true,
		"SSH":   true,
		"TLS":   true,
		"TTL":   true,
		"UI":    true,
		"UID":   true,
		"UUID":  true,
		"URI":   true,
		"URL":   true,
		"UTF8":  true,
		"VM":    true,
		"XML":   true,
	}
)

//生成mysql里所有表的go结构体形式
func GetMySQLAllTablesStruct(db *DBase) (ret string, err error) {
	buf := bytes.Buffer{}
	tables, err := GetAllMySQLTables(db)
	if err != nil {
		return
	}
	for _, table := range tables {
		str, err := GetMySQLTableStruct(db, table)
		if err != nil {
			return ret, err
		}
		buf.WriteString(str)
	}
	return buf.String(), nil
}

//提取所有的表
func GetAllMySQLTables(db *DBase) (ret []string, err error) {
	querySQL := fmt.Sprintf("SHOW TABLES")
	tmpRows, err := db.FetchRows(querySQL)
	if err != nil || len(tmpRows) <= 0 {
		return
	}
	if len(tmpRows) <= 0 {
		return
	}
	for _, table := range tmpRows {
		for _, v := range table {
			ret = append(ret, v.ToString())
		}
	}
	return
}

//mysql表信息
type mysqlTableInfo struct {
	TableName    string           //表名
	TableComment string           //表注释
	Fields       []mysqlFieldInfo //所有的字段列表
}

//mysql表的字段结构体，用于自动生成相应的go结构体用的
type mysqlFieldInfo struct {
	FieldName    string //字段名称
	DataType     string //数据类型
	IsUnsigned   bool   //是否为无符号类型
	FieldComment string //字段的注释
}

/**
CREATE TABLE `admin_menu` (
   `menu_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
   `menu_name` varchar(100) NOT NULL DEFAULT '' COMMENT '菜单名称',
   `menu_path` varchar(100) NOT NULL DEFAULT '' COMMENT '菜单路径',
   `icon_name` varchar(100) NOT NULL DEFAULT '' COMMENT '图标名称',
   `icon_color` varchar(100) NOT NULL DEFAULT '' COMMENT '图标的颜色',
   `parent_menu_id` int(11) NOT NULL DEFAULT '0' COMMENT '父菜单ID',
   `child_menu_num` int(11) NOT NULL DEFAULT '0' COMMENT '子菜单数量',
   `menu_sort` int(11) NOT NULL DEFAULT '0' COMMENT '同级菜单的排序值',
   PRIMARY KEY (`menu_id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8 COMMENT='菜单表'
*/
//获取表的结构
func GetMySQLTableStruct(db *DBase, tableName string) (ret string, err error) {
	var tableInfo mysqlTableInfo
	querySQL := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	tmpRows, err := db.FetchRows(querySQL)
	if err != nil || len(tmpRows) <= 0 {
		return
	}
	ct := tmpRows[0]["Create Table"]
	createSQL := ct.ToString()
	createSQL = strings.Replace(createSQL, "\n", "", -1)
	regexp1, _ := regexp.Compile(`CREATE TABLE.*?\((.*)\).*?COMMENT='(.*?)'`)
	tmpStrs := regexp1.FindStringSubmatch(createSQL)
	if len(tmpStrs) <= 2 {
		return
	}
	fieldStr := tmpStrs[1]
	tableInfo.TableComment = tmpStrs[2]
	tableInfo.TableName = tableName
	splitReg, _ := regexp.Compile(`,\s*`)
	fieldReg, _ := regexp.Compile("`(.*?)`\\s+(.*?)[(\\s].*")
	commentReg, _ := regexp.Compile(`.*?(?i:COMMENT\s+'(.*?)')`)
	unsignedReg, _ := regexp.Compile(`.*?(?i:unsigned(.*?))`)
	fieldStrList := splitReg.Split(fieldStr, -1)
	for _, fieldInfo := range fieldStrList {
		var f mysqlFieldInfo
		fieldAttr := fieldReg.FindStringSubmatch(strings.TrimSpace(fieldInfo))
		if len(fieldAttr) <= 2 {
			continue
		}
		f.FieldName = fieldAttr[1]
		f.DataType = fieldAttr[2]
		c := commentReg.FindStringSubmatch(fieldInfo)
		if len(c) > 1 {
			f.FieldComment = c[1]
		}
		isUnsigned := unsignedReg.FindStringSubmatch(fieldInfo)
		if len(isUnsigned) > 1 {
			f.IsUnsigned = true
		}
		tableInfo.Fields = append(tableInfo.Fields, f)
	}

	//开始生成go结构体
	buf := bytes.Buffer{}
	if len(tableInfo.TableComment) <= 0 {
		tableInfo.TableComment = fmt.Sprintf("table %s", tableName)
	}
	buf.WriteString(fmt.Sprintf("//%s\n", tableInfo.TableComment))
	buf.WriteString(fmt.Sprintf("type %s struct {\n", FormatFieldNameToGolangType(tableName)))
	for _, fieldInfo := range tableInfo.Fields {
		mysqlFieldName := fieldInfo.FieldName
		goFieldName := FormatFieldNameToGolangType(mysqlFieldName)
		goType := mysqlTypeToGoType(fieldInfo.DataType, fieldInfo.IsUnsigned)
		mysqlComment := fieldInfo.FieldComment
		if len(mysqlComment) <= 0 {
			mysqlComment = mysqlFieldName
		}
		buf.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" db:\"%s\"`//%s\n", goFieldName, goType, mysqlFieldName, mysqlFieldName, mysqlComment))
	}
	buf.WriteString("}\n\n")
	return buf.String(), nil
}

//格式字段的名称
func FormatFieldNameToGolangType(fieleName string) string {
	//如果首字符为数字则要转换一下
	first := fieleName[:1]
	firstChar, err := strconv.ParseInt(first, 10, 8)
	if err == nil {
		fieleName = intToWordMap[firstChar] + "_" + fieleName[1:]
	}
	if fieleName == "_" {
		return fieleName
	}
	//去掉前面的下划线
	for len(fieleName) > 0 && fieleName[0] == '_' {
		fieleName = fieleName[1:]
	}
	//判断是否全为小写
	allLower := true
	for _, r := range fieleName {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	runes := []rune(fieleName)
	//如果全为小写的话，判断是不是为常用的缩写
	if allLower {
		if u := strings.ToUpper(fieleName); commonInitialisms[u] {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}
	w, i := 0, 0
	for i+1 <= len(runes) {
		eow := false
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}
			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			eow = true
		}
		i++
		if !eow {
			continue
		}
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			copy(runes[w:], []rune(u))
		} else if strings.ToLower(word) == word {
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	//再处理一次
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	return string(runes)
}

//mysql类型转为go类型
func mysqlTypeToGoType(mysqlType string, isUnsigned bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if isUnsigned {
			return "uint"
		}
		return "int"
	case "bigint":
		if isUnsigned {
			return "uint64"
		}
		return "int64"
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
		return "string"
	case "date", "datetime", "time", "timestamp":
		return "string"
	case "decimal", "double":
		return "float64"
	case "float":
		return "float32"
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return "[]byte"
	}
	return "interface{}"
}
