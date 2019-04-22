/*
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     mysql
 * @date        2018-01-25 19:19
 */
package goUtils

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//存储MySQL的连接账号信息
type MySQLConf struct {
	Host            string        //连接地址
	Port            uint16        //端口号，默认3306
	User            string        //用户名
	Passwd          string        //密码
	DbName          string        //DB名称
	Charset         string        //设置的字符编码，默认utf8
	Timeout         time.Duration //连接超时时间，单位秒，默认5秒
	AutoCommit      bool          //是否自动提交，默认为true
	MaxIdleConns    int           //允许最大空闲连接数，默认为2
	MaxOpenConns    int           //最多允许打开多少个连接，默认0不限制
	ConnMaxLiftTime time.Duration //连接的最大生存时间，默认0不限制，单位秒
}

func MakeMySQLConf() MySQLConf {
	return MySQLConf{
		Port:            3306,
		Charset:         "utf8",
		Timeout:         5,
		AutoCommit:      true,
		MaxIdleConns:    2,
		MaxOpenConns:    0,
		ConnMaxLiftTime: 0,
	}
}

//设置连接地址
func (mc MySQLConf) SetHost(h string) MySQLConf {
	mc.Host = h
	return mc
}

//设置端口号
func (mc MySQLConf) SetPort(p uint16) MySQLConf {
	mc.Port = p
	return mc
}

//设置用户名
func (mc MySQLConf) SetUser(u string) MySQLConf {
	mc.User = u
	return mc
}

//设置密码
func (mc MySQLConf) SetPasswd(p string) MySQLConf {
	mc.Passwd = p
	return mc
}

//设置数据库名称
func (mc MySQLConf) SetDbName(d string) MySQLConf {
	mc.DbName = d
	return mc
}

//设置连接时的字符编码
func (mc MySQLConf) SetCharset(c string) MySQLConf {
	mc.Charset = c
	return mc
}

//设置连接时的超时时间
func (mc MySQLConf) SetTimeout(t time.Duration) MySQLConf {
	mc.Timeout = t
	return mc
}

//设置是否自动提交
func (mc MySQLConf) SetAutoCommit(b bool) MySQLConf {
	mc.AutoCommit = b
	return mc
}

//设置允许的最多多少个空闲连接
func (mc MySQLConf) SetMaxIdleConns(c int) MySQLConf {
	mc.MaxIdleConns = c
	return mc
}

//设置最多允许打开的连接数
func (mc MySQLConf) SetMaxOpenConns(c int) MySQLConf {
	mc.MaxOpenConns = c
	return mc
}

//设置每个连接最长的生存周期
func (mc MySQLConf) SetConnMaxLiftTime(t time.Duration) MySQLConf {
	mc.ConnMaxLiftTime = t
	return mc
}

var sqlTokenMap = make(map[string]string)
var delimiter = []string{"AND", "and", "OR", "or", ","}

func init() {
	sqlTokenMap["lt"] = "<"
	sqlTokenMap["gt"] = ">"
	sqlTokenMap["eq"] = "="
	sqlTokenMap["neq"] = "!="
	sqlTokenMap["lte"] = "<="
	sqlTokenMap["gte"] = ">="
	sqlTokenMap["in"] = "IN"
	sqlTokenMap["is"] = "IS"
	sqlTokenMap["notin"] = "NOT IN"
	sqlTokenMap["llike"] = "LIKE"
	sqlTokenMap["rlike"] = "LIKE"
	sqlTokenMap["like"] = "LIKE"
	sqlTokenMap["find"] = "FIND_IN_SET"
}

//检测分隔符是否合法
func checkDelimiter(delim string) bool {
	for _, a := range delimiter {
		if a == delim {
			return true
		}
	}
	return false
}

//条件连接符是否合法
func checkSQLToken(t string) bool {
	for k := range sqlTokenMap {
		if k == t {
			return true
		}
	}
	return false
}

/**
 * 仅用在预编译的查询语句中
 * 格式化要查询的条件语句，只支持简单语句
 */
func FormatCond(cond map[string]interface{}, delim string) (sqlCond string, param []interface{}) {
	if !checkDelimiter(delim) {
		return sqlCond, param
	}
	delim = strings.ToUpper(delim)
	var tmpCond []string
	var condItem = make(map[string]ElemType)
	for k, v := range cond {
		condItem[k] = MakeElemType(v)
	}

	//遍历所有的查询条件,k为字段名，v为相应的值
	for k, v := range condItem {
		tmpToken := "="
		tmpSym := "="
		key := k
		//如果字段名里包含":"
		if strings.Index(k, ":") > 0 {
			tmpS := strings.Split(k, ":")
			tmpSym = tmpS[1]
			key = tmpS[0]
			if checkSQLToken(tmpSym) {
				tmpToken = sqlTokenMap[tmpSym]
			}
		}
		//如果字段符号是in/notin,允许的值有string/slice/map
		vlen, verr := v.Len()
		if tmpSym == "in" || tmpSym == "notin" {
			var vslice []ElemType
			if v.IsString() && vlen > 0 { //如果是字符串则用"，"切成slice
				tmp := strings.Split(v.ToString(), ",")
				for _, t := range tmp {
					if len(t) <= 0 {
						continue
					}
					vslice = append(vslice, MakeElemType(t))
				}

			} else if v.IsSimpleType() { //其他的简单类型，直接填上去即可
				vslice = append(vslice, v)
			} else { //否则，对于复杂类型，只要可以转成slice就可以
				tos, toerr := v.ToSlice()
				if toerr == nil {
					vslice = append(vslice, tos...)
				}
			}
			vsliceLen := len(vslice)
			if vsliceLen <= 0 {
				continue
			}
			param = append(param, ConvertArgs(vslice)...)
			//用问号填充查询字段的值
			var tmpQ []string
			for i := 0; i < vsliceLen; i++ {
				tmpQ = append(tmpQ, "?")
			}
			//扔到最终的SQL条件里
			cd := fmt.Sprintf("`%s` %s (%s)", key, tmpToken, strings.Join(tmpQ, ","))
			tmpCond = append(tmpCond, cd)
			continue
		}
		//只允许简单类型
		if tmpToken == "LIKE" && v.IsSimpleType() {
			likeV := v.ToString()
			if len(likeV) <= 0 { //like的值不能为空，且要转为字符串
				continue
			}
			switch tmpSym {
			case "rlike":
				likeV += "%" //右like，加后面
			case "llike":
				likeV = "%" + likeV //左like，加前面
			case "like":
				likeV = "%" + likeV + "%" //双边like，前后都加
			}
			param = append(param, likeV)
			cd := fmt.Sprintf("`%s` %s ?", key, tmpToken)
			tmpCond = append(tmpCond, cd)
			continue
		}
		if tmpSym == "is" && vlen > 0 && verr == nil {
			cd := fmt.Sprintf("`%s` %s ?", key, tmpToken)
			tmpCond = append(tmpCond, cd)
			param = append(param, v.RawData())
			continue
		}
		//对于find_in_set来说，只允许简单类型
		if tmpSym == "find" && v.IsSimpleType() {
			if len(v.ToString()) <= 0 {
				continue
			}
			cd := fmt.Sprintf("%s(?,`%s`)", tmpToken, key)
			tmpCond = append(tmpCond, cd)
			param = append(param, v.RawData())
			continue
		}
		//其余的全部要求只能是简单类型
		if v.IsSimpleType() {
			cd := fmt.Sprintf("`%s` %s ?", key, tmpToken)
			tmpCond = append(tmpCond, cd)
			param = append(param, v.RawData())
			continue
		}
	}
	sqlCond = strings.Join(tmpCond, " "+delim+" ")
	return sqlCond, param
}

//转换查询SQL用的参数
func ConvertArgs(param []ElemType) []interface{} {
	args := make([]interface{}, len(param))
	for i := range param {
		args[i] = param[i].RawData()
	}
	return args
}

//将[]byte根据MySQL的字段类型转为相应的值
func convertMySQLType(b []byte, t *sql.ColumnType) interface{} {
	buf := bytes.NewBuffer(b)
	tstr := buf.String()
	scankind := t.ScanType().Kind().String()
	dbType := t.DatabaseTypeName()

	//根据接收类型做判断
	switch scankind {
	case "uint", "uint8", "uint16", "uint32", "uint64": //uint系列
		ret, err := strconv.ParseUint(tstr, 10, 64)
		if err != nil {
			return tstr
		}
		return ret
	case "int", "int8", "int16", "int32", "int64": //int系列
		ret, err := strconv.ParseInt(tstr, 10, 64)
		if err != nil {
			return tstr
		}
		return ret
	case "float32", "float64": //float系列
		ret, err := strconv.ParseFloat(tstr, 64)
		if err != nil {
			return tstr
		}
		return ret
	}

	//如果是varchar类型，可以返回字符串，对于text/blob/binary系列不好处理
	if dbType == "varchar" || dbType == "text" {
		return tstr
	}
	if scankind == "sql.RawBytes" {
		return b
	}
	return tstr
}

//提取字段
func filterTableFields(fields ...string) string {
	f := "*"
	if len(fields) > 0 {
		var tmp []string
		for _, tf := range fields {
			replace := strings.NewReplacer(
				"`", "",
				"'", "",
				"\"", "",
				",", "",
				"\\", "",
				"/", "",
			)
			tf = replace.Replace(tf)
			tmp = append(tmp, tf)
		}
		f = strings.Join(tmp, "`,`")
		f = "`" + f + "`"
	}
	return f
}
