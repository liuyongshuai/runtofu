// 对map的常用操作
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-12-05 13:37

package goUtils

type mapReduceFunc func(interface{}, interface{}) (interface{}, interface{}, bool)

//提取map的key
func MapKeys(m map[interface{}]interface{}) []interface{} {
	keys := make([]interface{}, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

//提取map的value
func MapValues(m map[interface{}]interface{}) []interface{} {
	values := make([]interface{}, 0)
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

//是否包含某个值
func MapIsSet(m map[interface{}]interface{}, key interface{}) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

//是否为空map
func MapIsEmpty(m map[interface{}]interface{}) bool {
	if (m == nil) || len(m) == 0 {
		return true
	}
	return false
}

//合并两个map，后面的覆盖前面相同的key
func MapMerge(args ...map[interface{}]interface{}) map[interface{}]interface{} {
	newmap := make(map[interface{}]interface{})
	for _, value := range args {
		if value == nil {
			continue
		}
		for sk, sv := range value {
			newmap[sk] = sv
		}
	}
	return newmap
}

//遍历map
func MapIterator(m map[interface{}]interface{}, fn mapReduceFunc) map[interface{}]interface{} {
	ret := make(map[interface{}]interface{})
	if m != nil {
		for k, v := range m {
			k1, v1, ok := fn(k, v)
			if !ok {
				continue
			}
			ret[k1] = v1
		}
	}
	return ret
}
