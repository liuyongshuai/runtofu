/**
 * 有关经纬度的geohash编码、解码、寻找周围格子的方法
 *
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-04-17 13:28
 */
package negoutils

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	GEOHASH_BASE32              = "0123456789bcdefghjkmnpqrstuvwxyz"
	MAX_LATITUDE        float64 = 90
	MIN_LATITUDE        float64 = -90
	MAX_LONGITUDE       float64 = 180
	MIN_LONGITUDE       float64 = -180
	GEOHASH_TYPE_NORMAL         = 1 //普通的切格子方法，如"wttex9"
	GEOHASH_TYPE_BITS           = 2 //切出来的结果是整数
)

var (
	geoHashBits      = []int{16, 8, 4, 2, 1}
	geoHashBitsLen   = len(geoHashBits)
	geoHashBase32    = []byte(GEOHASH_BASE32)
	geoHashBase32Pos = make(map[byte]int)
)

// 初始化操作
func init() {
	for i, c := range geoHashBase32 {
		geoHashBase32Pos[c] = i
	}
}

// 统一封装的计算geohash的方法，
// 当 gtype==GEOHASH_TYPE_NORMAL 时的精度情况见GeoHashEncode()方法的说明
// 当 gtype==GEOHASH_TYPE_BITS 时的精度情况见GeoHashBitsEncode()方法的说明
func GeneralGeoHashEncode(lat, lng float64, precision, gtype int) (string, *GeoRectangle) {
	if gtype == GEOHASH_TYPE_NORMAL {
		return GeoHashEncode(lat, lng, precision)
	} else if gtype == GEOHASH_TYPE_BITS {
		geo, rect := GeoHashBitsEncode(lat, lng, uint8(precision))
		return fmt.Sprintf("%d", geo), rect
	}
	return "", &GeoRectangle{}
}

// 统一封装的解码geohash方法
// precision精度参数只有当gtype == GEOHASH_TYPE_BITS时才用得着，其他的传0即可
func GeneralGeoHashDecode(geohash string, precision, gtype int) *GeoRectangle {
	if gtype == GEOHASH_TYPE_NORMAL {
		return GeoHashDecode(geohash)
	} else if gtype == GEOHASH_TYPE_BITS {
		geo, err := strconv.ParseUint(geohash, 10, 64)
		if err != nil {
			return &GeoRectangle{}
		}
		return GeoHashBitsDecode(geo, uint8(precision))
	}
	return &GeoRectangle{}
}

/*
*
geoHash算法详见介绍：https://www.cnblogs.com/LBSer/p/3310455.html
最后生成一个geohash的字符串，和一个包括geohash的小格子，可以取到小格子四个角经纬度信息及宽高
北京的相关精度如下（不同纬度的宽度不一样，从南向北，小格子越来越窄）：

	precision:4	dist:30022m x 19567m
	precision:5	dist:3750m x 4891m
	precision:6	dist:937m x 611m
	precision:7	dist:117m x 152m
	precision:8	dist:29m x 19m
	precision:9	dist:3m x 4m
*/
func GeoHashEncode(lat, lng float64, precision int) (string, *GeoRectangle) {
	if lat < MIN_LATITUDE || lat > MAX_LATITUDE {
		return "", &GeoRectangle{}
	}
	if lng < MIN_LONGITUDE || lng > MAX_LONGITUDE {
		return "", &GeoRectangle{}
	}
	var buf bytes.Buffer
	var minLat, maxLat = MIN_LATITUDE, MAX_LATITUDE
	var minLng, maxLng = MIN_LONGITUDE, MAX_LONGITUDE
	var mid float64 = 0

	bit, ch, isEven := 0, 0, true
	for i := 0; i < precision; {
		if isEven { //偶数位放经度
			if mid = (minLng + maxLng) / 2; mid < lng {
				ch |= geoHashBits[bit]
				minLng = mid
			} else {
				maxLng = mid
			}
		} else { //奇数位处理纬度
			if mid = (minLat + maxLat) / 2; mid < lat {
				ch |= geoHashBits[bit]
				minLat = mid
			} else {
				maxLat = mid
			}
		}
		isEven = !isEven
		if bit < geoHashBitsLen-1 {
			bit++
		} else {
			buf.WriteByte(geoHashBase32[ch])
			bit, ch = 0, 0
			i++
		}
	}

	b := &GeoRectangle{MinLat: minLat, MaxLat: maxLat, MinLng: minLng, MaxLng: maxLng}
	return buf.String(), b
}

// geoHash转换为坐标点，返回的是一个格子
// 当对一对经纬度转为geoHash，再次转为经纬度时只能保证是同一个格子
// 得到的只是格子四个角的经纬度坐标，并不能精确的还原之前的经纬度
func GeoHashDecode(geohash string) *GeoRectangle {
	ret := &GeoRectangle{}
	isEven := true
	lats := [2]float64{MIN_LATITUDE, MAX_LATITUDE}
	lngs := [2]float64{MIN_LONGITUDE, MAX_LONGITUDE}
	for _, ch := range geohash {
		pos, ok := geoHashBase32Pos[byte(ch)]
		if !ok {
			return ret
		}
		for i := 0; i < geoHashBitsLen; i++ {
			ti := pos & geoHashBits[i]
			if ti > 0 {
				ti = 0
			} else {
				ti = 1
			}
			if isEven {
				lngs[ti] = (lngs[0] + lngs[1]) / 2.0
			} else {
				lats[ti] = (lats[0] + lats[1]) / 2.0
			}
			isEven = !isEven
		}
	}
	ret.MinLat = lats[0]
	ret.MaxLat = lats[1]
	ret.MinLng = lngs[0]
	ret.MaxLng = lngs[1]
	return ret
}

// 计算给定的经纬度点在指定精度下周围8个区域的geoHash编码，包括自身，一共9个点
func GetNeighborsGeoCodes(lat, lng float64, precision int) []string {
	if lat < MIN_LATITUDE || lat > MAX_LATITUDE {
		return []string{}
	}
	if lng < MIN_LONGITUDE || lng > MAX_LONGITUDE {
		return []string{}
	}
	geoHashList := make([]string, 9)

	//自身的区域
	cur, b := GeoHashEncode(lat, lng, precision)
	geoHashList[0] = cur

	//上下左右四个格子
	centerUp, _ := GeoHashEncode((b.MinLat+b.MaxLat)/2+b.LatSpan(), (b.MinLng+b.MaxLng)/2, precision)
	centerBottom, _ := GeoHashEncode((b.MinLat+b.MaxLat)/2-b.LatSpan(), (b.MinLng+b.MaxLng)/2, precision)
	leftCenter, _ := GeoHashEncode((b.MinLat+b.MaxLat)/2, (b.MinLng+b.MaxLng)/2-b.LngSpan(), precision)
	rightCenter, _ := GeoHashEncode((b.MinLat+b.MaxLat)/2, (b.MinLng+b.MaxLng)/2+b.LngSpan(), precision)

	//四个角的格子
	leftUp, _ := GeoHashEncode((b.MinLat+b.MaxLat)/2+b.LatSpan(), (b.MinLng+b.MaxLng)/2-b.LngSpan(), precision)
	leftBottom, _ := GeoHashEncode((b.MinLat+b.MaxLat)/2-b.LatSpan(), (b.MinLng+b.MaxLng)/2-b.LngSpan(), precision)
	rightUp, _ := GeoHashEncode((b.MinLat+b.MaxLat)/2+b.LatSpan(), (b.MinLng+b.MaxLng)/2+b.LngSpan(), precision)
	rightBottom, _ := GeoHashEncode((b.MinLat+b.MaxLat)/2-b.LatSpan(), (b.MinLng+b.MaxLng)/2+b.LngSpan(), precision)

	//八个格子赋值
	geoHashList[1] = centerUp
	geoHashList[2] = centerBottom
	geoHashList[3] = leftCenter
	geoHashList[4] = rightCenter
	geoHashList[5] = leftUp
	geoHashList[6] = leftBottom
	geoHashList[7] = rightUp
	geoHashList[8] = rightBottom
	return geoHashList
}

/*
*
将geohash编码为uint64类型，源自：https://github.com/yinqiwen/geohash-int
北京的相关精度如下（不同纬度的宽度不一样，从南向北，小格子越来越窄）：

	precision=12	dist=7500 x 4891
	precision=13	dist=3750 x 2445
	precision=14	dist=1875 x 1222
	precision=15	dist=937 x 611
	precision=16	dist=468 x 305
	precision=17	dist=234 x 152
	precision=18	dist=117 x 76
	precision=19	dist=58 x 38
	precision=20	dist=29 x 19
	precision=21	dist=14 x 9
*/
func GeoHashBitsEncode(lat, lng float64, precision uint8) (geo uint64, rect *GeoRectangle) {
	if lat < MIN_LATITUDE || lat > MAX_LATITUDE {
		rect = &GeoRectangle{}
		return
	}
	if lng < MIN_LONGITUDE || lng > MAX_LONGITUDE {
		rect = &GeoRectangle{}
		return
	}
	if precision < 1 || precision > 32 {
		rect = &GeoRectangle{}
		return
	}
	var minLat, maxLat = MIN_LATITUDE, MAX_LATITUDE
	var minLng, maxLng = MIN_LONGITUDE, MAX_LONGITUDE
	var i uint8
	for i = 0; i < precision; i++ {
		var latBit, lngBit uint64
		if maxLat-lat >= lat-minLat {
			latBit = 0
			maxLat = (maxLat + minLat) / 2
		} else {
			latBit = 1
			minLat = (maxLat + minLat) / 2
		}
		if maxLng-lng >= lng-minLng {
			lngBit = 0
			maxLng = (maxLng + minLng) / 2
		} else {
			lngBit = 1
			minLng = (maxLng + minLng) / 2
		}
		geo <<= 1
		geo += lngBit
		geo <<= 1
		geo += latBit
	}
	rect = &GeoRectangle{MinLat: minLat, MaxLat: maxLat, MinLng: minLng, MaxLng: maxLng}
	return
}

// 对编码成位的geohash解码成小矩形
// 源自：https://github.com/yinqiwen/geohash-int
func GeoHashBitsDecode(geohash uint64, precision uint8) *GeoRectangle {
	rect := GeoRectangle{MinLat: MIN_LATITUDE, MinLng: MIN_LONGITUDE, MaxLat: MAX_LATITUDE, MaxLng: MAX_LONGITUDE}
	var i uint8
	for i = 0; i < precision; i++ {
		var latBit, lngBit uint64
		lngBit = (geohash >> ((precision-i)*2 - 1)) & 0x01
		latBit = (geohash >> ((precision-i)*2 - 2)) & 0x01
		if latBit == 0 {
			rect.MaxLat = (rect.MaxLat + rect.MinLat) / 2
		} else {
			rect.MinLat = (rect.MaxLat + rect.MinLat) / 2
		}
		if lngBit == 0 {
			rect.MaxLng = (rect.MaxLng + rect.MinLng) / 2
		} else {
			rect.MinLng = (rect.MaxLng + rect.MinLng) / 2
		}
	}
	return &rect
}

// 获取附近的9个小格子
func GeoHashBitsNeighbors(lat, lng float64, precision uint8) (ret []uint64) {
	geohash, _ := GeoHashBitsEncode(lat, lng, precision)
	if geohash == 0 {
		return
	}
	ret = append(
		ret,
		geohash, //当前的小格子
		geohashMoveBits(geohash, precision, 0, 1),   //north
		geohashMoveBits(geohash, precision, 0, -1),  //south
		geohashMoveBits(geohash, precision, 1, 0),   //east
		geohashMoveBits(geohash, precision, -1, 0),  //west
		geohashMoveBits(geohash, precision, -1, -1), //south_west
		geohashMoveBits(geohash, precision, 1, -1),  //south_east
		geohashMoveBits(geohash, precision, -1, 1),  //north_west
		geohashMoveBits(geohash, precision, 1, 1),   //north_east
	)
	return ret
}

// 移位
func geohashMoveBits(geohash uint64, precision uint8, dx, dy int8) uint64 {
	if dx != 0 {
		var x = geohash & 0xaaaaaaaaaaaaaaaa
		var y = geohash & 0x5555555555555555
		var zz uint64 = 0x5555555555555555 >> (64 - precision*2)
		if dx > 0 {
			x = x + zz + 1
		} else {
			x = x | zz
			x = x - (zz + 1)
		}
		x &= 0xaaaaaaaaaaaaaaaa >> (64 - precision*2)
		geohash = x | y
	}
	if dy != 0 {
		var x = geohash & 0xaaaaaaaaaaaaaaaa
		var y = geohash & 0x5555555555555555
		var zz uint64 = 0xaaaaaaaaaaaaaaaa >> (64 - precision*2)
		if dy > 0 {
			y = y + zz + 1
		} else {
			y = y | zz
			y = y - (zz + 1)
		}
		y &= 0x5555555555555555 >> (64 - precision*2)
		geohash = x | y
	}
	return geohash
}
