// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        utils_time.go
// @date        2025-03-25 10:20

package negoutils

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"
)

var (
	TimeLayoutDate     = "2006-01-02"
	TimeLayoutDateTime = "2006-01-02 15:04:05"
)

func CurTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func Timestamp2Str(sec, nsec int64, format string) string {
	dt := time.Unix(sec, nsec)
	if format == "" {
		format = TimeLayoutDateTime
	}
	return dt.Format(format)
}

func Str2Time(dtStr, format string) (dt time.Time, err error) {
	var TimeFormat string
	if format == "" {
		format = TimeLayoutDateTime
	}
	dt, err = time.ParseInLocation(TimeFormat, dtStr, time.Local)
	return
}

func Str2Timestamp(dtStr, format string) (ts int64, err error) {
	var TimeFormat string
	if format == "" {
		format = TimeLayoutDateTime
	}
	dt, err := time.ParseInLocation(TimeFormat, dtStr, time.Local)
	if err != nil {
		return 0, err
	}
	return dt.Unix(), nil
}

func MonthStart() time.Time {
	y, m, _ := time.Now().Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
}

func TodayStart() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func TodayEnd() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 23, 59, 59, 1e9-1, time.Local)
}

// 将字符串日期转为时间戳，秒
// t要转化的时间戳，如"2018-04-09"、"2018-08-04 12:34"
func StrToTimestamp(t string) int64 {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	reg, _ := regexp.Compile(`[\-:\s]+`)
	tmp := reg.Split(t, -1)
	if len(tmp) <= 0 {
		return 0
	}
	if len(tmp) > 6 {
		tmp = tmp[0:6]
	}
	for k, v := range tmp {
		vv, _ := MakeElemType(v).ToInt()
		if k == 0 {
			tmp[k] = fmt.Sprintf("%04d", vv)
		} else {
			tmp[k] = fmt.Sprintf("%02d", vv)
		}
	}
	tl := []string{"", "01", "01", "00", "00", "00"}
	if len(tmp) < 6 {
		tmp = append(tmp, tl[len(tmp):]...)
	}
	tstr := fmt.Sprintf("%s-%s-%s %s:%s:%s", tmp[0], tmp[1], tmp[2], tmp[3], tmp[4], tmp[5])
	utm, err := time.ParseInLocation("2006-01-02 15:04:05", tstr, loc)
	if err != nil {
		return 0
	}
	return utm.Unix()
}

// 格式化时间
// t为时间戳，秒数
// format为要格式化的格式，如"Y-m-d H:i:s"
func FormatCTime(t interface{}, format string) template.HTML {
	tm, err := MakeElemType(t).ToInt64()
	if err != nil {
		return template.HTML("")
	}
	local, err2 := time.LoadLocation("Asia/Chongqing")
	if err2 != nil {
		fmt.Println(err2)
	}
	replacer := strings.NewReplacer(
		"Y", "2006",
		"m", "01",
		"d", "02",
		"H", "15",
		"i", "04",
		"s", "05",
	)
	format = replacer.Replace(format)
	ret := time.Unix(tm, 0).In(local).Format(format)
	return template.HTML(ret)
}
