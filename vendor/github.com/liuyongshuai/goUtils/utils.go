/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-03-27 15:20
 */
package goUtils

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	LocalHostIP         = ""
	LocalHostIpArr      []string
	LocalHostIpTraceId  = ""
	SequenceIDGenerator *SnowFlakeIdGenerator
	preTraceID          = ""
	colorFns            = []ColorFunc{Green, LightGreen, Cyan, LightCyan, Red, LightRed, Yellow, Black, DarkGray, LightGray, White, Blue, LightBlue, Purple, LightPurple, Brown}
)

func init() {
	//提取本机IP记录到日志里
	ips := LocalIP()
	for _, ip := range ips {
		if IsPrivateIP(ip) {
			LocalHostIpArr = append(LocalHostIpArr, ip)
			LocalHostIpTraceId = fmt.Sprintf("%s%x", LocalHostIpTraceId, Ip2long(ip))
		}
	}
	LocalHostIP = strings.Join(LocalHostIpArr, ",")
	//base62转码初始化
	for k, v := range base62CharToInt {
		base62IntToChar[v] = k
	}
	//ID生成器
	var err error
	SequenceIDGenerator, err = NewIDGenerator().
		SetTimeBitSize(40).
		SetSequenceBitSize(22).
		SetWorkerIdBitSize(1).
		SetWorkerId(1).Init()
	if err != nil {
		errMsg := fmt.Sprintf("Init IDGenerator failed, err=%v", err)
		BitchWarning(errMsg)
	}
}

// 在指定浮点数范围内生成随机数
func RandFloat64InRange(min, max float64) float64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

// 处理请求的时间proc_time
// 时间要求time.Now().UnixNano()
func ProcTime(st, et int64) float64 {
	var t1 int64 = 0
	var t2 int64 = 0
	if st > et {
		t1 = et
		t2 = st
	} else {
		t1 = st
		t2 = et
	}
	ret := float64(t2-t1) / float64(1000*1000*1000)
	return ret
}

//超级警告，下划线、闪烁提示
func BitchWarning(msg string) string {
	var tmp = []string{msg}
	for _, fn := range colorFns {
		tmp = append(tmp, fn(msg, 1, 1))
	}
	return strings.Join(tmp, "\n")
}

//超级警告，只有颜色
func FuckWarning(msg string) string {
	var tmp = []string{msg}
	for _, fn := range colorFns {
		tmp = append(tmp, fn(msg))
	}
	return strings.Join(tmp, "\n")
}

//生成一个假的traceId
func FakeTraceId() (traceId string) {
	for {
		ReRandSeed()
		traceId = fmt.Sprintf("%x%s%x", time.Now().UnixNano(), LocalHostIpTraceId, rand.Int63())
		if preTraceID != traceId {
			preTraceID = traceId
			break
		}
	}
	return traceId
}

//重新设置随机数种子
func ReRandSeed() {
	genId, err := SequenceIDGenerator.NextId()
	if err != nil {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(genId)
	}
}

//根据业务特点，过滤非法的ID并去重，一般用于批量根据ID提取信息时
func FilterIds(ids []interface{}) (ret []int64) {
	tmap := map[int64]struct{}{}
	for _, id := range ids {
		v, err := TryBestToInt64(id)
		if err != nil || v <= 0 {
			continue
		}
		tmap[v] = struct{}{}
	}
	for i := range tmap {
		ret = append(ret, i)
	}
	return
}

//返回最大的一个int型
func MaxInt64(args ...interface{}) (int64, error) {
	if len(args) <= 0 {
		return 0, ErrorInvalidInputType
	}
	var m int64 = math.MinInt64
	var tmps []int64
	for _, arg := range args {
		a, e := TryBestToInt64(arg)
		if e != nil {
			continue
		}
		tmps = append(tmps, a)
	}
	if len(tmps) <= 0 {
		return 0, ErrorInvalidInputType
	}
	for _, t := range tmps {
		if t > m {
			m = t
		}
	}
	return m, nil
}

//返回最小的一个int型
func MinInt64(args ...interface{}) (int64, error) {
	if len(args) <= 0 {
		return 0, ErrorInvalidInputType
	}
	var m int64 = math.MaxInt64
	var tmps []int64
	for _, arg := range args {
		a, e := TryBestToInt64(arg)
		if e != nil {
			continue
		}
		tmps = append(tmps, a)
	}
	if len(tmps) <= 0 {
		return 0, ErrorInvalidInputType
	}
	for _, t := range tmps {
		if t < m {
			m = t
		}
	}
	return m, nil
}

//获取当前终端的宽、高信息：字符数，非终端时（如IDE的执行环境）会报错
func GetTerminalSize() (width, height int, err error) {
	return terminal.GetSize(int(os.Stdout.Fd()))
}
