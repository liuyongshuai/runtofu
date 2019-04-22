// 多边形的布尔操作一般应包括：UNION、DIFFERENCE、INTERSECTION、XOR。
// 但本处的多边形一般为商家配送多边形，不支持非相邻的边有交点、中间有空洞、自相交的情况出现，所以只保留了部分布尔操作。
//
// 本工具依据论文《A new algorithm for computing Boolean operations on polygons》所描述的算法来完成。
// 论文见：http://www.cs.ucr.edu/~vbz/cs230papers/martinez_boolean.pdf
// 参考实现：
// 		在平面坐标系上的实现：
// 			https://github.com/akavel/polyclip-go
//		C++实现：
//			http://www4.ujaen.es/~fmartin/bool_op.html
//		相关说明：
// 			http://www.angusj.com/delphi/clipper.php
//		2D多边形布尔操作说明：
// 			http://thebuildingcoder.typepad.com/blog/2013/09/boolean-operations-for-2d-polygons.html
//			http://thebuildingcoder.typepad.com/blog/2009/02/boolean-operations-for-2d-polygons.html
//		线段相交平面扫描法：
//			https://blog.csdn.net/avan_lau/article/details/10967597
// 			http://geomalgorithms.com/a09-_intersect-3.html
//		参考博客：
//			https://www.cnblogs.com/wuhanhoutao/archive/2008/03/08/1096224.html
//		计算几何的总结：
//			http://www.twinklingstar.cn/category/computational-geometry/
// 核心思想如下：
//		所有的边的交点作为endpoint放到一个队列里，此队列初始值为多边形的顶点（也是边与边的交点），以后碰到不同多边形边与边的交点再放进去，对于处理完的边则要剔除掉。
//		这些顶点或交点也叫一个事件，分为边的相邻事件、边的相交事件。
//		将这些边按一定的规则排序，从左到右用sweepline逐个扫描，有关边的排序规则见代码说明。
//		sweepline是一个垂线，它扫过endpoint，也会有边按顺序有交点。将扫描过的endpoing放互eventqueue里面，里面的event（endpoint）是动态变化且有序的。
//
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-10-08 11:51

package goUtils

import (
	"math"
	"sort"
)

//多边形的类型
type polygonType int

const (
	POLYGON_TYPE_SUBJECT  polygonType = iota //目标多边形
	POLYGON_TYPE_CLIPPING                    //被裁剪的多边形
)

//边的类型，主要用于重叠的边
type edgeType int

const (
	EDGE_TYPE_NORMAL           edgeType = iota
	EDGE_TYPE_NON_CONTRIBUTING  //没用的边，最终的结果里不会包含的边
	EDGE_TYPE_SAME_TRANSITION   //针对inOut值的
	EDGE_TYPE_DIFFERENT_TRANSITION
)

//布尔操作类型
type boolOpType int

const (
	OP_TYPE_UNION boolOpType = iota
	OP_TYPE_INTERSECTION
	OP_TYPE_DIFFERENCE
)

//多边形裁剪结构体
type polyBoolOperation struct {
	subject  GeoPolygon
	clipping GeoPolygon
	eventQueue
}

//构造虚拟的扫描线，从左到右依次扫描所有的线段。
//它是基于事件的，这里的事件分为两类：线段的端点及线段与线段的交点，分别称为端点事件、相交事件。
//初始时，事件队列是所有线段的端点。
//如果两个线段相交，那么它们在扫描线列表中是相邻的！
func (pbo *polyBoolOperation) compute(operation boolOpType) (ret []GeoPolygon) {
	subjectPointNum := len(pbo.subject.Points)
	clippingPointNum := len(pbo.clipping.Points)
	if subjectPointNum == 0 && clippingPointNum == 0 {
		return
	}
	//有多边形为空
	if clippingPointNum*subjectPointNum == 0 {
		switch operation {
		case OP_TYPE_UNION:
			if subjectPointNum == 0 {
				ret = append(ret, pbo.clipping.Clone())
			} else {
				ret = append(ret, pbo.subject.Clone())
			}
		case OP_TYPE_DIFFERENCE:
			ret = append(ret, pbo.subject.Clone())
		}
		return
	}

	//两个多边形的外包矩形没有重叠部分，交集必然为空，并集的话直接合并即可
	subjectBoxRect := pbo.subject.GetBoundsRect()
	clippingBoxRect := pbo.clipping.GetBoundsRect()
	if !subjectBoxRect.IsIntersect(clippingBoxRect, false) {
		switch operation {
		case OP_TYPE_UNION:
			ret = append(ret, pbo.subject.Clone())
			ret = append(ret, pbo.clipping.Clone())
		case OP_TYPE_DIFFERENCE:
			ret = append(ret, pbo.subject.Clone())
		}
		return
	}

	//将所有的边添加到事件队列中去，按从左到右的顺序
	subjectLines := pbo.subject.GetPolygonBorders()
	for _, l := range subjectLines {
		addProcessedSegment(&pbo.eventQueue, l, POLYGON_TYPE_SUBJECT)
	}
	clippingLines := pbo.clipping.GetPolygonBorders()
	for _, l := range clippingLines {
		addProcessedSegment(&pbo.eventQueue, l, POLYGON_TYPE_CLIPPING)
	}

	connector := connector{}

	//从左到右扫描所有多边形的所有的边
	sweepLines := sweepline{}

	//两外包矩形里的最大经度的最小值
	minMaxLng := math.Min(subjectBoxRect.MaxLng, clippingBoxRect.MaxLng)

	for !pbo.eventQueue.IsEmpty() {
		//上一个点、下一个点
		var prev, next *endpoint

		//处理一个端点
		e := pbo.eventQueue.popQueue()

		//此点越界，对交集没有用了就
		switch {
		case operation == OP_TYPE_INTERSECTION && e.p.Lng > minMaxLng:
			fallthrough
		case operation == OP_TYPE_INTERSECTION && e.p.Lng > subjectBoxRect.MaxLng:
			return connector.toPolygon()
		}

		//如果此点为左点，先插入扫描线列表里
		if e.left {
			//前一个点（线段）、后一个点（线段）都有可能跟当前线段相交
			//它们也决定了当前点（线段）的inside/inOut标志值
			pos := sweepLines.insert(e)
			prev = nil
			if pos > 0 {
				prev = sweepLines[pos-1]
			}
			next = nil
			if pos < len(sweepLines)-1 {
				next = sweepLines[pos+1]
			}

			//计算inside、inOut标志信息
			//这两个标志位只有left endpoint才有用
			//强调：inside表示当前线段是否在另一个多边形内部（边上也算）
			switch {
			//如果前面没有线段，e是最上面的边
			case prev == nil:
				e.inside, e.inout = false, false
				//prev!=normal的话，有可能是共线的边，也有可能前面有两条共线的边
			case prev.edgeType != EDGE_TYPE_NORMAL:
				if pos < 2 { //当前线段跟prev，部分重合，在求交点时被push进去了
					e.inside, e.inout = false, false
					if prev.polygonType != e.polygonType {
						e.inside = true
					} else {
						e.inout = true
					}
				} else { //前两个线段都部分重合
					prevTwo := sweepLines[pos-2]
					if prev.polygonType == e.polygonType {
						e.inout = !prev.inout
						e.inside = !prevTwo.inout
					} else {
						e.inout = !prevTwo.inout
						e.inside = !prev.inout
					}
				}
				//prev跟当前点属于同一个多边形
			case e.polygonType == prev.polygonType:
				//在eventQueue相邻的两边，如果有一个在另一个多边形里面，那另一条边也在另一个多边形里面
				e.inside = prev.inside
				e.inout = !prev.inout
				//prev 跟当前线段，不属于同一个多边形
			default:
				e.inside = !prev.inout
				e.inout = prev.inside
			}

			//当前线段有可能跟下一个线段相交
			if next != nil {
				pbo.possibleIntersection(e, next)
			}
			//当前线段有可能跟前一个线段相交
			if prev != nil {
				pbo.possibleIntersection(prev, e)
			}
		} else { //如果是右点，则当前点所属的线段要从S中移除，表示其处理完了已经
			otherPos := -1
			for i := range sweepLines {
				if sweepLines[i].equals(e.other) {
					otherPos = i
					break
				}
			}

			if otherPos != -1 {
				prev = nil
				if otherPos > 0 {
					prev = sweepLines[otherPos-1]
				}
				next = nil
				if otherPos < len(sweepLines)-1 {
					next = sweepLines[otherPos+1]
				}
			}

			//检查当前线段是否为布尔操作的一部分
			switch e.edgeType {
			case EDGE_TYPE_NORMAL: //所有的操作都用得着
				switch operation {
				case OP_TYPE_INTERSECTION: //如果是两多边形相交的话，则此线段要另一多边形内部的要保留
					if e.other.inside {
						connector.add(e.segment())
					}
				case OP_TYPE_UNION: //如果是两多边形相并的话，则此线段要另一多边形外部的要保留
					if !e.other.inside {
						connector.add(e.segment())
					}
				case OP_TYPE_DIFFERENCE:
					if (e.polygonType == POLYGON_TYPE_SUBJECT && !e.other.inside) ||
						(e.polygonType == POLYGON_TYPE_CLIPPING && e.other.inside) {
						connector.add(e.segment())
					}
				}
			case EDGE_TYPE_SAME_TRANSITION: //只有求交、求并用得着
				if operation == OP_TYPE_INTERSECTION || operation == OP_TYPE_UNION {
					connector.add(e.segment())
				}
			case EDGE_TYPE_DIFFERENT_TRANSITION:
				if operation == OP_TYPE_DIFFERENCE {
					connector.add(e.segment())
				}
			}

			//删除当前点所属的线段，它已经处理完了
			if otherPos != -1 {
				sweepLines.remove(sweepLines[otherPos])
			}
			//将当前点所属的线段移除了，则它的前一个、后一个线段就变成相邻的了，也要计算它们是否可能相交
			if next != nil && prev != nil {
				pbo.possibleIntersection(next, prev)
			}
		}
	}

	return connector.toPolygon()
}

//寻找可能的交点
func (pbo *polyBoolOperation) possibleIntersection(e1, e2 *endpoint) {
	l1 := MakeGeoLine(e1.p, e1.other.p)
	l2 := MakeGeoLine(e2.p, e2.other.p)
	ips, _ := l1.GetIntersectPoints(l2)
	intersectNum := len(ips)
	if intersectNum <= 0 {
		return
	}
	interPoint := ips[0]
	//只有一个交点，但交点在端点，即是相邻的两条边
	if intersectNum == 1 && (e1.p.IsEqual(e2.p) || e1.other.p.IsEqual(e2.other.p)) {
		return
	}

	//同一个多边形的两条边相交，且有两个交点，目前在商家的配送范围里是不允许的
	if intersectNum == 2 && e1.polygonType == e2.polygonType {
		return
	}

	//只有一个交点，这是正常情况
	if intersectNum == 1 {
		//交点不是线段e1的端点
		if !e1.p.IsEqual(interPoint) && !e1.other.p.IsEqual(interPoint) {
			//将e1分成两个线段，
			pbo.divideSegment(e1, interPoint)
		}
		//交点不是线段e2的端点
		if !e2.p.IsEqual(interPoint) && !e2.other.p.IsEqual(interPoint) {
			//将e2分成两个线段
			pbo.divideSegment(e2, interPoint)
		}
		return
	}

	/**
	  以下处理全是交点个数>=2的情况，即线段部分重合，就要用得上edgeType来区分了
	*/

	sortedEvents := make([]*endpoint, 0)
	switch {
	case e1.p.IsEqual(e2.p): //起点重合
		sortedEvents = append(sortedEvents, nil)
	case endpointCompare(e1, e2):
		sortedEvents = append(sortedEvents, e2, e1)
	default:
		sortedEvents = append(sortedEvents, e1, e2)
	}

	switch {
	case e1.other.p.IsEqual(e2.other.p): //端点重合
		sortedEvents = append(sortedEvents, nil)
	case endpointCompare(e1.other, e2.other):
		sortedEvents = append(sortedEvents, e2.other, e1.other)
	default:
		sortedEvents = append(sortedEvents, e1.other, e2.other)
	}

	//两个线段根本就是相等的，两个端点完全相同，保留一个就可以了，此处将e1干掉了
	if len(sortedEvents) == 2 {
		e1.edgeType, e1.other.edgeType = EDGE_TYPE_NON_CONTRIBUTING, EDGE_TYPE_NON_CONTRIBUTING
		if e1.inout == e2.inout {
			e2.edgeType, e2.other.edgeType = EDGE_TYPE_SAME_TRANSITION, EDGE_TYPE_SAME_TRANSITION
		} else {
			e2.edgeType, e2.other.edgeType = EDGE_TYPE_DIFFERENT_TRANSITION, EDGE_TYPE_DIFFERENT_TRANSITION
		}
		return
	}

	//两线段共享一个端点，即起点或终点相同的部分重合情况，此时，一线段完全包含另一线段
	//重合的那部分线段保留待用，其中一个线段多出来那部分就可以丢弃了
	//如（nil,e1.other,e2.other）、(e1,e2,nil)
	if len(sortedEvents) == 3 {
		sortedEvents[1].edgeType, sortedEvents[1].other.edgeType = EDGE_TYPE_NON_CONTRIBUTING, EDGE_TYPE_NON_CONTRIBUTING
		var idx int
		if sortedEvents[0] != nil { //共享右端点
			idx = 0
		} else { //共享左端点
			idx = 2
		}
		if e1.inout == e2.inout {
			sortedEvents[idx].other.edgeType = EDGE_TYPE_SAME_TRANSITION
		} else {
			sortedEvents[idx].other.edgeType = EDGE_TYPE_DIFFERENT_TRANSITION
		}
		if sortedEvents[0] != nil {
			pbo.divideSegment(sortedEvents[0], sortedEvents[1].p)
		} else {
			pbo.divideSegment(sortedEvents[2].other, sortedEvents[1].p)
		}
		return
	}

	//此时len(sortedEvents) == 4，线段只有部分区域是重合的，但端点不同，一线段不完全包含另一线段
	//如(e1,e2,e1.other,e2.other)，总共分为三段
	if sortedEvents[0] != sortedEvents[3].other {
		//第二个线段的类型就不能是normal了
		sortedEvents[1].edgeType = EDGE_TYPE_NON_CONTRIBUTING
		if e1.inout == e2.inout {
			sortedEvents[2].edgeType = EDGE_TYPE_SAME_TRANSITION
		} else {
			sortedEvents[2].edgeType = EDGE_TYPE_DIFFERENT_TRANSITION
		}
		pbo.divideSegment(sortedEvents[0], sortedEvents[1].p)
		//这是重合部分的线段
		pbo.divideSegment(sortedEvents[1], sortedEvents[2].p)
		return
	}

	//剩下的情况，端点不相同，是一个线段完全包含另一个
	//如(e1,e2,e2.other,e1.other)
	sortedEvents[1].edgeType, sortedEvents[1].other.edgeType = EDGE_TYPE_NON_CONTRIBUTING, EDGE_TYPE_NON_CONTRIBUTING
	pbo.divideSegment(sortedEvents[0], sortedEvents[1].p)
	if e1.inout == e2.inout {
		sortedEvents[3].other.edgeType = EDGE_TYPE_SAME_TRANSITION
	} else {
		sortedEvents[3].other.edgeType = EDGE_TYPE_DIFFERENT_TRANSITION
	}
	pbo.divideSegment(sortedEvents[3].other, sortedEvents[2].p)
}

//将一条边根据其上的一个点分成两部分
func (pbo *polyBoolOperation) divideSegment(e *endpoint, p GeoPoint) {
	//(e,p)这里e是左端点，p其实是右端点
	r := &endpoint{
		p:           p,
		left:        false,
		polygonType: e.polygonType,
		other:       e,
		edgeType:    e.edgeType,
	}
	//(p,e.other)这里p是左端点，e.other其实是右端点
	l := &endpoint{
		p:           p,
		left:        true,
		polygonType: e.polygonType,
		other:       e.other,
		edgeType:    e.other.edgeType,
	}

	//还是要比较一下
	if endpointCompare(l, e.other) {
		e.other.left = true
		e.left = false
	}

	e.other.other = l
	e.other = r

	pbo.eventQueue.pushQueue(l)
	pbo.eventQueue.pushQueue(r)
}

//当扫描线穿过多边形的时候的两多边形的边交点的位置
type endpoint struct {
	p      GeoPoint  //当前的点
	left   bool      //当前点是否为线段的左点，线段表示为 (p, other->p)
	polygonType      //当前点所属的多边形
	other  *endpoint //当前点所属线段的另一个点的结构体
	inout  bool      //针对该点所属的边（p, other->p）的所属的多边形而言，当从点（p.Lng, -无穷）发出的射线穿过此边时，是否为从里到外过渡
	edgeType         //边的类型，见上注释
	inside bool      //只有当left=true时才用得着，表示当前线段是否在另一个多边形的内部
}

//比较两个端点，必须所有的参数都相同才算是相同
func (ep *endpoint) equals(e *endpoint) bool {
	return ep.p.IsEqual(e.p) &&
		ep.left == e.left &&
		ep.polygonType == e.polygonType &&
		ep.other == e.other &&
		ep.inout == e.inout &&
		ep.edgeType == e.edgeType &&
		ep.inside == e.inside
}

//返回端点所属的线段
func (ep *endpoint) segment() GeoLine {
	return MakeGeoLine(ep.p, ep.other.p)
}

//判断两个点的上下关系
func (ep *endpoint) below(x GeoPoint) bool {
	var l GeoLine
	if ep.left {
		l = MakeGeoLine(ep.p, ep.other.p)
	} else {
		l = MakeGeoLine(ep.other.p, ep.p)
	}
	return x.IsBelow(l)
}

//判断两个点的上下关系
func (ep *endpoint) above(x GeoPoint) bool {
	return !ep.below(x)
}

//事件队列，每个边的端点为一个事件，两条边的交点也为一个事件
type eventQueue struct {
	elements []*endpoint
	sorted   bool
}

//将一个端点事件入队列
func (q *eventQueue) pushQueue(e *endpoint) {
	//如果没有排序的话，直接塞进去即可
	if !q.sorted {
		q.elements = append(q.elements, e)
		return
	}

	//如果还没有排序，用插入排序算法来排序
	length := len(q.elements)
	if length == 0 {
		q.elements = append(q.elements, e)
		return
	}

	//插入并排序
	q.elements = append(q.elements, nil)
	i := length - 1
	for i >= 0 && endpointCompare(e, q.elements[i]) {
		q.elements[i+1] = q.elements[i]
		i--
	}
	q.elements[i+1] = e
}

//弹出一个
func (q *eventQueue) popQueue() *endpoint {
	if !q.sorted {
		sort.Sort(eventQueueComparer(q.elements))
		q.sorted = true
	}
	x := q.elements[len(q.elements)-1]
	q.elements = q.elements[:len(q.elements)-1]
	return x
}

//端点事件队列是否为空
func (q *eventQueue) IsEmpty() bool {
	return len(q.elements) == 0
}

//比较队列中点的比较器，实现了sort.Interface接口
type eventQueueComparer []*endpoint

func (q eventQueueComparer) Len() int           { return len(q) }
func (q eventQueueComparer) Less(i, j int) bool { return endpointCompare(q[i], q[j]) }
func (q eventQueueComparer) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

//点链，一堆点的连接序列
type geoPointChain struct {
	closed bool
	points []GeoPoint
}

//实例化一个点链
func newChain(s GeoLine) *geoPointChain {
	return &geoPointChain{
		closed: false,
		points: []GeoPoint{s.Point1, s.Point2}}
}

//将指定的点放到链的前端
func (c *geoPointChain) pushFront(p GeoPoint) {
	c.points = append([]GeoPoint{p}, c.points...)
}

//将指定的点放到链的末端
func (c *geoPointChain) pushBack(p GeoPoint) {
	c.points = append(c.points, p)
}

//判断可否将一个边连接成链
func (c *geoPointChain) linkSegment(s GeoLine) bool {
	//第一个点
	firstPoint := c.points[0]
	//最后一个点
	lastPoint := c.points[len(c.points)-1]

	switch true {
	//起点跟第一个点相同
	case s.Point1.IsEqual(firstPoint):
		//如果终点再跟最后一个点相同则为闭合的多边形
		if s.Point2.IsEqual(lastPoint) {
			c.closed = true
		} else {
			c.pushFront(s.Point2)
		}
		return true
	//如果终点跟最后一个点相同
	case s.Point2.IsEqual(lastPoint):
		//起点再跟第一个相同，也是闭合的多边形
		if s.Point1.IsEqual(firstPoint) {
			c.closed = true
		} else {
			c.pushBack(s.Point1)
		}
		return true
	//终点跟第一个点相同
	case s.Point2.IsEqual(firstPoint):
		//如果起点跟最后一个点相同，形成了闭合的多边形
		if s.Point1.IsEqual(lastPoint) {
			c.closed = true
		} else {
			c.pushFront(s.Point1)
		}
		return true
	//起点跟最后一个点相同
	case s.Point1.IsEqual(lastPoint):
		//终点跟第一个点相同，形成了闭合的多边形
		if s.Point2.IsEqual(firstPoint) {
			c.closed = true
		} else {
			c.pushBack(s.Point2)
		}
		return true
	}
	return false
}

//将另一个点链连接到当前链上
func (c *geoPointChain) linkChain(other *geoPointChain) bool {
	//第一个点
	firstPoint := c.points[0]
	//最后一个点
	lastPoint := c.points[len(c.points)-1]
	//另一个点链的第一个点
	otherFirstPoint := other.points[0]
	//另一个点链的最后一个点
	otherLastPoint := other.points[len(other.points)-1]

	//其他点链第一个点跟当前点链的最后一个相同，直接拼接起来即可
	if otherFirstPoint.IsEqual(lastPoint) {
		c.points = append(c.points, other.points[1:]...)
		goto success
	}

	//其他点链的最后一个点跟当前点链的第一个点相同，也是直接拼接起来即可
	if otherLastPoint.IsEqual(firstPoint) {
		c.points = append(other.points, c.points[1:]...)
		goto success
	}

	//其他点链的第一个点跟当前的第一个点相同，将其他点链反转后再拼接起来，但要干掉当前点链的第一个点，因为重合了嘛
	if otherFirstPoint.IsEqual(firstPoint) {
		c.points = append(reverseGeoPoints(other.points), c.points[1:]...)
		goto success
	}

	//其他点链的最后一个点跟当前点链的最后一个点链相同，反转其他点链并拼接到当前的后面即可
	if otherLastPoint.IsEqual(lastPoint) {
		c.points = append(c.points[:len(c.points)-1], reverseGeoPoints(other.points)...)
		goto success
	}

	return false

success:
	other.points = []GeoPoint{}
	return true
}

//模拟扫描线，将扫到的点暂存起来，这些点（包括顶点、交点）也叫一个事件，即顶点事件、交点事件
//这些点所属的线段是有顺序的，从左到右排序，其比较规则见相关方法
type sweepline []*endpoint

//移除一个点事件
func (s *sweepline) remove(key *endpoint) {
	for i, el := range *s {
		if el.equals(key) {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}

//扫描线列表里添加一个点，返回其下标
func (s *sweepline) insert(item *endpoint) int {
	length := len(*s)
	if length == 0 {
		*s = append(*s, item)
		return 0
	}

	*s = append(*s, &endpoint{})
	i := length - 1
	for i >= 0 && segmentCompare(item, (*s)[i]) {
		(*s)[i+1] = (*s)[i]
		i--
	}
	(*s)[i+1] = item
	return i + 1
}

//暂存中间结果，并最终拼装成结果多边形
type connector struct {
	openPolys   []geoPointChain
	closedPolys []geoPointChain
}

//添加可能的一条边
func (c *connector) add(s GeoLine) {
	//遍历所有开放的多边形的线段链，判断可否组成多边形
	for j := range c.openPolys {
		chain := &c.openPolys[j]
		if !chain.linkSegment(s) {
			continue
		}

		//如果链路已经闭合了
		if chain.closed {
			//只有两个点，不能组成一个多边形，非法情况，强制给打开
			if len(chain.points) == 2 {
				chain.closed = false
				return
			}
			//将打开的多边形迁移到闭合的多边形里去
			c.closedPolys = append(c.closedPolys, c.openPolys[j])
			c.openPolys = append(c.openPolys[:j], c.openPolys[j+1:]...)
			return
		}

		//如果多边形没有闭合
		k := len(c.openPolys)
		for i := j + 1; i < k; i++ {
			//尝试将打开的链路拼接到闭合的链上去，并删除掉拼接后的打开的链路
			if chain.linkChain(&c.openPolys[i]) {
				c.openPolys = append(c.openPolys[:i], c.openPolys[i+1:]...)
				return
			}
		}
		return
	}

	//此边不能跟任何开放的多边形连接在一起
	c.openPolys = append(c.openPolys, *newChain(s))
}

//连接器最终转为多边形
func (c *connector) toPolygon() []GeoPolygon {
	var poly []GeoPolygon
	for _, chain := range c.closedPolys {
		var ps []GeoPoint
		for _, p := range chain.points {
			ps = append(ps, p)
		}
		poly = append(poly, MakeGeoPolygon(ps))
	}
	return poly
}

//比较两个endpoint的大小，返回true时e1在e2后面处理，先处理e2
func endpointCompare(e1, e2 *endpoint) bool {
	//如果经度不一样，先比较经度
	if e1.p.Lng != e2.p.Lng {
		return e1.p.Lng > e2.p.Lng
	}

	//如果经度一样，纬度小的先处理
	if e1.p.Lat != e2.p.Lat {
		return e1.p.Lat > e2.p.Lat
	}

	//同一个点，一个是左点，一个是右点，先处理右点
	if e1.left != e2.left {
		return e1.left
	}

	//同一个点，都是左点，或都是右点，先处理底部的（有关below的判断见其他方法）
	return e1.above(e2.other.p)
}

//将一堆点列表反转
func reverseGeoPoints(list []GeoPoint) []GeoPoint {
	length := len(list)
	other := make([]GeoPoint, length)
	for i := range list {
		other[length-i-1] = list[i]
	}
	return other
}

//有关线段的比较
func segmentCompare(e1, e2 *endpoint) bool {
	e1Line := MakeGeoLine(e1.p, e1.other.p)
	switch {
	case e1 == e2: //重合了
		return false
	case !e1Line.IsContainPoint(e2.p): //点e2不在e1所属线段上
		fallthrough
	case !e1Line.IsContainPoint(e2.other.p): //点e2的另一个点也不在e1所属线段上，说明两上线段不共线，
		//如果他们的左点相同，要按右点排序
		if e1.p.IsEqual(e2.p) {
			return e1.below(e2.other.p)
		}
		//左点不相同，e1所属线段在e2所属线段之后插入扫描列表上去
		if endpointCompare(e1, e2) {
			return e2.above(e1.p)
		}
		//e2所属线段在e1所属线段之后插入扫描列表上去
		return e1.below(e2.p)
	case e1.p.IsEqual(e2.p): //线段共线
		return false
	}
	return endpointCompare(e1, e2)
}

//将两个多边形所有的边放到一个队列里去，并区分来自哪个多边形
func addProcessedSegment(q *eventQueue, segment GeoLine, polyType polygonType) {
	//线段的两个端点相等，异常数据
	if segment.Point1.IsEqual(segment.Point2) {
		return
	}

	//拼装端点结构体数据
	e1 := &endpoint{
		p:           segment.Point1,
		left:        true,
		polygonType: polyType,
	}
	e2 := &endpoint{
		p:           segment.Point2,
		left:        true,
		polygonType: polyType,
		other:       e1,
	}
	e1.other = e2

	//比较两个点大小，来决定哪个是左点
	switch {
	case e1.p.Lng < e2.p.Lng:
		e2.left = false
	case e1.p.Lng > e2.p.Lng:
		e1.left = false
	case e1.p.Lat < e2.p.Lat: //竖线，下面的端点判为左
		e2.left = false
	default:
		e1.left = false
	}

	//将两个点添加到事件队列中去，并按从左到右的顺序添加
	q.pushQueue(e1)
	q.pushQueue(e2)
}
