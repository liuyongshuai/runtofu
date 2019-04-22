/*
 * 各种类型的数据相互转来转去
 *
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-01-25 19:19
 */
package goUtils

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
)

//任意类型的数据转为结构体
func MakeElemType(d interface{}) ElemType {
	dval := reflect.ValueOf(d)
	for {
		if dval.Kind() == reflect.Ptr {
			dval = dval.Elem()
			continue
		}
		break
	}
	return ElemType{Data: d, RefVal: reflect.ValueOf(d)}
}

//基本的元素类型
type ElemType struct {
	Data   interface{}   //数据元素
	RefVal reflect.Value //通过反射获取的value
}

//转换为bool类型，如果是bool型则直接返回
func (et ElemType) ToBool() (bool, error) {
	return TryBestToBool(et.Data)
}

//转换为int类型
func (et ElemType) ToInt() (int, error) {
	tmp, err := et.ToInt64()
	if err != nil {
		return 0, err
	}
	return int(tmp), nil
}

//转换为int8类型
func (et ElemType) ToInt8() (int8, error) {
	tmp, err := et.ToInt64()
	if err != nil {
		return 0, err
	}
	if tmp > math.MaxInt8 || tmp < math.MinInt8 {
		return 0, fmt.Errorf("toInt8 failed, %v overflow [math.MinInt8,math.MaxInt8]", et.Data)
	}
	return int8(tmp), nil
}

//转换为int16类型
func (et ElemType) ToInt16() (int16, error) {
	tmp, err := et.ToInt64()
	if err != nil {
		return 0, err
	}
	if tmp > math.MaxInt16 || tmp < math.MinInt16 {
		return 0, fmt.Errorf("toInt16 failed, %v overflow [math.MinInt16,math.MaxInt16]", et.Data)
	}
	return int16(tmp), nil
}

//转换为int32类型
func (et ElemType) ToInt32() (int32, error) {
	tmp, err := et.ToInt64()
	if err != nil {
		return 0, err
	}
	if tmp > math.MaxInt32 || tmp < math.MinInt32 {
		return 0, fmt.Errorf("toInt32 failed, %v overflow [math.MinInt32,math.MaxInt32]", et.Data)
	}
	return int32(tmp), nil
}

//转换为int64类型
func (et ElemType) ToInt64() (int64, error) {
	return TryBestToInt64(et.Data)
}

//转换为uint类型
func (et ElemType) ToUint() (uint, error) {
	tmp, err := et.ToUint64()
	if err != nil {
		return 0, err
	}
	return uint(tmp), nil
}

//转换为uint8类型
func (et ElemType) ToUint8() (uint8, error) {
	tmp, err := et.ToUint64()
	if err != nil {
		return 0, err
	}
	if tmp > math.MaxUint8 {
		return 0, fmt.Errorf("toUint8 failed, %v overflow", et.Data)
	}
	return uint8(tmp), nil
}

//转换为uint16类型
func (et ElemType) ToUint16() (uint16, error) {
	tmp, err := et.ToUint64()
	if err != nil {
		return 0, err
	}
	if tmp > math.MaxUint16 {
		return 0, fmt.Errorf("toUint16 failed, %v overflow", et.Data)
	}
	return uint16(tmp), nil
}

//转换为uint32类型
func (et ElemType) ToUint32() (uint32, error) {
	tmp, err := et.ToUint64()
	if err != nil {
		return 0, err
	}
	if tmp > math.MaxUint32 {
		return 0, fmt.Errorf("toUint32 failed, %v overflow", et.Data)
	}
	return uint32(tmp), nil
}

//转换为uint64类型
func (et ElemType) ToUint64() (uint64, error) {
	return TryBestToUint64(et.Data)
}

//转换为string类型
func (et ElemType) ToString() string {
	str, _ := TryBestToString(et.Data)
	return str
}

//转换为float类型
func (et ElemType) ToFloat32() (float32, error) {
	f64, err := et.ToFloat64()
	if err != nil {
		return 0, err
	}
	if f64 > math.MaxFloat32 {
		return 0, fmt.Errorf("toFloat32 failed, %v overflow", et.Data)
	}
	return float32(f64), nil
}

//转换为float64类型
func (et ElemType) ToFloat64() (float64, error) {
	return TryBestToFloat(et.Data)
}

//转换为slice类型
//原始数据若为array/slice，则直接返回
//原始数据为map时只返回[]value
//原始数据若为数字、字符串等简单类型则将其放到slice中返回，即强制转为slice
//否则，报错
func (et ElemType) ToSlice() ([]ElemType, error) {
	switch et.RefVal.Kind() {
	case reflect.Slice, reflect.Array: //in为slice类型
		vlen := et.RefVal.Len()
		ret := make([]ElemType, vlen)
		for i := 0; i < vlen; i++ {
			ret[i] = MakeElemType(et.RefVal.Index(i).Interface())
		}
		return ret, nil
	case reflect.Map: //in为map类型，取map的value，要不要取map的key呢？
		var ret []ElemType
		ks := et.RefVal.MapKeys()
		for _, k := range ks {
			kiface := et.RefVal.MapIndex(k).Interface()
			ret = append(ret, MakeElemType(kiface))
		}
		return ret, nil
	case reflect.String: //字符串类型
		tmp := []byte(et.RefVal.String())
		var ret []ElemType
		for _, t := range tmp {
			ret = append(ret, MakeElemType(t))
		}
		return ret, nil
	default: //其他的类型一律强制转为slice
		return []ElemType{et}, nil
	}
}

//转换为map类型
//如果原始数据是map则直接返回
//如果是json字符串则尝试去解析
//否则，报错
func (et ElemType) ToMap() (map[ElemType]ElemType, error) {
	ret := make(map[ElemType]ElemType)
	if et.RefVal.Kind() == reflect.Map {
		ks := et.RefVal.MapKeys()
		for _, k := range ks {
			kiface := MakeElemType(k.Interface())
			viface := et.RefVal.MapIndex(k).Interface()
			ret[kiface] = MakeElemType(viface)
		}
		return ret, nil
	}
	if et.RefVal.Kind() == reflect.String {
		str := et.RefVal.String()
		var vmap interface{}
		err := json.Unmarshal([]byte(str), &vmap)
		if err != nil {
			return ret, err
		}
		inRefVal := reflect.ValueOf(vmap)
		if inRefVal.Kind() == reflect.Map {
			ks := inRefVal.MapKeys()
			for _, k := range ks {
				kiface := MakeElemType(k.Interface())
				viface := inRefVal.MapIndex(k).Interface()
				ret[kiface] = MakeElemType(viface)
			}
			return ret, nil
		}
	}
	return ret, fmt.Errorf("cannot convert %v to map", et.Data)
}

//提取原始数据的长度，只有string/slice/map/array/chan
func (et ElemType) Len() (int, error) {
	switch et.RefVal.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return et.RefVal.Len(), nil
	default:
		return 0, fmt.Errorf("invalid type for len %v", et.Data)
	}
}

//判断原始数据的类型是否为int
func (et ElemType) IsInt() bool {
	return et.Kind() == reflect.Int
}

//判断原始数据的类型是否为int8
func (et ElemType) IsInt8() bool {
	return et.Kind() == reflect.Int8
}

//判断原始数据的类型是否为int16
func (et ElemType) IsInt16() bool {
	return et.Kind() == reflect.Int16
}

//判断原始数据的类型是否为int32
func (et ElemType) IsInt32() bool {
	return et.Kind() == reflect.Int32
}

//判断原始数据的类型是否为int64
func (et ElemType) IsInt64() bool {
	return et.Kind() == reflect.Int64
}

//判断原始数据的类型是否为uint
func (et ElemType) IsUint() bool {
	return et.Kind() == reflect.Uint
}

//判断原始数据的类型是否为uint8
func (et ElemType) IsUint8() bool {
	return et.Kind() == reflect.Uint8
}

//判断原始数据的类型是否为uint16
func (et ElemType) IsUint16() bool {
	return et.Kind() == reflect.Uint16
}

//判断原始数据的类型是否为uint32
func (et ElemType) IsUint32() bool {
	return et.Kind() == reflect.Uint32
}

//判断原始数据的类型是否为uint64
func (et ElemType) IsUint64() bool {
	return et.Kind() == reflect.Uint64
}

//判断原始数据的类型是否为float32
func (et ElemType) IsFloat32() bool {
	return et.Kind() == reflect.Float32
}

//判断原始数据的类型是否为float64
func (et ElemType) IsFloat64() bool {
	return et.Kind() == reflect.Float64
}

//判断原始数据的类型是否为string
func (et ElemType) IsString() bool {
	return et.Kind() == reflect.String
}

//判断原始数据的类型是否为slice
func (et ElemType) IsSlice() bool {
	return et.Kind() == reflect.Slice
}

//判断原始数据的类型是否为map
func (et ElemType) IsMap() bool {
	return et.Kind() == reflect.Map
}

//判断原始数据的类型是否为array
func (et ElemType) IsArray() bool {
	return et.Kind() == reflect.Array
}

//判断原始数据的类型是否为chan
func (et ElemType) IsChan() bool {
	return et.Kind() == reflect.Chan
}

//判断原始数据的类型是否为bool
func (et ElemType) IsBool() bool {
	return et.Kind() == reflect.Bool
}

//是否为字符切片
func (et ElemType) IsByteSlice() bool {
	return reflect.TypeOf(et.Data).String() == "[]uint8"
}

//是否为简单类型：int/uint/string/bool/float....
func (et ElemType) IsSimpleType() bool {
	return et.IsInt() || et.IsInt8() || et.IsInt16() || et.IsInt32() || et.IsInt64() ||
		et.IsUint() || et.IsUint8() || et.IsUint16() || et.IsUint32() || et.IsUint64() ||
		et.IsString() || et.IsFloat32() || et.IsFloat64() || et.IsBool()
}

//是否为复合类型：slice/array/map/chan
func (et ElemType) IsComplexType() bool {
	return et.IsSlice() || et.IsMap() || et.IsArray() || et.IsChan()
}

//原始数据的类型
func (et ElemType) Kind() reflect.Kind {
	return et.RefVal.Kind()
}

//获取原始数据
func (et ElemType) RawData() interface{} {
	return et.Data
}
