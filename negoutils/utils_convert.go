// 各种类型之间转来转去
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-22 18:27

package negoutils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type Basickind int

const (
	//转换时最大值
	MaxInt64Float  = float64(math.MaxInt64)
	MinInt64Float  = float64(math.MinInt64)
	MaxUint64Float = float64(math.MaxUint64)
	//基本类型归纳，类型转换时用得着
	InvalidKind Basickind = iota
	BoolKind
	ComplexKind
	IntKind
	FloatKind
	StringKind
	UintKind
	PtrKind
	ContainerKind
	FuncKind
)

var (
	//一些错误信息
	ErrorOverflowMaxInt64  = errors.New("this value overflow math.MaxInt64")
	ErrorOverflowMaxUint64 = errors.New("this value overflow math.MaxUint64")
	ErrorLessThanMinInt64  = errors.New("this value less than math.MinInt64")
	ErrorLessThanZero      = errors.New("this value less than zero")
	ErrorBadComparisonType = errors.New("invalid type for comparison")
	ErrorBadComparison     = errors.New("incompatible types for comparison")
	ErrorNoComparison      = errors.New("missing argument for comparison")
	ErrorInvalidInputType  = errors.New("invalid input type")
)

// 转换成特定类型，便于判断
func GetBasicKind(v reflect.Value) (Basickind, error) {
	switch v.Kind() {
	case reflect.Bool:
		return BoolKind, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return IntKind, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return UintKind, nil
	case reflect.Float32, reflect.Float64:
		return FloatKind, nil
	case reflect.Complex64, reflect.Complex128:
		return ComplexKind, nil
	case reflect.String:
		return StringKind, nil
	case reflect.Ptr:
		return PtrKind, nil
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		return ContainerKind, nil
	case reflect.Func:
		return FuncKind, nil
	}
	return InvalidKind, ErrorInvalidInputType
}

// int64转byte
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// bytes转int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// float64转byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, bits)
	return bs
}

// bytes转float
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

// 字符串转为字节切片
func StrToByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// 字节切片转为字符串
func ByteToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 尽最大努力将给定的类型转换为uint64
// 如"45.67"->45、"98.4abc3"->98、"34.87"->34
func TryBestToUint64(value interface{}) (uint64, error) {
	ret, err := tryBestConvertAnyTypeToInt(value, true)
	if err != nil {
		return 0, err
	}
	val := reflect.ValueOf(ret)
	if val.Kind() != reflect.Uint64 {
		return 0, ErrorInvalidInputType
	}
	return val.Uint(), nil
}

// 尽最大努力将给定的类型转换为int64
// 如"45.67"->45、"98.4abc3"->98、"34.87"->34
func TryBestToInt64(value interface{}) (int64, error) {
	ret, err := tryBestConvertAnyTypeToInt(value, false)
	if err != nil {
		return 0, err
	}
	val := reflect.ValueOf(ret)
	if val.Kind() != reflect.Int64 {
		return 0, ErrorInvalidInputType
	}
	return val.Int(), nil
}

// 从左边开始提取数据及小数点
func getFloatStrFromLeft(val string) string {
	val = strings.TrimSpace(val)
	valBytes := StrToByte(val)
	buf := bytes.Buffer{}
	for _, b := range valBytes {
		if b >= 48 && b <= 57 || b == 46 {
			buf.WriteByte(b)
			continue
		}
		break
	}
	return buf.String()
}

// 尽最大努力将任意类型转为int64或uint64
func tryBestConvertAnyTypeToInt(value interface{}, isUnsigned bool) (interface{}, error) {
	val := reflect.ValueOf(value)
	basicKind, err := GetBasicKind(val)
	if err != nil {
		return 0, err
	}
	switch basicKind {
	case IntKind:
		v := val.Int()
		if isUnsigned {
			if v >= 0 {
				return uint64(v), nil
			}
			return 0, ErrorLessThanZero
		}
		return v, nil
	case UintKind:
		v := val.Uint()
		if isUnsigned {
			return v, nil
		}
		if v > math.MaxInt64 {
			return 0, ErrorOverflowMaxInt64
		}
		return int(v), nil
	case StringKind: //取连续的最长的数字或小数点
		floatStr := getFloatStrFromLeft(val.String())
		if len(floatStr) <= 0 {
			if isUnsigned {
				return uint64(0), nil
			}
			return int64(0), nil
		}
		//先转成float，因为将"45.33"直接转为int/uint时会报错
		f, err := strconv.ParseFloat(floatStr, 10)
		if err != nil {
			return 0, err
		}
		return tryBestConvertAnyTypeToInt(f, isUnsigned)
		//float特殊处理，会有科学记数法表示形式
	case FloatKind:
		f := val.Float()
		if isUnsigned {
			if f > MaxUint64Float {
				return 0, ErrorOverflowMaxUint64
			}
			if f < 0 {
				return 0, ErrorLessThanZero
			}
			return uint64(f), nil
		}
		if f > MaxInt64Float {
			return 0, ErrorOverflowMaxInt64
		}
		if f < MinInt64Float {
			return 0, ErrorLessThanMinInt64
		}
		return int64(f), nil
	case BoolKind:
		b := val.Bool()
		tmp := 0
		if b {
			tmp = 1
		}
		if isUnsigned {
			return uint64(tmp), nil
		}
		return int64(tmp), nil
		//指针类型递归调用，直到取本值为止
	case PtrKind:
		if val.IsNil() {
			if isUnsigned {
				return uint64(0), nil
			}
			return int64(0), nil
		}
		return tryBestConvertAnyTypeToInt(val.Elem().Interface(), isUnsigned)
	default:
		return 0, ErrorInvalidInputType
	}
}

// 尽最大努力转换为字符串
func TryBestToString(value interface{}) (string, error) {
	val := reflect.ValueOf(value)
	basicKind, err := GetBasicKind(val)
	if err != nil {
		return "", err
	}
	switch basicKind {
	case IntKind:
		return strconv.FormatInt(val.Int(), 10), nil
	case UintKind:
		return strconv.FormatUint(val.Uint(), 10), nil
	case StringKind:
		return val.String(), nil
	case FloatKind:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
	case BoolKind:
		return strconv.FormatBool(val.Bool()), nil
	case PtrKind:
		if val.IsNil() {
			return "nil", nil
		}
		return TryBestToString(val.Elem().Interface())
	case ContainerKind:
		result, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(result), err
	default:
		return val.String(), nil
	}
}

// 尽最大努力转换为float64
func TryBestToFloat(value interface{}) (float64, error) {
	val := reflect.ValueOf(value)
	basicKind, err := GetBasicKind(val)
	if err != nil {
		return 0, err
	}
	switch basicKind {
	case IntKind:
		return float64(val.Int()), nil
	case UintKind:
		return float64(val.Uint()), nil
	case StringKind:
		floatStr := getFloatStrFromLeft(val.String())
		if len(floatStr) <= 0 {
			return 0, nil
		}
		return strconv.ParseFloat(floatStr, 10)
	case FloatKind:
		return val.Float(), nil
	case BoolKind:
		if val.Bool() {
			return 1, nil
		}
		return 0, nil
	case PtrKind:
		if val.IsNil() {
			return 0, nil
		}
		return TryBestToFloat(val.Elem().Interface())
	default:
		return 0, ErrorInvalidInputType
	}
}

// 尽最大努力转为bool类型
func TryBestToBool(value interface{}) (bool, error) {
	val := reflect.ValueOf(value)
	basicKind, err := GetBasicKind(val)
	if err != nil {
		return false, err
	}
	switch basicKind {
	case FloatKind:
		return val.Float() != 0, nil
	case IntKind:
		return val.Int() != 0, nil
	case UintKind:
		return val.Uint() != 0, nil
	case StringKind:
		v := strings.TrimSpace(val.String())
		if len(v) > 0 {
			return true, nil
		}
		return false, nil
	case BoolKind:
		return val.Bool(), nil
	case PtrKind:
		if val.IsNil() {
			return false, nil
		}
		return TryBestToBool(val.Elem().Interface())
	case FuncKind:
		return !val.IsNil(), nil
	}

	//对于Array, Chan, Map, Slice长度>0即可
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice, reflect.Map:
		return val.Len() != 0, nil
	}
	return false, ErrorInvalidInputType
}
