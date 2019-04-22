/**
 * 所有线段操作的方法，逐步技术储备性质的添加中
 *
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-09-30 17:12
 */
package goUtils

import (
	"fmt"
	"math"
)

//构造直线
func MakeGeoLine(p1 GeoPoint, p2 GeoPoint) GeoLine {
	return GeoLine{Point1: p1, Point2: p2}
}

//一条直接
type GeoLine struct {
	Point1 GeoPoint `json:"point1"` //起点
	Point2 GeoPoint `json:"point2"` //终点
}

//跟另一直线是否相同
func (gl *GeoLine) IsEqual(line GeoLine) bool {
	if line.Point1.IsEqual(gl.Point1) && line.Point2.IsEqual(gl.Point2) {
		return true
	}
	if line.Point1.IsEqual(gl.Point2) && line.Point2.IsEqual(gl.Point1) {
		return true
	}
	return false
}

//直线的长度
func (gl *GeoLine) Length() float64 {
	return EarthDistance(gl.Point1, gl.Point2)
}

//直线的长度
func (gl *GeoLine) FormatStr() string {
	return fmt.Sprintf("%s-%s", gl.Point1.FormatStr(), gl.Point2.FormatStr())
}

//获取直线的最小外包矩形，如果是条平行线或竖线的话，可能会有问题
func (gl *GeoLine) GetBoundsRect() GeoRectangle {
	return GeoRectangle{
		MaxLat: math.Max(gl.Point2.Lat, gl.Point1.Lat),
		MaxLng: math.Max(gl.Point2.Lng, gl.Point1.Lng),
		MinLat: math.Min(gl.Point2.Lat, gl.Point1.Lat),
		MinLng: math.Min(gl.Point2.Lng, gl.Point1.Lng),
	}
}

//是否包含某个点，基本思路：
//点为Q，线段为P1P2，判断点Q在线段上的依据是：(Q-P1)×(P2-P1)=0
//且Q在以P1P2为对角定点的矩形内
func (gl *GeoLine) IsContainPoint(p GeoPoint) bool {
	rect := gl.GetBoundsRect()
	if !rect.IsPointInRect(p) {
		return false
	}
	if p.IsEqual(gl.Point2) || p.IsEqual(gl.Point1) {
		return true
	}
	p1 := VectorDifference(gl.Point1, gl.Point2)
	p2 := VectorDifference(p, gl.Point1)
	cross := VectorCrossProduct(p1, p2)
	return cross == 0
}

//两直线间的夹角，用向量来解决，角度有可能为负
func (gl *GeoLine) AngleWithLine(l GeoLine) float64 {
	r1 := VectorDifference(l.Point1, l.Point2)
	r2 := VectorDifference(gl.Point1, gl.Point2)
	//向量值
	cross := VectorCrossProduct(r1, r2)
	//两向量模的积
	m := math.Sqrt(r1.Lat*r1.Lat+r1.Lng*r1.Lng) * math.Sqrt(r2.Lat*r2.Lat+r2.Lng*r2.Lng)
	angle := math.Asin(cross / m)
	//弧度转为角度
	angle = angle * 180 / math.Pi
	return angle
}

//与另一条直线是否相交、平行
func (gl *GeoLine) IsIntersectWithLine(line GeoLine) (isIntersect bool, isParallel bool) {
	_, isParallel, isIntersect = gl.GetIntersectPoint(line)
	return
}

//两直线相交的判断，考虑了部分共线的情况
//参考：https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
//	interPoints：是所有的交点，如果两线段有部分重合，此时会有两个交点
//	isParallel：两直线是否平行
//若len(interPoints)==1且isParallel=true表示平行且端点重合
//若len(interPoints)==2且isParallel=true表示两线段部分重合
func (gl *GeoLine) GetIntersectPoints(line GeoLine) (interPoints []GeoPoint, isParallel bool) {
	//如果其中一线是点
	if gl.Point1.IsEqual(gl.Point2) {
		if line.IsContainPoint(gl.Point1) {
			interPoints = append(interPoints, gl.Point1)
			isParallel = true
		}
		return
	}
	if line.Point1.IsEqual(line.Point2) {
		if gl.IsContainPoint(line.Point1) {
			interPoints = append(interPoints, line.Point1)
			isParallel = true
		}
		return
	}
	//如果两线段完全重合
	if line.Point1.IsEqual(gl.Point1) && line.Point2.IsEqual(gl.Point2) ||
		line.Point2.IsEqual(gl.Point1) && line.Point1.IsEqual(gl.Point2) {
		interPoints = append(interPoints, line.Point1, line.Point2)
		isParallel = true
		return
	}

	p := gl.Point1
	//一线段的向量差
	r := VectorDifference(gl.Point2, gl.Point1)
	q := line.Point1
	//另一线段的向量差
	s := VectorDifference(line.Point2, line.Point1)
	//两线段向量差的叉乘
	rCrossS := VectorCrossProduct(r, s)
	qMinusP := VectorDifference(q, p)
	//两线段平行的情况：可能共线、也可能完全无关（完全可以区分出来共线和平行不相交的情况，这里偷懒了）
	if rCrossS == 0 {
		isParallel = true
		//VectorCrossProduct(qMinusP, r)==0表示共线，否则就是平行但不相交
		if gl.IsContainPoint(line.Point1) {
			interPoints = append(interPoints, line.Point1)
		}
		if gl.IsContainPoint(line.Point2) {
			interPoints = append(interPoints, line.Point2)
		}
		if line.IsContainPoint(gl.Point1) {
			interPoints = append(interPoints, gl.Point1)
		}
		if line.IsContainPoint(gl.Point2) {
			interPoints = append(interPoints, gl.Point2)
		}
		//稍微去一下重，考虑到只有一个交点且平行的线段，该交点会被添加两次。
		//还有可能一线段完全包含另一线段
		if len(interPoints) > 0 {
			tmp := map[GeoPoint]struct{}{}
			for _, p := range interPoints {
				tmp[p] = struct{}{}
			}
			interPoints = []GeoPoint{}
			for p := range tmp {
				interPoints = append(interPoints, p)
			}
		}
		return
	}
	//不共线的话就剩下正常的相交的情况了
	t := VectorCrossProduct(qMinusP, s) / rCrossS
	u := VectorCrossProduct(qMinusP, r) / rCrossS
	//正常的相交的情况，若r × s ≠ 0 and 0 ≤ t ≤ 1 and 0 ≤ u ≤ 1表示两线段在点相交： p + t r = q + u s
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		//两线段上的交点
		p1 := GeoPoint{Lat: gl.Point1.Lat + t*r.Lat, Lng: gl.Point1.Lng + t*r.Lng}
		p2 := GeoPoint{Lat: line.Point1.Lat + u*s.Lat, Lng: line.Point1.Lng + u*s.Lng}
		p := p1
		//如果在计算的时候有点小小的误差，这里直接取中间得了，理论上这两个点应该相等
		if !p1.IsEqual(p2) {
			p = MidPoint(p1, p2)
		}
		isParallel = false
		interPoints = append(interPoints, p)
	}
	return
}

//【在处理线段共线上有点问题】求两直线交点的坐标
//参考：https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
//返回值：交点、是否平行、是否相交
func (gl *GeoLine) GetIntersectPoint(line GeoLine) (interPoint GeoPoint, isParallel bool, isIntersect bool) {
	ps, para := gl.GetIntersectPoints(line)
	isParallel = para
	if len(ps) > 0 {
		isIntersect = true
		interPoint = ps[0]
	}
	return
}

//拷贝【看起来没啥用】
func (gl *GeoLine) Clone() GeoLine {
	return GeoLine{
		Point1: gl.Point1.Clone(),
		Point2: gl.Point2.Clone(),
	}
}
