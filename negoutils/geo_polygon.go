/**
 * 所有（配送）多边形的操作的方法，逐步技术储备性质的添加中
 *
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-05-29 18:12
 */
package negoutils

import (
	"math"
)

// 根据经纬度字符串构造多边形，字符串格式"lat,lng",根据ES对经纬度串的要求来的
func MakeGeoPolygonByStr(locs []string) GeoPolygon {
	var points []GeoPoint
	for _, loc := range locs {
		lat, lng := SplitGeoPoint(loc)
		points = append(points, GeoPoint{Lat: lat, Lng: lng})
	}
	return MakeGeoPolygon(points)
}

// 构造多边形
func MakeGeoPolygon(points []GeoPoint) GeoPolygon {
	pl := len(points)
	if pl >= 3 {
		if !points[0].IsEqual(points[pl-1]) {
			points = append(points, points[0])
		}
	}
	ret := GeoPolygon{Points: points, GeoHashType: GEOHASH_TYPE_NORMAL}
	return ret
}

// 一个多边形
type GeoPolygon struct {
	Points      []GeoPoint   `json:"points"`       //一堆顶点，必须是首尾相连有序的
	Borders     []GeoLine    `json:"borders"`      //所有的边
	Rect        GeoRectangle `json:"rect"`         //最小外包矩形
	GeoHashType int          `json:"geohash_type"` //切格子的类型
}

// 获取所有的顶点
func (gp *GeoPolygon) GetPoints() []GeoPoint {
	return gp.Points
}

// 设置切格子的方式
func (gp *GeoPolygon) SetGeoHashType(gtype int) {
	if gtype != GEOHASH_TYPE_BITS && gtype != GEOHASH_TYPE_NORMAL {
		return
	}
	gp.GeoHashType = gtype
}

// 是否有边相交（相怜的不算）
func (gp *GeoPolygon) IsBorderInterect() bool {
	if !gp.Check() {
		return false
	}
	borders := gp.GetPolygonBorders()
	for _, line1 := range borders {
		for _, line2 := range borders {
			if line1.Point1.IsEqual(line2.Point1) || line1.Point1.IsEqual(line2.Point2) {
				continue
			}
			if line1.Point2.IsEqual(line2.Point1) || line1.Point1.IsEqual(line2.Point2) {
				continue
			}
			isInter, _ := line1.IsIntersectWithLine(line2)
			if isInter {
				return true
			}
		}
	}
	return false
}

// 多边形的所有的边
func (gp *GeoPolygon) GetPolygonBorders() (ret []GeoLine) {
	if !gp.Check() {
		return
	}
	if len(gp.Borders) > 0 {
		return gp.Borders
	}
	points := gp.GetPoints()
	l := len(points)
	p0 := points[0]
	for i := 1; i < l; i++ {
		p := points[i]
		ret = append(ret, GeoLine{Point1: p0, Point2: p})
		p0 = p
	}
	ret = append(ret, GeoLine{Point1: points[l-1], Point2: points[0]})
	gp.Borders = ret
	return
}

// 添加点
func (gp *GeoPolygon) AddPoint(p GeoPoint) {
	gp.Points = append(gp.Points, p)
}

// 将多边形处理成字符串的切片格式
func (gp *GeoPolygon) FormatStringArray() (ret []string) {
	for _, p := range gp.Points {
		ret = append(ret, p.FormatStr())
	}
	return
}

// 将多边形处理成字符串的切片格式
func (gp *GeoPolygon) FormatStringMultiArray() (ret [][]string) {
	var tmp []string
	for _, p := range gp.Points {
		tmp = append(tmp, p.FormatStr())
	}
	ret = append(ret, tmp)
	return
}

// 判断是否是合法的多边形
func (gp *GeoPolygon) Check() bool {
	if len(gp.Points) < 3 {
		return false
	}
	//如果多边形边长超过100km，玩去，没法处理了！
	rect := gp.GetBoundsRect()
	width := rect.Width()
	height := rect.Height()
	var maxDist float64 = 100000
	if width >= maxDist || height >= maxDist {
		return false
	}
	if gp.GeoHashType != GEOHASH_TYPE_NORMAL && gp.GeoHashType != GEOHASH_TYPE_BITS {
		gp.GeoHashType = GEOHASH_TYPE_NORMAL
	}
	return true
}

// 获取最小外包矩形
func (gp *GeoPolygon) GetBoundsRect() GeoRectangle {
	if gp.Rect.Width() > 0 || gp.Rect.Height() > 0 {
		return gp.Rect
	}
	var maxLat = MIN_LATITUDE
	var maxLng = MIN_LONGITUDE
	var minLat = MAX_LATITUDE
	var minLng = MAX_LONGITUDE
	for _, p := range gp.Points {
		maxLat = math.Max(maxLat, p.Lat)
		minLat = math.Min(minLat, p.Lat)
		maxLng = math.Max(maxLng, p.Lng)
		minLng = math.Min(minLng, p.Lng)
	}
	rect := GeoRectangle{MaxLat: maxLat, MaxLng: maxLng, MinLat: minLat, MinLng: minLng}
	gp.Rect = rect
	return rect
}

// 判断点是否在多边形内部，此处使用最简单的射线法判断
// 边数较多时性能不高，只适合在写入时小批量判断
// 计算射线与多边形各边的交点，如果是偶数，则点在多边形外，否则在多边形内。
// 还会考虑一些特殊情况，如点在多边形顶点上，点在多边形边上等特殊情况。
// 参考【有BUG】：http://api.map.baidu.com/library/GeoUtils/1.2/src/GeoUtils.js
func (gp *GeoPolygon) IsPointInPolygon(p GeoPoint) bool {
	if !p.Check() || !gp.Check() {
		return false
	}
	//判断最小外包矩形
	rect := gp.GetBoundsRect()
	if !rect.IsPointInRect(p) {
		return false
	}

	//交点总数
	var interCount = 0
	//相邻的两个顶点
	var p1, p2 GeoPoint
	//顶点个数
	PNum := len(gp.Points)

	//逐个顶点的判断
	p1 = gp.Points[0]
	points := gp.Points

	//遍历所有的边，寻找跟此点向右发出的射线的相交情况
	//所有跟此射线相交的边，必须有一部分在射线的下方
	for i := 1; i < PNum; i++ {
		//其他顶点
		p2 = points[i%PNum]
		//正好落在了顶点上
		if p1.IsEqual(p) || p2.IsEqual(p) {
			return true
		}
		maxLat := math.Max(p1.Lat, p2.Lat)
		minLat := math.Min(p1.Lat, p2.Lat)
		minLng := math.Min(p1.Lng, p2.Lng)
		maxLng := math.Max(p1.Lng, p2.Lng)
		//射线没有交点
		if p.Lat < minLat || p.Lat > maxLat {
			p1 = p2
			continue
		}
		//射线有可能有交点
		if p.Lat > minLat && p.Lat < maxLat {
			//点在此边的左边
			if p.Lng <= math.Max(p1.Lng, p2.Lng) {
				//此边为一条横线
				if p1.Lat == p2.Lat && p.Lng >= minLng {
					return true
				}
				//一条竖线
				if p1.Lng == p2.Lng {
					if p1.Lng == p.Lng {
						return true
					} else {
						interCount++
					}
				} else {
					//判断是否相交
					xInters := (p.Lat-p1.Lat)*(p2.Lng-p1.Lng)/(p2.Lat-p1.Lat) + p1.Lng
					if math.Abs(p.Lng-xInters) < FLOAT_DIFF {
						return true
					}
					if p.Lng < xInters {
						interCount++
					}
				}
			}
		} else {
			//正好在一条横线上
			if p.Lat == p1.Lat && p1.Lat == p2.Lat && p.Lng >= minLng && p.Lng <= maxLng {
				return true
			}
			//如果交点在边的顶点上
			if p.Lat == p2.Lat && p.Lng <= p2.Lng {
				//本处的多边形是始于第一个点，终于第一个点，此点在points出现两次，因为要做这样的判断
				p3 := points[(i+1)%PNum]
				if p3.IsEqual(p2) {
					p3 = points[(i+2)%PNum]
				}
				//如果另一个边的另一个顶点在射线的上方，忽略此边，只算一次。否则都算
				if p.Lat >= math.Min(p1.Lat, p3.Lat) && p.Lat <= math.Max(p1.Lat, p3.Lat) {
					interCount++
				} else {
					interCount += 2
				}
			}
		}
		p1 = p2
	}

	if interCount%2 == 0 {
		return false
	} else {
		return true
	}
}

// 对多边形暴力切格子，笨而慢，但比较完备，只做double check用！
func (gp *GeoPolygon) ViolentSplitGeoHashRect(precision int) (inRect, interRect []string) {
	if !gp.Check() {
		return
	}

	//切格子用的geohash精度
	stp := precision

	splitInfo := gp.getSplitGeoHashRect(stp)
	basePoint := splitInfo.basePoint //向下、向右遍历小格子的基准点
	verNum := splitInfo.verNum       //竖直方向上的小格子数量
	horiNum := splitInfo.horiNum     //水平方向上的小格子数量
	diffLat := splitInfo.diffLat     //小格子的纬度跨度值
	diffLng := splitInfo.diffLng     //小格子的经度跨度值
	interRect = splitInfo.interRect
	if len(interRect) > 0 {
		return
	}

	//多边形的各边
	polygonBorders := gp.GetPolygonBorders()
	polygonPoints := gp.GetPoints()

	//从左到右、从上到下遍历所有的小格子
	for vi := 0; vi < verNum; vi++ {
		//最左边的小格子的经度都一样，纬度逐次减小
		baseLat := basePoint.Lat - float64(vi)*diffLat
		baseLng := basePoint.Lng
		for hi := 0; hi < horiNum; hi++ {
			//当前小格子的geo值及小格子矩形
			geo, tRect := GeneralGeoHashEncode(baseLat, baseLng+float64(hi)*diffLng, stp, gp.GeoHashType)
			rectMidPoint := tRect.MidPoint()
			rectPoints := tRect.GetRectVertex()
			rectBorders := tRect.GetRectBorders()
			//完全在多边形内部的情况：四个顶点在内部、与多边形边的交点只能是此边顶点
			inPointNum := 0
			for _, p := range rectPoints {
				if gp.IsPointInPolygon(p) {
					inPointNum++
				}
			}
			if gp.IsPointInPolygon(rectMidPoint) {
				inPointNum++
			}
			isContinue := false
			//小格子的边和多边形各边的交点情况
			var vePoints []GeoPoint
			var bPoints []GeoPoint
			for _, rectB := range rectBorders {
				for _, polygonB := range polygonBorders {
					interP, isParallel, isIntersect := rectB.GetIntersectPoint(polygonB)
					if isIntersect && !isParallel {
						//如果交点位于小格子
						if interP.IsEqual(tRect.RightBottomPoint()) ||
							interP.IsEqual(tRect.RightUpPoint()) ||
							interP.IsEqual(tRect.LeftBottomPoint()) ||
							interP.IsEqual(tRect.LeftUpPoint()) {
							vePoints = append(vePoints, interP)
						} else {
							bPoints = append(bPoints, interP)
						}
					}
				}
			}
			//4个顶点全在多边形内部，且所有的交点都是多边形的顶点
			if inPointNum == 5 {
				if len(bPoints) <= 0 {
					inRect = append(inRect, geo)
					continue
				}
				interRect = append(interRect, geo)
				continue
			}
			//既没有交点在边上，也没有在顶点上的，没有相关性
			if len(vePoints) <= 0 && len(bPoints) <= 0 {
				continue
			}
			//交点全在小格子的顶点上
			if len(vePoints) > 0 && len(bPoints) <= 0 {
				for _, p := range polygonPoints {
					if tRect.IsPointRealInRect(p) {
						interRect = append(interRect, geo)
						isContinue = true
						break
					}
				}
				if isContinue {
					continue
				}
			}
			//相交在小格子的某个边上，任一个多边形顶点在小格子内即可
			if len(bPoints) > 0 {
				for _, p := range polygonPoints {
					if tRect.IsPointRealInRect(p) {
						interRect = append(interRect, geo)
						isContinue = true
						break
					}
				}
			}
			if isContinue {
				continue
			}
			//斜对角线
			for _, polygonB := range polygonBorders {
				interP, isParallel, isIntersect := polygonB.GetIntersectPoint(tRect.LeftUp2RightBottomLine())
				if isIntersect &&
					!isParallel &&
					!interP.IsEqual(tRect.LeftUpPoint()) &&
					!interP.IsEqual(tRect.RightBottomPoint()) {
					interRect = append(interRect, geo)
					break
				}
				interP, isParallel, isIntersect = polygonB.GetIntersectPoint(tRect.LeftBottom2RightUpLine())
				if isIntersect &&
					!isParallel &&
					!interP.IsEqual(tRect.LeftBottomPoint()) &&
					!interP.IsEqual(tRect.RightUpPoint()) {
					interRect = append(interRect, geo)
					break
				}
			}
		}
	}

	return
}

/*
*
用类似射线法的思想去将多边形切成多个小格子
*/
func (gp *GeoPolygon) RaySplitGeoHashRect(stp int) (inRect, interRect []string) {
	if !gp.Check() {
		return
	}
	splitInfo := gp.getSplitGeoHashRect(stp)
	basePoint := splitInfo.basePoint //向下、向右遍历小格子的基准点
	verNum := splitInfo.verNum       //竖直方向上的小格子数量
	horiNum := splitInfo.horiNum     //水平方向上的小格子数量
	diffLat := splitInfo.diffLat     //小格子的纬度跨度值
	diffLng := splitInfo.diffLng     //小格子的经度跨度值
	geoRect := splitInfo.geoRect     //最小的geohash的外包矩形
	interRect = splitInfo.interRect
	if len(interRect) > 0 {
		return
	}

	//线段与多边形的交点情况，缓存交点用的，保证每个切线只计算一次
	lineInterCache := map[string]map[GeoLine]GeoPoint{}

	//从左到右、从上到下遍历所有的小格子
	for verInterator := 0; verInterator < verNum; verInterator++ {
		//最左边的小格子的经度都一样，纬度逐次减小
		baseLat := basePoint.Lat - float64(verInterator)*diffLat
		baseLng := basePoint.Lng
		for horiInterator := 0; horiInterator < horiNum; horiInterator++ {
			//当前小格子的geo值及小格子矩形
			tLng := baseLng + float64(horiInterator)*diffLng
			geo, tRect := GeneralGeoHashEncode(baseLat, tLng, stp, gp.GeoHashType)

			//小格子上边框的延长线、及与多边形每个边的交点情况
			topLine := GeoLine{
				Point1: GeoPoint{Lat: tRect.MaxLat, Lng: geoRect.MinLng - 1},
				Point2: GeoPoint{Lat: tRect.MaxLat, Lng: geoRect.MaxLng + 1},
			}
			topStr := topLine.FormatStr()
			topInters, ok := lineInterCache[topStr]
			if !ok {
				topInters = gp.interPointsWithHorizontalLine(topLine)
				lineInterCache[topStr] = topInters
			}

			//小格子下边框延长线、及与多边形各边交点情况
			bottomLine := GeoLine{
				Point1: GeoPoint{Lat: tRect.MinLat, Lng: geoRect.MinLng - 1},
				Point2: GeoPoint{Lat: tRect.MinLat, Lng: geoRect.MaxLng + 1},
			}
			bottomStr := bottomLine.FormatStr()
			bottomInters, ok := lineInterCache[bottomStr]
			if !ok {
				bottomInters = gp.interPointsWithHorizontalLine(bottomLine)
				lineInterCache[bottomStr] = bottomInters
			}

			//小格子左边框延长线、及与多边形各边交点情况
			leftLine := GeoLine{
				Point1: GeoPoint{Lat: geoRect.MaxLat + 1, Lng: tRect.MinLng},
				Point2: GeoPoint{Lat: geoRect.MinLat - 1, Lng: tRect.MinLng},
			}
			leftStr := leftLine.FormatStr()
			leftInters, ok := lineInterCache[leftStr]
			if !ok {
				leftInters = gp.interPointsWithVertialLine(leftLine)
				lineInterCache[leftStr] = leftInters
			}

			//小格子右边框延长线、及与多边形各边交点情况
			rightLine := GeoLine{
				Point1: GeoPoint{Lat: geoRect.MaxLat + 1, Lng: tRect.MaxLng},
				Point2: GeoPoint{Lat: geoRect.MinLat - 1, Lng: tRect.MaxLng},
			}
			rightStr := rightLine.FormatStr()
			rightInters, ok := lineInterCache[rightStr]
			if !ok {
				rightInters = gp.interPointsWithVertialLine(rightLine)
				lineInterCache[rightStr] = rightInters
			}

			//TODO 既然得到了上下左右四线的交点情况，可以将此行或此列的小格子一起判断，节省判断次数！
			//TODO 此处存在同一行或同一列的小格子对交点重复判断的情况，待优化！

			isContinue := false

			//底边框跟多边形的交点不符合在多边形内部的情况
			for border, interPoint := range bottomInters {
				//如果交点位于小格子底边框的非顶点上，即只位于底边框线段中间某个位置上
				//且多边形的相应边的另一顶点在下边框的上方，此时必定是半包围的小格子
				if interPoint.Lng > tRect.MinLng && interPoint.Lng < tRect.MaxLng {
					if border.Point1.Lat > interPoint.Lat || border.Point2.Lat > interPoint.Lat {
						interRect = append(interRect, geo)
						isContinue = true
						break
					}
				}
				//考虑到特殊情况，即此边正好和小格子对角线部分重合
				//如果交点在小格子的某个角上，同时又跟对角线上对应的点相交
				topInterPoint, ok := topInters[border]
				if !ok {
					continue
				}
				if interPoint.Lng == tRect.MinLng && topInterPoint.Lng == tRect.MaxLng ||
					interPoint.Lng == tRect.MaxLng && topInterPoint.Lng == tRect.MinLng {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}

			//左边垂线的交点在边框上
			for border, interPoint := range leftInters {
				//如果交点位于小格子左边框的非顶点上，即只位于左边框线段中间某个位置上
				//且多边形的相应边的另一顶点在左边框的右方，此时必定是半包围的小格子
				if interPoint.Lat < tRect.MaxLat && interPoint.Lat > tRect.MinLat {
					if border.Point1.Lng > interPoint.Lng || border.Point2.Lng > interPoint.Lng {
						interRect = append(interRect, geo)
						isContinue = true
						break
					}
				}
				//考虑到特殊情况，即此边正好和小格子对角线部分重合
				//如果交点在小格子的某个角上，同时又跟对角线上对应的点相交
				rightInterPoint, ok := rightInters[border]
				if !ok {
					continue
				}
				if interPoint.Lat == tRect.MaxLat && rightInterPoint.Lat == tRect.MinLat ||
					interPoint.Lat == tRect.MinLat && rightInterPoint.Lat == tRect.MaxLat {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}

			//右边框的交点在边框上
			for border, interPoint := range rightInters {
				//如果交点位于小格子右边框的非顶点上，即只位于右边框线段中间某个位置上
				//且多边形的相应边的另一顶点在右边框的左方，此时必定是半包围的小格子
				if interPoint.Lat < tRect.MaxLat && interPoint.Lat > tRect.MinLat {
					if border.Point1.Lng < interPoint.Lng || border.Point2.Lng < interPoint.Lng {
						interRect = append(interRect, geo)
						isContinue = true
						break
					}
				}
			}
			if isContinue {
				continue
			}

			//对于上下边框，判断小格子左右两边的交点情况
			leftNum := 0
			rightNum := 0
			for border, interPoint := range topInters {
				//上边框向左的射线跟多边形的交点情况
				if interPoint.Lng <= tRect.MinLng {
					leftNum++
					continue
				}
				//上边框向右的射线跟多边形的交点情况
				if interPoint.Lng >= tRect.MaxLng {
					rightNum++
					continue
				}
				//如果交点位于小格子上边框的非顶点上，即只位于上边框线段中间某个位置上
				//且多边形的相应边的另一顶点在上边框的下方，此时必定是半包围的小格子
				if border.Point1.Lat < interPoint.Lat || border.Point2.Lat < interPoint.Lat {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}
			if leftNum%2 == 1 && rightNum%2 == 1 {
				inRect = append(inRect, geo)
				continue
			}
		}
	}

	return
}

// 切格子开始遍历时需要的信息
type splitRectBaseInfo struct {
	basePoint        GeoPoint
	verNum, horiNum  int
	diffLat, diffLng float64
	geoRect, minRect GeoRectangle
	interRect        []string
}

// 获取切格子的基准点及水平方向、垂直方向上的格子个数
func (gp *GeoPolygon) getSplitGeoHashRect(precision int) (ret splitRectBaseInfo) {
	//先提取最小外包矩形，并取取四个点的geohash格子
	minRect := gp.GetBoundsRect()

	//一个临时小格子，计算经纬度差用的
	_, tmpRect := GeneralGeoHashEncode(minRect.MaxLat, minRect.MinLng, precision, gp.GeoHashType)
	tmpMidPoint := tmpRect.MidPoint()
	diffLat := tmpRect.MaxLat - tmpRect.MinLat
	diffLng := tmpRect.MaxLng - tmpRect.MinLng

	//如果外包矩形太小，取其中心点的
	var interRect []string
	if minRect.Width() <= tmpRect.Width() && minRect.Height() <= tmpRect.Height() {
		mid := minRect.MidPoint()
		geo, _ := GeneralGeoHashEncode(mid.Lat, mid.Lng, precision, gp.GeoHashType)
		interRect = append(interRect, geo)
		return
	}

	//左上角的格子，基准点
	_, leftUpRect := GeneralGeoHashEncode(tmpMidPoint.Lat, tmpMidPoint.Lng, precision, gp.GeoHashType)
	//左下角的小格子，经纬、纬度全最小
	_, leftBottom := GeneralGeoHashEncode(minRect.MinLat, minRect.MinLng, precision, gp.GeoHashType)
	//右上角的小格子，经纬、纬度全最大
	_, rightUp := GeneralGeoHashEncode(minRect.MaxLat, minRect.MaxLng, precision, gp.GeoHashType)

	//最小的geohash整数倍的外包矩形，可能要比最小外包矩形大一些
	geoRect := GeoRectangle{
		MaxLat: rightUp.MaxLat,
		MaxLng: rightUp.MaxLng,
		MinLat: leftBottom.MinLat,
		MinLng: leftBottom.MinLng,
	}

	//垂直、水平方向的格子数，可能出现xx.9999或者xxx.0001这样的情况
	tmpVNum := geoRect.Height() / rightUp.Height()
	tmpHNum := geoRect.Width() / rightUp.Width()
	verNum := int(tmpVNum)
	horiNum := int(tmpHNum)
	if math.Abs(tmpVNum-float64(verNum)) > 0.1 || verNum <= 0 {
		verNum++
	}
	if math.Abs(tmpHNum-float64(horiNum)) > 0.1 || horiNum <= 0 {
		horiNum++
	}

	//基准点，从这个点开始往右、往下推进
	basePoint := leftUpRect.MidPoint()
	ret = splitRectBaseInfo{
		basePoint: basePoint,
		verNum:    verNum,
		horiNum:   horiNum,
		diffLat:   diffLat,
		diffLng:   diffLng,
		geoRect:   geoRect,
		minRect:   minRect,
		interRect: interRect,
	}
	return
}

// 一条横线和多边形的交点，与横线部分重合的不算、交点在多边形顶点的时位于直接上方的不算
func (gp *GeoPolygon) interPointsWithHorizontalLine(line GeoLine) (ret map[GeoLine]GeoPoint) {
	ret = map[GeoLine]GeoPoint{}
	maxLng := math.Max(line.Point1.Lng, line.Point2.Lng)
	minLng := math.Min(line.Point1.Lng, line.Point2.Lng)
	lineLat := line.Point2.Lat
	borders := gp.GetPolygonBorders()
	for _, border := range borders {
		//由于line是一条直线，纬度相等，只不过经度有变化
		if border.Point2.Lat > lineLat && border.Point1.Lat > lineLat {
			continue
		}
		if border.Point2.Lat < lineLat && border.Point1.Lat < lineLat {
			continue
		}
		//平行线不要
		if border.Point1.Lat == border.Point2.Lat {
			continue
		}
		//如果交点在其顶点上，并且另一点的纬度大于横线的不要，否则就算有交点
		if border.Point1.Lat == lineLat && border.Point1.Lng >= minLng && border.Point1.Lng <= maxLng {
			if border.Point2.Lat <= lineLat {
				border.Point1.Lat = lineLat
				ret[border] = border.Point1
			}
			continue
		}
		if border.Point2.Lat == lineLat && border.Point2.Lng >= minLng && border.Point2.Lng <= maxLng {
			if border.Point1.Lat <= lineLat {
				border.Point1.Lat = lineLat
				ret[border] = border.Point2
			}
			continue
		}
		//普通的相交
		p, isParallel, isInter := border.GetIntersectPoint(line)
		if isInter && !isParallel {
			p.Lat = lineLat
			ret[border] = p
		}
	}
	return
}

// 一条垂线和多边形的交点
func (gp *GeoPolygon) interPointsWithVertialLine(line GeoLine) (ret map[GeoLine]GeoPoint) {
	ret = map[GeoLine]GeoPoint{}
	lineLng := line.Point2.Lng
	borders := gp.GetPolygonBorders()
	for _, border := range borders {
		if border.Point2.Lng > lineLng && border.Point1.Lng > lineLng {
			continue
		}
		if border.Point2.Lng < lineLng && border.Point1.Lng < lineLng {
			continue
		}
		//垂线不要
		if border.Point1.Lng == border.Point2.Lng {
			continue
		}
		//普通的相交
		p, isParallel, isInter := border.GetIntersectPoint(line)
		if isInter && !isParallel {
			p.Lng = lineLng
			ret[border] = p
		}
	}
	return
}

// 求两个多边形的并集，返回多个多边形
func (gp *GeoPolygon) UnionWithPoly(poly GeoPolygon) []GeoPolygon {
	pbo := polyBoolOperation{
		subject:  *gp,
		clipping: poly,
	}
	ret := pbo.compute(OP_TYPE_UNION)
	return ret
}

// 求两个多边形的交集，返回多个多边形
func (gp *GeoPolygon) IntersectionWithPoly(poly GeoPolygon) []GeoPolygon {
	pbo := polyBoolOperation{
		subject:  *gp,
		clipping: poly,
	}
	ret := pbo.compute(OP_TYPE_INTERSECTION)
	return ret
}

// 求两个多边形的差集，在当前多边形且不在参数多边形里面的
func (gp *GeoPolygon) DifferenceWithPoly(poly GeoPolygon) []GeoPolygon {
	pbo := polyBoolOperation{
		subject:  *gp,
		clipping: poly,
	}
	ret := pbo.compute(OP_TYPE_DIFFERENCE)
	return ret
}

// 拷贝【看起来没啥用】
func (gp *GeoPolygon) Clone() GeoPolygon {
	var points []GeoPoint
	var borders []GeoLine
	for _, p := range gp.GetPoints() {
		points = append(points, p.Clone())
	}
	for _, b := range gp.GetPolygonBorders() {
		borders = append(borders, b.Clone())
	}
	rect := gp.GetBoundsRect()
	ret := GeoPolygon{
		Points:      points,
		Borders:     borders,
		Rect:        rect.Clone(),
		GeoHashType: gp.GeoHashType,
	}
	return ret
}
