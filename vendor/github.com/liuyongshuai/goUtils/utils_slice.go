// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-22 18:47

package goUtils

import (
	"math/rand"
	"reflect"
)

type reduceCallbackFunc func(interface{}) interface{}
type filterCallbackFunc func(interface{}) bool

//将其他类型转为interface的切片，支持slice/map/[u]int[8-64]/float../string
func ToSliceIface(in interface{}) []interface{} {
	vin := reflect.ValueOf(in)
	switch vin.Kind() {
	case reflect.Slice, reflect.Array: //in为slice类型
		vlen := vin.Len()
		ret := make([]interface{}, vlen)
		for i := 0; i < vlen; i++ {
			ret[i] = vin.Index(i).Interface()
		}
		return ret
	case reflect.Map: //in为map类型
		ks := vin.MapKeys()
		vlen := vin.Len()
		ret := make([]interface{}, vlen)
		for _, k := range ks {
			ret = append(ret, vin.MapIndex(k).Interface())
		}
		return ret
	case reflect.String: //字符串类型
		tmp := []byte(vin.String())
		var ret []interface{}
		for _, t := range tmp {
			ret = append(ret, t)
		}
		return ret
	default:
		return []interface{}{vin.Interface()}
	}
}

//检查interface类型是否在slice里
func InSlice(val interface{}, sl []interface{}) bool {
	for _, sval := range sl {
		if sval == val {
			return true
		}
	}
	return false
}

//合并两个slice
func SliceMerge(slice1, slice2 []interface{}) (c []interface{}) {
	c = append(slice1, slice2...)
	return
}

//对给定的slice调用回调函数，生成一个新的slice
func SliceReduce(slice []interface{}, a reduceCallbackFunc) (dSlice []interface{}) {
	for _, v := range slice {
		dSlice = append(dSlice, a(v))
	}
	return
}

//从给定的slice里随机提取新的slice
func SliceRand(a []interface{}) (b interface{}) {
	ReRandSeed()
	randnum := rand.Intn(len(a))
	b = a[randnum]
	return
}

//计算int64型slice的和
func SliceSum(slice []interface{}) (sum int64) {
	for _, v := range slice {
		vi, _ := TryBestToInt64(v)
		sum += vi
	}
	return
}

//返回给定的slice在回调函数返回true的新slice
func SliceFilter(slice []interface{}, a filterCallbackFunc) (filteredSlice []interface{}) {
	for _, v := range slice {
		if a(v) {
			filteredSlice = append(filteredSlice, v)
		}
	}
	return
}

//计算两个slice的差集
func SliceDiff(slice1, slice2 []interface{}) (diffSlice []interface{}) {
	for _, v := range slice1 {
		if !InSlice(v, slice2) {
			diffSlice = append(diffSlice, v)
		}
	}
	return
}

//计算slice交集
func SliceIntersect(slice1, slice2 []interface{}) (intersectSlice []interface{}) {
	for _, v := range slice1 {
		if InSlice(v, slice2) {
			intersectSlice = append(intersectSlice, v)
		}
	}
	return
}

//将一个slice切成若干个大小的子slice
func SliceChunk(slice []interface{}, size int) (chunkSlice [][]interface{}) {
	if size >= len(slice) {
		chunkSlice = append(chunkSlice, slice)
		return
	}
	end := size
	for i := 0; i <= (len(slice) - size); i += size {
		chunkSlice = append(chunkSlice, slice[i:end])
		end += size
	}
	return
}

//填充slice
func SlicePad(slice []interface{}, size int, val interface{}) []interface{} {
	if size <= len(slice) {
		return slice
	}
	for i := 0; i < (size - len(slice)); i++ {
		slice = append(slice, val)
	}
	return slice
}

//slice去重
func SliceUnique(slice []interface{}) (uniqueSlice []interface{}) {
	for _, v := range slice {
		if !InSlice(v, uniqueSlice) {
			uniqueSlice = append(uniqueSlice, v)
		}
	}
	return
}

//打乱一个slice
func SliceShuffle(slice []interface{}) []interface{} {
	ReRandSeed()
	sliceLen := len(slice)
	for i := 0; i < sliceLen; i++ {
		a := rand.Intn(sliceLen)
		b := rand.Intn(sliceLen)
		slice[a], slice[b] = slice[b], slice[a]
	}
	return slice
}

//转为字符串slice
func ToStringSlice(arg []interface{}, ignoreErr bool) (ret []string, err error) {
	for _, v := range arg {
		vStr, err := TryBestToString(v)
		if err != nil {
			if ignoreErr {
				continue
			}
			return ret, err
		}
		ret = append(ret, vStr)
	}
	return
}

//转为int64地slice
func ToInt64Slice(arg []interface{}, ignoreErr bool) (ret []int64, err error) {
	for _, v := range arg {
		vInt64, err := TryBestToInt64(v)
		if err != nil {
			if ignoreErr {
				continue
			}
			return ret, err
		}
		ret = append(ret, vInt64)
	}
	return
}

//转为float64地slice
func ToFloat64Slice(arg []interface{}, ignoreErr bool) (ret []float64, err error) {
	for _, v := range arg {
		vFloat64, err := TryBestToFloat(v)
		if err != nil {
			if ignoreErr {
				continue
			}
			return ret, err
		}
		ret = append(ret, vFloat64)
	}
	return
}
