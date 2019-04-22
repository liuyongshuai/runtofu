/**
 * 所有矩形的操作的方法，逐步技术储备性质的添加中
 * 原本是多边形的一种特殊情况，但geohash小格子就是矩形的，所以此处单独拉出来了
 *
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-09-30 17:14
 */
package goUtils

import (
	"math"
)

//一个矩形
type GeoRectangle struct {
	MinLat float64 `json:"min_lat"` //最小纬度
	MinLng float64 `json:"min_lng"` //最小经度
	MaxLat float64 `json:"max_lat"` //最大纬度
	MaxLng float64 `json:"max_lng"` //最大经度
}

//校验是否为合法的矩形
func (gr *GeoRectangle) Check() bool {
	if gr.MaxLng <= gr.MinLng || gr.MaxLat <= gr.MinLat {
		return false
	}
	return true
}

//经度方向的跨度
func (gr *GeoRectangle) LngSpan() float64 {
	return gr.MaxLng - gr.MinLng
}

//纬度方向的跨度
func (gr *GeoRectangle) LatSpan() float64 {
	return gr.MaxLat - gr.MinLat
}

//判断给定的经纬度是否在小格子内，包括边界
func (gr *GeoRectangle) IsPointInRect(point GeoPoint) bool {
	if point.Lat <= gr.MaxLat &&
		point.Lat >= gr.MinLat &&
		point.Lng >= gr.MinLng &&
		point.Lng <= gr.MaxLng {
		return true
	}
	return false
}

//判断给定的经纬度是否完全在小格子内，不包括边界
func (gr *GeoRectangle) IsPointRealInRect(point GeoPoint) bool {
	if point.Lat < gr.MaxLat &&
		point.Lat > gr.MinLat &&
		point.Lng > gr.MinLng &&
		point.Lng < gr.MaxLng {
		return true
	}
	return false
}

//获取中心点坐标
func (gr *GeoRectangle) MidPoint() GeoPoint {
	point := MidPoint(
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
	)
	return point
}

//矩形X方向的边长，即纬度线方向，保持纬度相同即可，单位米
func (gr *GeoRectangle) Width() float64 {
	return EarthDistance(
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
	)
}

//矩形Y方向的边长，即经度线方向，保持经度相同即可，单位米
func (gr *GeoRectangle) Height() float64 {
	return EarthDistance(
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	)
}

//矩形的所有的点，从左下角开始，按顺时针方向来返回
func (gr *GeoRectangle) GetRectVertex() (ret []GeoPoint) {
	ret = append(ret,
		GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng}, //经纬度均最小
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng}, //纬度最大，经度最小
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng}, //经纬度均最大
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng}, //纬度最小，经度最大
	)
	return
}

//左下角的坐标：经纬度最小
func (gr *GeoRectangle) LeftBottomPoint() GeoPoint {
	return GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng}
}

//左上角的坐标：经度最小、纬度最大
func (gr *GeoRectangle) LeftUpPoint() GeoPoint {
	return GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng}
}

//右上角的坐标：经纬度最大
func (gr *GeoRectangle) RightUpPoint() GeoPoint {
	return GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng}
}

//右下角的坐标：经度最大、纬度最小
func (gr *GeoRectangle) RightBottomPoint() GeoPoint {
	return GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng}
}

//左边框线段，从上往下指
func (gr *GeoRectangle) LeftBorder() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
	}
}

//右边框线段，从上往下指
func (gr *GeoRectangle) RightBorder() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
		Point2: GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	}
}

//上边框线段，从左往右指
func (gr *GeoRectangle) TopBorder() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
	}
}

//下边框线段，从左往右指
func (gr *GeoRectangle) BottomBorder() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	}
}

//矩形的左上角、右下角的对象线
func (gr *GeoRectangle) LeftUp2RightBottomLine() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	}
}

//矩形的左下角、右上角的对象线
func (gr *GeoRectangle) LeftBottom2RightUpLine() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
	}
}

//矩形的所有的边
func (gr *GeoRectangle) GetRectBorders() (ret []GeoLine) {
	p := gr.GetRectVertex()
	ret = append(ret,
		GeoLine{Point1: p[0], Point2: p[1]},
		GeoLine{Point1: p[1], Point2: p[2]},
		GeoLine{Point1: p[2], Point2: p[3]},
		GeoLine{Point1: p[3], Point2: p[0]},
	)
	return
}

// 在矩形范围内随机生成点
func (gr *GeoRectangle) GetRandomGeoPoint() (gp GeoPoint) {
	gp.Lat = RandFloat64InRange(gr.MinLat, gr.MaxLat)
	gp.Lng = RandFloat64InRange(gr.MinLng, gr.MaxLng)
	return
}

// 合并两个矩形
func (gr *GeoRectangle) Union(rect *GeoRectangle) GeoRectangle {
	ret := GeoRectangle{
		MinLng: math.Min(gr.MinLng, rect.MinLng),
		MinLat: math.Min(gr.MinLat, rect.MinLat),
		MaxLng: math.Max(gr.MaxLng, rect.MaxLng),
		MaxLat: math.Max(gr.MaxLat, rect.MaxLat),
	}
	return ret
}

//两矩形是否相等
func (gr *GeoRectangle) IsEqual(rect GeoRectangle) bool {
	return gr.MinLng == rect.MinLng && gr.MaxLng == rect.MaxLng && gr.MinLat == rect.MinLat && gr.MaxLat == rect.MaxLat
}

//将矩形转为多边形对象
func (gr *GeoRectangle) ToPolygon() GeoPolygon {
	return MakeGeoPolygon(gr.GetRectVertex())
}

//两矩形是否有重合
func (gr *GeoRectangle) IsIntersect(rect GeoRectangle, isReal bool) bool {
	if gr.IsEqual(rect) {
		return true
	}
	v1s := gr.GetRectVertex()
	for _, v := range v1s {
		if isReal {
			if rect.IsPointRealInRect(v) {
				return true
			}
		} else {
			if rect.IsPointInRect(v) {
				return true
			}
		}

	}
	v2s := rect.GetRectVertex()
	for _, v := range v2s {
		if isReal {
			if gr.IsPointRealInRect(v) {
				return true
			}
		} else {
			if gr.IsPointInRect(v) {
				return true
			}
		}
	}
	tmp := [2]*GeoRectangle{gr, &rect}
	for i := 0; i < 2; i++ {
		a, b := tmp[0], tmp[1]
		if a.MinLng <= b.MinLng && a.MaxLng >= b.MaxLng {
			if a.MaxLat >= b.MaxLat && a.MinLat <= b.MinLat {
				return true
			}
			if a.MaxLat > b.MinLat && a.MaxLat < b.MaxLat {
				return true
			}
			if a.MinLat > b.MinLat && a.MinLat < b.MaxLat {
				return true
			}
			if a.MaxLat <= b.MaxLat && a.MinLat >= b.MinLat {
				return true
			}
		}
		tmp[0], tmp[1] = tmp[1], tmp[0]
	}
	return false
}

//拷贝【看起来没啥用】
func (gr *GeoRectangle) Clone() GeoRectangle {
	return GeoRectangle{
		MinLat: gr.MinLat,
		MinLng: gr.MinLng,
		MaxLat: gr.MaxLat,
		MaxLng: gr.MaxLng,
	}
}
