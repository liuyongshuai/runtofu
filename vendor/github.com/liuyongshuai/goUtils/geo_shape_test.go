/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-09-30 17:21
 */
package goUtils

import (
	"fmt"
	"os"
	"testing"
)

//定义一个矩形
var (
	//大概在北京四环左上角附近（西北角）
	testWestNorthPoint = MakeGeoPoint(39.989877, 116.28576)
	//大概在北京四环右下角附近（东南角）
	testEastSouthPoint = MakeGeoPoint(39.842258, 116.487986)
	//东北角
	testEastNorthPoint = MakeGeoPoint(39.989877, 116.487986)
	//西南角
	testWestSouthPoint = MakeGeoPoint(39.842258, 116.28576)

	//矩形的上边
	testTopLine = MakeGeoLine(testWestNorthPoint, testEastNorthPoint)
	//矩形的下边
	testBottomLine = MakeGeoLine(testWestSouthPoint, testEastSouthPoint)
	//矩形的左边
	testLeftLine = MakeGeoLine(testWestSouthPoint, testWestNorthPoint)
	//矩形的右边
	testRightLine = MakeGeoLine(testEastSouthPoint, testEastNorthPoint)
	//矩形的左上-右下的对角线
	testLeftTopRightBottom = MakeGeoLine(testWestNorthPoint, testEastSouthPoint)
	//矩形的左下-右上的对角线
	testLeftBottomRightTop = MakeGeoLine(testWestSouthPoint, testEastNorthPoint)
)

func TestGeoLine(t *testing.T) {
	testStart()

	//各边的夹角
	fmt.Printf("rightLine_bottomLine angle: %v\n", testRightLine.AngleWithLine(testBottomLine))
	fmt.Printf("rightLine_leftLine angle: %v\n", testRightLine.AngleWithLine(testLeftLine))
	fmt.Printf("leftBottom_rightTop_leftTop_rightBottom angle: %v\n", testLeftBottomRightTop.AngleWithLine(testLeftTopRightBottom))

	//topLineDistance:17247.51175508048
	//bottomLineDistance:17284.7283002726
	//leftLineDistance:16432.87191160098
	//rightLineDistance:16432.87191160098
	fmt.Printf("topLineDistance:%v\nbottomLineDistance:%v\nleftLineDistance:%v\nrightLineDistance:%v\n",
		testTopLine.Length(), testBottomLine.Length(), testLeftLine.Length(), testRightLine.Length(),
	)

	//上下两边，只平行，不相交
	isIntersect, isParallel := testTopLine.IsIntersectWithLine(testBottomLine)
	fmt.Printf("topLine bottomLine\tisParallel=%v\tisIntersect=%v\n", isParallel, isIntersect)
	//左右两边，，只平行，不相交
	isIntersect, isParallel = testLeftLine.IsIntersectWithLine(testRightLine)
	fmt.Printf("leftLine rightLine\tisParallel=%v\tisIntersect=%v\n", isParallel, isIntersect)
	//左边/下边，不平行，只相交,interPoint={39.842258 116.28576}
	interPoint, isParallel, isIntersect := testLeftLine.GetIntersectPoint(testBottomLine)
	fmt.Printf("leftLine bottomLine\tisParallel=%v\tisIntersect=%v\tinterPoint=%v\n", isParallel, isIntersect, interPoint)

	//两个对角线，只相交不平行，interPoint={39.9160675 116.38687300000001}
	interPoint, isParallel, isIntersect = testLeftTopRightBottom.GetIntersectPoint(testLeftBottomRightTop)
	fmt.Printf("leftTop_rightBottom leftBottom_rightTop\tisParallel=%v\tisIntersect=%v\tinterPoint=%v\n", isParallel, isIntersect, interPoint)
	fmt.Printf("interPoint inLeftTop_rightBottom=%v\tinLeftBottom_rightTop=%v\n",
		testLeftTopRightBottom.IsContainPoint(interPoint), testLeftBottomRightTop.IsContainPoint(interPoint),
	)

	//对角线交点
	fmt.Println(interPoint.IsBelow(testLeftTopRightBottom))
	fmt.Println(testWestSouthPoint.IsBelow(testLeftTopRightBottom))
	testEnd()
}

var (
	/**
	line1 line2平行但不相交
	line2 line3平行、共线、两个交点
	line4 line1完全相等
	line5 line1平行，共享一端点，只有一个交点应该
	line6 line1一个完全包含另一个，不共端点
	line7 line8不平行，有一个交点
	line7 line1不平行，有一个交点，且交点在端点
	*/
	testIntersectPointsLine1 = MakeGeoLine(MakeGeoPoint(39.989877, 116.28576), MakeGeoPoint(39.989877, 116.487986))
	testIntersectPointsLine2 = MakeGeoLine(MakeGeoPoint(39.842258, 116.28576), MakeGeoPoint(39.842258, 116.487986))
	testIntersectPointsLine3 = MakeGeoLine(MakeGeoPoint(39.842258, 116.08576), MakeGeoPoint(39.842258, 116.387986))
	testIntersectPointsLine4 = MakeGeoLine(MakeGeoPoint(39.989877, 116.28576), MakeGeoPoint(39.989877, 116.487986))
	testIntersectPointsLine5 = MakeGeoLine(MakeGeoPoint(39.989877, 115.28576), MakeGeoPoint(39.989877, 116.28576))
	testIntersectPointsLine6 = MakeGeoLine(MakeGeoPoint(39.989877, 116.08576), MakeGeoPoint(39.989877, 116.987986))
	testIntersectPointsLine7 = MakeGeoLine(testWestNorthPoint, testEastSouthPoint)
	testIntersectPointsLine8 = MakeGeoLine(testWestSouthPoint, testEastNorthPoint)
)

//两线段共线的测试
func TestGeoLine_GetIntersectPoints(t *testing.T) {
	testStart()
	interP12, isParallel := testIntersectPointsLine1.GetIntersectPoints(testIntersectPointsLine2)
	fmt.Fprintf(os.Stdout, "line1 line2（平行但不相交） \t：=>\t%v %v\n", interP12, isParallel)
	interP23, isParallel := testIntersectPointsLine2.GetIntersectPoints(testIntersectPointsLine3)
	fmt.Fprintf(os.Stdout, "line2 line3（平行、共线、两个交点） \t：=>\t%v %v\n", interP23, isParallel)
	interP41, isParallel := testIntersectPointsLine4.GetIntersectPoints(testIntersectPointsLine1)
	fmt.Fprintf(os.Stdout, "line4 line1（完全相等） \t：=>\t%v %v\n", interP41, isParallel)
	interP51, isParallel := testIntersectPointsLine5.GetIntersectPoints(testIntersectPointsLine1)
	fmt.Fprintf(os.Stdout, "line5 line1（平行，共享一端点，只有一个交点应该） \t：=>\t%v %v\n", interP51, isParallel)
	interP61, isParallel := testIntersectPointsLine6.GetIntersectPoints(testIntersectPointsLine1)
	fmt.Fprintf(os.Stdout, "line6 line1（一个完全包含另一个，不共端点）\t：=>\t %v %v\n", interP61, isParallel)
	interP78, isParallel := testIntersectPointsLine7.GetIntersectPoints(testIntersectPointsLine8)
	fmt.Fprintf(os.Stdout, "line7 line8（不平行，有一个交点） \t：=>\t%v %v\n", interP78, isParallel)
	interP71, isParallel := testIntersectPointsLine7.GetIntersectPoints(testIntersectPointsLine1)
	fmt.Fprintf(os.Stdout, "line7 line1（不平行，有一个交点，且交点在端点） \t：=>\t%v %v\n", interP71, isParallel)
	testEnd()
}
