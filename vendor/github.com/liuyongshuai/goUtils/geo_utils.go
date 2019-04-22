/**
 * 地理位置相关的小工具，主要是给ES里的地理位置查询用的
 *
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     es
 * @date        2018-04-20 14:24
 */
package goUtils

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

const (
	//地球半径
	EARTH_RADIUS = 6378137

	//浮点类型计算时候与0比较时候的容差
	FLOAT_DIFF = 2e-10

	//以下三个是用弧度计算距离用的
	RATIO0 = 0.9970707597894937
	RATIO1 = 0.0004532759255277588
	RATIO2 = -0.00017587656744607181
	RATIO3 = 0.0000005028600490023173
)

//格式化距离，最终都要输出以米为单位
// 输入"5.5km"=>"5500"
// 输入"5000m"=>"5000"
func FormatDistance(distance string) (ret float64) {
	ret = 0
	distance = strings.ToLower(distance)
	reg := regexp.MustCompile(`^[\d][\d\.]+km$`)
	if reg.MatchString(distance) { //以km为单位
		distance = strings.Replace(distance, "km", "", -1)
		dist, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			return
		}
		ret = dist * 1000
		return
	}
	reg = regexp.MustCompile(`^[\d][\d\.]+m$`)
	if reg.MatchString(distance) { //以m为单位
		distance = strings.Replace(distance, "m", "", -1)
		dist, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			return
		}
		ret = dist
		return
	}
	return
}

//计算两个经纬度间的中间位置
func MidPoint(point1, point2 GeoPoint) GeoPoint {
	if point2.IsEqual(point1) {
		return point2
	}
	lat1Arc := point1.Lat * math.Pi / 180.0
	lat2Arc := point2.Lat * math.Pi / 180.0
	lng1Arc := point1.Lng * math.Pi / 180.0
	diffLng := (point2.Lng - point1.Lng) * math.Pi / 180.0

	bx := math.Cos(lat2Arc) * math.Cos(diffLng)
	by := math.Cos(lat2Arc) * math.Sin(diffLng)

	lat3Rad := math.Atan2(math.Sin(lat1Arc)+math.Sin(lat2Arc), math.Sqrt(math.Pow(math.Cos(lat1Arc)+bx, 2)+math.Pow(by, 2)))
	lng3Rad := lng1Arc + math.Atan2(by, math.Cos(lat1Arc)+bx)

	lat3 := lat3Rad * 180.0 / math.Pi
	lng3 := lng3Rad * 180.0 / math.Pi

	//确保新点在这两点组成的直线上（新生成的点未必在这两点组成的直线上）
	line1 := MakeGeoLine(point1, point2)
	line2 := MakeGeoLine(MakeGeoPoint(90, lng3), MakeGeoPoint(-90, lng3))
	ret, isParallel, isIntersect := line1.GetIntersectPoint(line2)
	//如果平行了，即原两点的经度相同
	if isParallel && point1.Lng == point2.Lng {
		ret.Lng = point1.Lng
	} else if !isIntersect { //如果根本不相交，原样返回吧
		ret = MakeGeoPoint(lat3, lng3)
	}

	return ret
}

//在指定距离、角度上，返回另一个经纬度坐标
//跟经度线的角度，即便角度为90度，也不是一条横线，有误差
//lat、lng：源经纬度
//dist：距离，单位米
//angle：角度，如"45"
func PointAtDistAndAngle(point GeoPoint, dist, angle float64) GeoPoint {
	if dist <= 0 {
		return point
	}
	dr := dist / EARTH_RADIUS
	angle = angle * (math.Pi / 180.0)
	lat1 := point.Lat * (math.Pi / 180.0)
	lng1 := point.Lng * (math.Pi / 180.0)

	lat2 := math.Asin(math.Sin(lat1)*math.Cos(dr) + math.Cos(lat1)*math.Sin(dr)*math.Cos(angle))
	lng2 := lng1 + math.Atan2(math.Sin(angle)*math.Sin(dr)*math.Cos(lat1), math.Cos(dr)-(math.Sin(lat1)*math.Sin(lat2)))
	lng2 = math.Mod(lng2+3*math.Pi, 2*math.Pi) - math.Pi

	lat2 = lat2 * (180.0 / math.Pi)
	lng2 = lng2 * (180.0 / math.Pi)
	return GeoPoint{Lat: lat2, Lng: lng2}
}

//用三角函数来计算地球上的曲线距离，返回值为米
func EarthDistanceOld(point1, point2 GeoPoint) float64 {
	if point1.IsEqual(point2) {
		return 0
	}
	rad := math.Pi / 180.0
	lat1 := point1.Lat * rad
	lng1 := point1.Lng * rad
	lat2 := point2.Lat * rad
	lng2 := point2.Lng * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * float64(EARTH_RADIUS)
}

//用拟合出来的多项式来计算距离，返回米
func EarthDistance(point1, point2 GeoPoint) float64 {
	if point1.IsEqual(point2) {
		return 0
	}
	//经度差值
	dLng := point1.Lng - point2.Lng
	//纬度差值
	dLat := point1.Lat - point2.Lat
	//平均纬度
	avgLat := (point1.Lat + point2.Lat) / 2.0
	//计算东西方向的距离，采用三阶多项式
	ewDist := (RATIO3*avgLat*avgLat*avgLat + RATIO2*avgLat*avgLat + RATIO1*avgLat + RATIO0) * ToRadians(dLng) * EARTH_RADIUS
	//计算南北距离，采用一阶多项式
	snDist := ToRadians(dLat) * EARTH_RADIUS
	return math.Sqrt(ewDist*ewDist + snDist*snDist)
}

//将经纬度字符串转为lat、lng
func SplitGeoPoint(loc string) (float64, float64) {
	if len(loc) <= 0 || !strings.Contains(loc, ",") {
		return 0, 0
	}
	tmp := strings.Split(loc, ",")
	if len(tmp) != 2 {
		return 0, 0
	}
	lat, _ := strconv.ParseFloat(tmp[0], 64)
	lng, _ := strconv.ParseFloat(tmp[1], 64)
	return lat, lng
}

/**
随机生成指定数量的多边形，必须指定基本的矩形区域、顶点的最大个数及最小个数
所有的这些多边形的边，非相邻的不能有交点！
仅于校验之用，不保证性能！
*/
func GenPolygons(baseRect GeoRectangle, polygonNum, pointMinNum, pointMaxNum int) (ret []GeoPolygon) {
	width := baseRect.Width()
	height := baseRect.Height()
	//校验部分必要的条件是否满足
	if width <= 0 || height <= 0 {
		return
	}
	if polygonNum <= 0 {
		return
	}
	if pointMaxNum < pointMinNum {
		return
	}
	if pointMinNum < 3 {
		return
	}
	LogInfof("start GenPolygons：baseRect[%d x %d] polygonNum[%d] pointMinNum[%d] pointMaxNum[%d]",
		int(width), int(height), polygonNum, pointMinNum, pointMaxNum)
	diffNum := pointMaxNum - pointMinNum + 1
	var randFloat float64
	diffLat := baseRect.MaxLat - baseRect.MinLat
	diffLng := baseRect.MaxLng - baseRect.MinLng
	//开始随机生成多边形
	for i := 0; i < polygonNum; i++ {
		ReRandSeed()
		//顶点的数量
		vertexNum := rand.Intn(diffNum) + pointMinNum
		var points []GeoPoint
		//逐个顶点生成
		for vn := 0; vn < vertexNum; vn++ {
			//确保一个多边形的边都不相交
			for {
				ReRandSeed()
				randFloat = rand.Float64()
				//新生点的纬度
				lat := baseRect.MinLat + randFloat*diffLat
				ReRandSeed()
				randFloat = rand.Float64()
				//新生点的经度
				lng := baseRect.MinLng + randFloat*diffLng
				//将新生点放到多边形里去，检查非相邻边相交情况
				point := MakeGeoPoint(lat, lng)
				points = append(points, point)
				polygon := MakeGeoPolygon(points)
				//如果加上这个新生成的顶点后，多边形有不相邻的边相交的话，忽略掉重新生成一个新的顶点
				if polygon.IsBorderInterect() {
					points = points[0 : len(points)-1]
				} else {
					break
				}
			}
		}
		ret = append(ret, MakeGeoPolygon(points))
	}
	return
}

// 随机生成指定经纬度范围内的点
func RandomLatLng(minLat, maxLat, minLng, maxLng float64) (lat, lng float64) {
	lat = RandFloat64InRange(minLat, maxLat)
	lng = RandFloat64InRange(minLng, maxLng)
	return
}

//角度转为弧度
func ToRadians(d float64) float64 {
	return d * math.Pi / 180.0
}

//两点的向量差
func VectorDifference(p1 GeoPoint, p2 GeoPoint) GeoPoint {
	return GeoPoint{Lat: p1.Lat - p2.Lat, Lng: p1.Lng - p2.Lng}
}

//两向量叉乘
func VectorCrossProduct(p1 GeoPoint, p2 GeoPoint) float64 {
	cross := p1.Lat*p2.Lng - p1.Lng*p2.Lat
	if math.Abs(0-cross) < FLOAT_DIFF {
		cross = 0
	}
	return cross
}
