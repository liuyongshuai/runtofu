/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-03-27 15:20
 */
package negoutils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// 从字符串里提取单字，只要中文汉字
func ExtractCNWord(str string) (ret []string) {
	ret = make([]string, 0)
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			ret = append(ret, string(r))
		}
	}
	return
}

// 对比两个string切片，看内容是否一样
func CompareStringSlice(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, item := range a {
		if item != b[i] {
			return false
		}
	}
	return true
}

// 是否全是汉字
func IsAllChinese(str string) bool {
	if len(str) <= 0 {
		return false
	}
	ret := true
	for _, r := range str {
		if !unicode.Is(unicode.Scripts["Han"], r) {
			ret = false
			break
		}
	}
	return ret
}

// 是否为汉字、字母、数字
func IsNormalStr(str string) bool {
	if len(str) <= 0 {
		return false
	}
	reg := regexp.MustCompile("^[a-zA-Z0-9\u4e00-\u9fa5]+$")
	return reg.MatchString(str)
}

// 半角字符转全角字符【处理搜索的query用】
func ToDBC(str string) string {
	ret := ""
	for _, r := range str {
		if r == 32 {
			ret += string(rune(12288))
		} else if r < 127 {
			ret += string(r + 65248)
		} else {
			ret += string(r)
		}
	}
	return ret
}

// 全角字符转半角字符【处理搜索的query用】
func ToCBD(str string) string {
	ret := ""
	for _, r := range str {
		if r == 12288 {
			ret += string(r - 12256)
			continue
		}
		if r > 65280 && r < 65375 {
			ret += string(r - 65248)
		} else {
			ret += string(r)
		}
	}
	return ret
}

// md5转换
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

// 打乱一个字符串slice
func StrSliceShuffle(slice []string) []string {
	sl := len(slice)
	if sl <= 0 {
		return slice
	}
	for i := 0; i < sl; i++ {
		ReRandSeed()
		a := rand.Intn(sl)
		ReRandSeed()
		b := rand.Intn(sl)
		slice[a], slice[b] = slice[b], slice[a]
	}
	return slice
}

// 截取字符串
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 {
		start += length
	} else if start > length {
		start = start % length
	}
	if end < 0 {
		end += length
	} else if end > length {
		end = end % length
	}
	if start > end || end < 0 || start < 0 {
		return ""
	}
	return string(rs[start : end+1])
}

// 去重（不保证原顺序）
func UniqueStrSlice(slice []string) []string {
	tmp := map[string]bool{}
	for _, v := range slice {
		tmp[v] = true
	}
	var ret []string
	for t := range tmp {
		ret = append(ret, t)
	}
	return ret
}

// 检查是否在slice里面
func InStrSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// 字符串hash为uint64
func StrHashSum64(str string) uint64 {
	return xxhash.Sum64(StrToByte(str))
}

var (
	alphaNumAll = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)
	alphaNum1   = []byte(`0123456789abcdefghijklmnopqrstuvwxyz`)
	alphaNum2   = []byte(`abcdefghijklmnopqrstuvwxyz`)
	alphaNum3   = []byte(`0123456789`)
)

// 生成一堆随机数
func RandomStr(n int, alphabets ...byte) string {
	if len(alphabets) == 0 {
		alphabets = alphaNumAll
	}
	var byteSlice = make([]byte, n)
	var randBy bool
	ReRandSeed()
	if num, err := rand.Read(byteSlice); num != n || err != nil {
		randBy = true
	}
	for i, b := range byteSlice {
		if randBy {
			ReRandSeed()
			byteSlice[i] = alphabets[rand.Intn(len(alphabets))]
		} else {
			byteSlice[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return string(byteSlice)
}

var base62CharToInt = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
var base62IntToChar = make(map[string]int)

// base62转换
func Base62Encode(num int64) string {
	baseStr := ""
	for {
		if num <= 0 {
			break
		}
		i := num % 62
		baseStr += base62CharToInt[i]
		num = (num - i) / 62
	}
	return baseStr
}

// base62解码
func Base62Decode(b62Str string) int64 {
	var rs int64 = 0
	for i := 0; i < len(b62Str); i++ {
		rs += int64(base62IntToChar[string(b62Str[i])]) * int64(math.Pow(62, float64(i)))
	}
	return rs
}

// 高仿PHP的 preg_replace_callback
// pattern：正则表达式
// originStr：要处理的字符串
// fn：参数是字符串切片，0表示整个匹配的字符串，1表示正则里的第一个捕获项、2表示第二个、依次类推。。。。
func PregReplaceCallback(pattern, originStr string, fn func([]string) string) (string, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return originStr, err
	}
	originLen := len(originStr)
	if originLen <= 0 {
		return originStr, err
	}
	buf := bytes.Buffer{}
	var endIndex, preIndex int
	//提取出所有的子字符串
	subList := reg.FindAllStringSubmatchIndex(originStr, -1)
	if len(subList) <= 0 {
		return originStr, nil
	}
	for _, subInfo := range subList {
		//subs的结构类似：[9,12,23,45]，表示匹配项的起止位置，必须是偶数个
		subLen := len(subInfo)
		//这玩意吧，啥玩意也没匹配上
		if subLen%2 != 0 || subLen < 2 {
			return originStr, fmt.Errorf("invalid sub match")
		}
		var matches []string
		//提取所有的匹配项
		for i := 0; i < subLen; i += 2 {
			si := subInfo[i]
			ei := subInfo[i+1]
			matches = append(matches, originStr[si:ei])
		}
		startIndex := subInfo[0]
		endIndex = subInfo[1]
		//上一次循环的结束位置：本次的开始位置的字符填充到结果里
		buf.WriteString(originStr[preIndex:startIndex])
		buf.WriteString(fn(matches))
		preIndex = endIndex
	}
	//如果没到字符串的末尾，全部填充进去即可
	if endIndex < originLen {
		buf.WriteString(originStr[endIndex:])
	}
	if buf.Len() <= 0 {
		return originStr, nil
	}
	return buf.String(), nil
}

// 判断字符串是否全为数字
func IsAllNumber(str string) bool {
	if len(str) <= 0 {
		return false
	}
	if len(strings.Trim(str, "0123456789")) > 0 {
		return false
	}
	return true
}

// 解析字符串，高仿PHP的http://php.net/manual/zh/function.parse-str.php
// "first=value&arr[]=foo+bar&arr[]=baz";
func ParseStr(str string) (ret map[string][]ElemType) {
	ret = make(map[string][]ElemType)
	tmpRet := map[string]*[]ElemType{}
	fArr := strings.Split(str, "&")
	if len(fArr) <= 0 {
		return
	}
	for _, arg := range fArr {
		//可以包含等号，也可以不包含
		if !strings.Contains(arg, "=") {
			tmpRet[arg] = &[]ElemType{MakeElemType("")}
			continue
		}
		//截取第一个鹄等号前的
		field := Substr(arg, 0, strings.Index(arg, "=")-1)
		val := Substr(arg, strings.Index(arg, "=")+1, -1)
		//如果字段名里有[]
		if !strings.HasSuffix(field, "[]") {
			tmpRet[field] = &[]ElemType{MakeElemType(val)}
			continue
		}
		realFieldName := strings.TrimRight(field, "[]")
		if tmpArr, ok := tmpRet[realFieldName]; ok {
			*tmpArr = append(*tmpArr, MakeElemType(val))
			continue
		}
		tmpRet[realFieldName] = &[]ElemType{MakeElemType(val)}
	}
	//腾挪一下结果
	for k, vals := range tmpRet {
		var tmp []ElemType
		for _, val := range *vals {
			tmp = append(tmp, val)
		}
		ret[k] = tmp
	}
	return
}

// 打乱一个字符串
func StrShuffle(str string) string {
	rs := []rune(str)
	sliceLen := len(rs)
	for i := 0; i < sliceLen; i++ {
		ReRandSeed()
		a := rand.Intn(sliceLen)
		ReRandSeed()
		b := rand.Intn(sliceLen)
		rs[a], rs[b] = rs[b], rs[a]
	}
	return string(rs)
}
