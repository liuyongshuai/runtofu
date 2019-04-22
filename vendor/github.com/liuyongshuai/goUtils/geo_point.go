/**
 * 所有点（经纬度坐标）操作的方法
 *
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-04-23 17:00
 */
package goUtils

import (
	"fmt"
	"math"
)

//字符串格式的经纬度构造点："lat,lng"
func MakeGeoPointFromStr(loc string) GeoPoint {
	lat, lng := SplitGeoPoint(loc)
	return GeoPoint{Lat: lat, Lng: lng}
}

//构造点
func MakeGeoPoint(lat, lng float64) GeoPoint {
	return GeoPoint{Lat: lat, Lng: lng}
}

//一个点
type GeoPoint struct {
	Lat float64 `json:"lat"` //纬度
	Lng float64 `json:"lng"` //经度
}

//返回字符串表示的形式
func (gp *GeoPoint) FormatStr() string {
	return fmt.Sprintf("%v,%v", gp.Lat, gp.Lng)
}

//返回数组表示的形式
func (gp *GeoPoint) FormatArray() [2]float64 {
	return [2]float64{gp.Lat, gp.Lng}
}

//根据指定的距离、角度构造另一个点
func (gp *GeoPoint) PointAtDistAndAngle(distance, angle float64) GeoPoint {
	return PointAtDistAndAngle(*gp, distance, angle)
}

//以指定点为中心，构造一个正多这形，通常用于格子查询，方便对此多边形切格子
//需要指定边数、各边的顶点到中心点的距离
func (gp *GeoPoint) BuildPolygon(borderNum, distance int) GeoPolygon {
	angle := 360 / float64(borderNum)
	var points []GeoPoint
	for i := 0; i < borderNum; i++ {
		p := gp.PointAtDistAndAngle(float64(distance), angle*float64(i)+angle/2)
		points = append(points, p)
	}
	return MakeGeoPolygon(points)
}

//跟另一个点是否相等
func (gp *GeoPoint) IsEqual(p GeoPoint) bool {
	return gp.Lat == p.Lat && gp.Lng == p.Lng
}

//判断一个点的经纬度是否合法
func (gp *GeoPoint) Check() bool {
	if gp.Lng > MAX_LONGITUDE ||
		gp.Lng < MIN_LONGITUDE ||
		gp.Lat > MAX_LATITUDE ||
		gp.Lat < MIN_LATITUDE {
		return false
	}
	return true
}

//是否在指定线段的下方（南边），采用向量的办法
//1：点在线段上方（包括线段的上延长线上）
//-1：点在线段下方(包括线段的下延长线上)
//0：点在线段上
func (gp *GeoPoint) IsBelow(l GeoLine) bool {
	if gp.Lat > math.Max(l.Point1.Lat, l.Point2.Lat) {
		return false
	}
	if gp.Lat < math.Min(l.Point1.Lat, l.Point2.Lat) {
		return true
	}
	r := VectorDifference(l.Point2, *gp)
	s := VectorDifference(l.Point1, *gp)
	cross := VectorCrossProduct(r, s)
	return cross < 0
}

//拷贝【看起来没啥用】
func (gp *GeoPoint) Clone() GeoPoint {
	return GeoPoint{Lat: gp.Lat, Lng: gp.Lng}
}

//获取纬度，看起来没啥用的方法
func (gp *GeoPoint) GetLat() float64 {
	return gp.Lat
}

//获取经度，看起来没啥用的方法
func (gp *GeoPoint) GetLng() float64 {
	return gp.Lng
}
