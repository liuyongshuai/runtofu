/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-04-17 13:54
 */
package negoutils

import (
	"fmt"
	"os"
	"testing"
)

var Testpoints = []GeoPoint{
	{Lng: 116.350514, Lat: 39.912832},
	{Lng: 116.300209, Lat: 39.939281},
	{Lng: 116.342753, Lat: 39.955654},
}

func TestGeoHashEncode(t *testing.T) {
	testStart()

	for _, pInfo := range Testpoints {
		fmt.Fprintf(os.Stdout, "Lat:%f\tLng:%f\n", pInfo.Lat, pInfo.Lng)
		for i := 4; i <= 10; i++ {
			_, square := GeoHashEncode(pInfo.Lat, pInfo.Lng, i)
			xdist := int(square.Width())
			ydist := int(square.Height())
			dt := fmt.Sprintf("%vm x %vm", xdist, ydist)
			fmt.Fprintf(os.Stdout, "\tprecision:%d\tdist:%-14s\n", i, dt)
		}
	}
	testEnd()
}

func TestGeoHashDecode(t *testing.T) {
	testStart()

	//一坨的地理位置信息
	points := []GeoPoint{
		{Lat: 31.675095455, Lng: 120.309268963407},
		{Lat: 31.54224843, Lng: 120.378810983558},
	}
	for _, pInfo := range points {
		fmt.Fprintf(os.Stdout, "Lat:%f\tLng:%f\n", pInfo.Lat, pInfo.Lng)
		for i := 4; i <= 10; i++ {
			geo, deGeo := GeoHashEncode(pInfo.Lat, pInfo.Lng, i)
			p := deGeo.MidPoint()
			midLat, midLng := p.Lat, p.Lng
			fmt.Fprintf(os.Stdout, "\tprecision:%d\tgeoHash:%-10s\t\n", i, geo)
			fmt.Fprintf(os.Stdout, "\t\tmaxLat:%v\tminLat:%v\n", deGeo.MaxLat, deGeo.MinLat)
			fmt.Fprintf(os.Stdout, "\t\tminLng:%v\tmaxLng:%v\n", deGeo.MinLng, deGeo.MaxLng)
			fmt.Fprintf(os.Stdout, "\t\tmidLat:%v\tmidLng:%v\n", midLat, midLng)
			fmt.Println("\tdeGeo")
			deGeo = GeoHashDecode(geo)
			fmt.Fprintf(os.Stdout, "\t\tmaxLat:%v\tminLat:%v\n", deGeo.MaxLat, deGeo.MinLat)
			fmt.Fprintf(os.Stdout, "\t\tminLng:%v\tmaxLng:%v\n", deGeo.MinLng, deGeo.MaxLng)
			fmt.Fprintf(os.Stdout, "\t\tmidLat:%v\tmidLng:%v\n", midLat, midLng)

		}
	}
	testEnd()
}

func TestDiffLatLng(t *testing.T) {
	testStart()

	distance := 54392.02
	fmt.Fprintf(os.Stdout, "distance=%v\n", distance)
	var angle float64 = 33
	for index, point := range Testpoints {
		ang := float64(index) * angle
		p := PointAtDistAndAngle(point, distance, ang)
		dist := EarthDistance(p, point)
		fmt.Fprintf(os.Stdout, "oriLat:%f\toriLng:%f\tdiff:%v\tdist:%f\tangle=%v\n", point.Lat, point.Lng, p, dist, ang)
	}
	testEnd()
}

func TestMidPoint(t *testing.T) {
	testStart()
	point := GeoPoint{Lat: 39.43373712, Lng: 120.378810983558}
	for _, p := range Testpoints {
		midPoint := MidPoint(point, p)
		dist1 := EarthDistanceOld(point, p)
		dist2 := EarthDistanceOld(point, midPoint)
		dist3 := EarthDistanceOld(p, midPoint)
		fmt.Fprintf(os.Stdout, "allDist:%f\tdist1:%f\tdist2:%f\n", dist1, dist2, dist3)
	}
	testEnd()
}

func TestFormatDistance(t *testing.T) {
	testStart()
	str := []string{
		"5000.00m",
		"44.1km",
		"43535345",
		"km222",
		"44km.01",
		"44.44KM",
	}
	for _, s := range str {
		ret := FormatDistance(s)
		fmt.Fprintf(os.Stdout, "str:%s\tdist:%v\n", s, ret)
	}
	testEnd()
}

func TestGeoGridBox_InGridBox(t *testing.T) {
	testStart()
	box := GeoRectangle{
		MaxLng: 120.378810983558,
		MinLng: 120.309268963407,
		MaxLat: 31.547996,
		MinLat: 31.047996,
	}
	fmt.Println(box.Width())
	fmt.Println(box.Height())
	fmt.Println(box.IsPointInRect(GeoPoint{Lat: 31.547996, Lng: 120.378810}))
	fmt.Println(box.IsPointInRect(GeoPoint{Lat: 31.947996, Lng: 120.378810}))
	fmt.Println(box.IsPointInRect(GeoPoint{Lat: 31.547996, Lng: 120.978810}))
	fmt.Println(box.IsPointInRect(GeoPoint{Lat: 31.947996, Lng: 120.978810}))
	fmt.Println(box.IsPointInRect(box.MidPoint()))
	testEnd()
}

// 点是否在多边形内
func TestGeoPolygon_IsPointInPolygon(t *testing.T) {
	testStart()
	polygon := getSpecialPolygon4()
	is10 := polygon.IsPointInPolygon(GeoPoint{Lat: 39.9957275390625, Lng: 116.3177490234375})
	fmt.Println(is10)
	polygon.AddPoint(GeoPoint{Lat: 39.972907, Lng: 116.322631})
	polygon.AddPoint(GeoPoint{Lat: 39.937953, Lng: 116.346777})
	polygon.AddPoint(GeoPoint{Lat: 39.902095, Lng: 116.322056})
	polygon.AddPoint(GeoPoint{Lat: 39.883051, Lng: 116.39622})
	polygon.AddPoint(GeoPoint{Lat: 39.96583, Lng: 116.488206})
	polygon.AddPoint(GeoPoint{Lat: 39.992368, Lng: 116.436464})
	polygon.AddPoint(GeoPoint{Lat: 39.972907, Lng: 116.322631})
	is1 := polygon.IsPointInPolygon(GeoPoint{Lat: 39.946804, Lng: 116.383572})
	fmt.Println(is1)
	is2 := polygon.IsPointInPolygon(GeoPoint{Lat: 40.043204, Lng: 116.394495})
	fmt.Println(is2)
	is3 := polygon.IsPointInPolygon(GeoPoint{Lat: 39.992368, Lng: 116.436464})
	fmt.Println(is3)
	testEnd()
}

type testLineIntersect struct {
	name  string   //地点名称
	point GeoPoint //这是一个点
}

func adfadfasdfasdfajlkjkhkjh(ps ...testLineIntersect) bool {
	pl := len(ps)
	for i := 0; i < pl; i++ {
		for j := i + 1; j < pl; j++ {
			if ps[i].point.IsEqual(ps[j].point) {
				return false
			}
		}
	}
	return true
}

func TestGeoLine_IsIntersect(t *testing.T) {
	testStart()
	points := []testLineIntersect{
		{name: "北大东门", point: GeoPoint{Lat: 39.998006, Lng: 116.322415}},
		{name: "西郊线香山站", point: GeoPoint{Lat: 39.999899, Lng: 116.211025}},
		{name: "宋家庄地铁站", point: GeoPoint{Lat: 39.851969, Lng: 116.434991}},
		{name: "首经贸地铁站", point: GeoPoint{Lat: 39.850391, Lng: 116.326583}},
		{name: "望京地铁站", point: GeoPoint{Lat: 40.004818, Lng: 116.474552}},
	}
	for _, p1 := range points {
		for _, p2 := range points {
			for _, p3 := range points {
				for _, p4 := range points {
					if !adfadfasdfasdfajlkjkhkjh(p1, p2, p3, p4) {
						continue
					}
					line1 := MakeGeoLine(p1.point, p2.point)
					line2 := MakeGeoLine(p4.point, p3.point)
					inter, _ := line1.IsIntersectWithLine(line2)
					fmt.Fprintf(
						os.Stdout,
						"是否相交[%v] line1[%s(%v,%v) %s(%v,%v)] line2[%s(%v,%v) %s(%v,%v)]\n",
						inter,
						p1.name, p1.point.Lat, p1.point.Lng,
						p2.name, p2.point.Lat, p2.point.Lng,
						p3.name, p3.point.Lat, p3.point.Lng,
						p4.name, p4.point.Lat, p4.point.Lng,
					)
				}
			}
		}
	}

	testEnd()
}

func TestGeoLine_GetIntersectPoint(t *testing.T) {
	testStart()
	p1 := GeoPoint{Lat: 39.998006, Lng: 116.322415}
	p2 := GeoPoint{Lat: 39.851969, Lng: 116.434991}
	p3 := GeoPoint{Lat: 40.004818, Lng: 116.474552}
	p4 := GeoPoint{Lat: 39.850391, Lng: 116.326583}
	line1 := MakeGeoLine(p1, p2)
	line2 := MakeGeoLine(p3, p4)
	a, p, b := line1.GetIntersectPoint(line2)
	fmt.Println(a, p, b)
	p1 = GeoPoint{Lat: 39.998006, Lng: 116.322415}
	p2 = GeoPoint{Lat: 39.851969, Lng: 116.434991}
	p3 = GeoPoint{Lat: 39.908006, Lng: 116.322415}
	p4 = GeoPoint{Lat: 39.908006, Lng: 116.922415}
	line1 = MakeGeoLine(p1, p2)
	line2 = MakeGeoLine(p3, p4)
	a, p, b = line1.GetIntersectPoint(line2)
	fmt.Println(a, p, b)
	testEnd()
}

func TestGeoLine_IsContainPoint(t *testing.T) {
	testStart()
	basePoint := GeoPoint{Lat: 39.983855, Lng: 116.27635}
	p1 := PointAtDistAndAngle(basePoint, 3000, 88)
	p2 := PointAtDistAndAngle(basePoint, 9000, 88)
	line := MakeGeoLine(p1, p2)
	fmt.Println(line.IsContainPoint(PointAtDistAndAngle(basePoint, 9000, 88)))
	fmt.Println(line.IsContainPoint(PointAtDistAndAngle(basePoint, 8888, 88)))
	fmt.Println(line.IsContainPoint(PointAtDistAndAngle(basePoint, 10000, 88)))
	fmt.Println(line.IsContainPoint(GeoPoint{Lat: 39.783855, Lng: 116.27635}))
	testEnd()
}

func TestGeoPloygon_FormatStringArray(t *testing.T) {
	testStart()
	ployon := GetTestPolygon1()
	s := ployon.FormatStringArray()
	for _, ss := range s {
		fmt.Println(ss)
	}

	testEnd()
}

func TestGeoHashDecodeBits(t *testing.T) {
	testStart()
	lat := 39.956981
	lng := 116.440488
	var i uint8
	for i = 1; i < 32; i++ {
		geo, rect := GeoHashBitsEncode(lat, lng, i)
		rect1 := GeoHashBitsDecode(geo, i)
		fmt.Fprintf(os.Stdout,
			"precision=%d\tdist=%v x %v\tdist=%v x %v\n",
			i, int(rect.Width()), int(rect.Height()), int(rect1.Width()), int(rect1.Height()))
	}
	testEnd()
}

func TestGeoHashBitsNeighbors(t *testing.T) {
	testStart()
	lat := 39.956981
	lng := 116.440488
	var pre uint8 = 15
	nebs := GeoHashBitsNeighbors(lat, lng, pre)
	for _, n := range nebs {
		rect := GeoHashBitsDecode(n, pre)
		fmt.Fprintf(os.Stdout, "dist=%v x %v\n", int(rect.Width()), int(rect.Height()))
		fmt.Println(rect.GetRectVertex())
	}
	testEnd()
}

func TestGeo_JustTest(t *testing.T) {
	testStart()

	testEnd()
}
