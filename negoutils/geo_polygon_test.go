/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date        2018-05-23 15:37
 */
package negoutils

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"math"
	"os"
	"testing"
	"time"
)

var customDeliveryGeoHashPrecision = 6
var testSplitPolygonGeoHashPrecision = customDeliveryGeoHashPrecision
var splitPolygonRectGeoHashTypeForWendao = GEOHASH_TYPE_NORMAL

func TestGeoPolygon_SplitGeoHashRect(t *testing.T) {
	testStart()
	var polygon GeoPolygon
	testSplitPolygonGeoHashPrecision = 6
	isRay := true //采用射线法切割多边形
	polygon = getSpecialPolygon1()
	polygon.SetGeoHashType(splitPolygonRectGeoHashTypeForWendao)
	splitGeoHashRect(polygon, "polygon_special1", 14, isRay)
	polygon = getSpecialPolygon2()
	polygon.SetGeoHashType(splitPolygonRectGeoHashTypeForWendao)
	splitGeoHashRect(polygon, "polygon_special2", 14, isRay)
	polygon = getSpecialPolygon3()
	polygon.SetGeoHashType(splitPolygonRectGeoHashTypeForWendao)
	splitGeoHashRect(polygon, "polygon_special3", 14, isRay)
	polygon = getSpecialPolygon4()
	polygon.SetGeoHashType(splitPolygonRectGeoHashTypeForWendao)
	splitGeoHashRect(polygon, "polygon_special4", 14, isRay)
	polygon = getSpecialPolygon5()
	polygon.SetGeoHashType(splitPolygonRectGeoHashTypeForWendao)
	splitGeoHashRect(polygon, "polygon_special5", 14, isRay)
	polygon = getSpecialPolygon6()
	polygon.SetGeoHashType(splitPolygonRectGeoHashTypeForWendao)
	splitGeoHashRect(polygon, "polygon_special6", 14, isRay)
	return
	var polygons []GeoPolygon
	polygons = append(polygons,
		GetTestPolygon1(),
		getPolygon2(),
		getPolygon3(),
		getPolygon4(),
		getPolygon5(),
		getPolygon6(),
		getPolygon7(),
		getPolygon8(),
		getPolygon9(),
		getPolygon10(),
		getPolygon11(),
		getPolygon12(),
	)
	for i, polygon := range polygons {
		splitGeoHashRect(polygon, fmt.Sprintf("polygon%d", i+1), 13, isRay)
	}

	var total float64 = 0
	size := 10
	randomPolygonList := GenPolygons(GeoRectangle{
		MaxLat: 40.033261,
		MinLat: 39.822564,
		MaxLng: 116.554322,
		MinLng: 116.190975,
	}, size, 3, 20)
	for i, polygon := range randomPolygonList {
		total += splitGeoHashRect(polygon, fmt.Sprintf("polygon%d", i+len(polygons)+1), 13, isRay)
	}
	avg := total / float64(size)
	fmt.Println(avg)
	testEnd()
}

// 凸多边形：http://10.96.112.48/polygon1.html
func GetTestPolygon1() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.385297, Lat: 39.993252},
		{Lng: 116.325505, Lat: 39.974235},
		{Lng: 116.290435, Lat: 39.931314},
		{Lng: 116.346777, Lat: 39.879508},
		{Lng: 116.436464, Lat: 39.911836},
		{Lng: 116.451987, Lat: 39.93751},
		{Lng: 116.449687, Lat: 39.971138},
		{Lng: 116.415767, Lat: 39.994579},
	})
	return polygon
}

// http://10.96.112.48/polygon2.html
func getPolygon2() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.399669, Lat: 40.004307},
		{Lng: 116.360575, Lat: 39.952114},
		{Lng: 116.281812, Lat: 39.954326},
		{Lng: 116.3623, Lat: 39.916706},
		{Lng: 116.309983, Lat: 39.863559},
		{Lng: 116.401969, Lat: 39.892352},
		{Lng: 116.503729, Lat: 39.861344},
		{Lng: 116.469234, Lat: 39.929101},
		{Lng: 116.529025, Lat: 39.978215},
		{Lng: 116.440488, Lat: 39.956981},
	})
	return polygon
}

// http://10.96.112.48/polygon3.html
func getPolygon3() GeoPolygon {
	p := MakeGeoPoint(39.923664, 116.403424)
	var points []GeoPoint
	stp := 20
	num := 360 / stp
	for i := 0; i < num; i++ {
		dist := 4000
		if i%2 == 0 {
			dist = 8000
		}
		angle := float64(i * stp)
		p0 := PointAtDistAndAngle(p, float64(dist), angle)
		points = append(points, p0)
	}
	return MakeGeoPolygon(points)
}

// http://10.96.112.48/polygon4.html
func getPolygon4() GeoPolygon {
	//{39.869384765625 116.279296875 39.9957275390625 116.455078125}
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lat: 39.869384765625, Lng: 116.279296875},
		{Lat: 39.9957275390625, Lng: 116.279296875},
		{Lat: 39.9957275390625, Lng: 116.455078125},
		{Lat: 39.869384765625, Lng: 116.455078125},
	})
	return polygon

}

// http://10.96.112.48/polygon5.html
func getPolygon5() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.315013, Lat: 39.969147},
		{Lng: 116.340453, Lat: 39.993584},
		{Lng: 116.364456, Lat: 39.96771},
		{Lng: 116.37193, Lat: 39.967488},
		{Lng: 116.38429, Lat: 39.994358},
		{Lng: 116.411024, Lat: 39.964724},
		{Lng: 116.423241, Lat: 39.994247},
		{Lng: 116.463916, Lat: 39.966935},
		{Lng: 116.423816, Lat: 39.940387},
		{Lng: 116.441926, Lat: 39.91405},
		{Lng: 116.399382, Lat: 39.89202},
		{Lng: 116.350514, Lat: 39.912832},
		{Lng: 116.300209, Lat: 39.939281},
		{Lng: 116.342753, Lat: 39.955654},
	})
	return polygon
}

// http://10.96.112.48/polygon6.html
func getPolygon6() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.314438, Lat: 39.968926},
		{Lng: 116.331829, Lat: 39.991594},
		{Lng: 116.350227, Lat: 39.963949},
		{Lng: 116.369343, Lat: 39.992921},
		{Lng: 116.386159, Lat: 39.96406},
		{Lng: 116.42281, Lat: 39.9948},
		{Lng: 116.463485, Lat: 39.966493},
		{Lng: 116.496543, Lat: 39.95134},
		{Lng: 116.442069, Lat: 39.929543},
		{Lng: 116.485332, Lat: 39.913718},
		{Lng: 116.448681, Lat: 39.878068},
		{Lng: 116.425541, Lat: 39.91416},
		{Lng: 116.414617, Lat: 39.846499},
		{Lng: 116.390327, Lat: 39.896338},
		{Lng: 116.360144, Lat: 39.846499},
		{Lng: 116.33456, Lat: 39.886705},
		{Lng: 116.28713, Lat: 39.854808},
		{Lng: 116.319756, Lat: 39.90298},
		{Lng: 116.281668, Lat: 39.931093},
		{Lng: 116.331254, Lat: 39.952778},
		{Lng: 116.277356, Lat: 39.976446},
	})
	return polygon
}

// http://10.96.112.48/polygon7.html
func getPolygon7() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.30797, Lat: 39.991926},
		{Lng: 116.345627, Lat: 40.006739},
		{Lng: 116.385297, Lat: 39.993695},
		{Lng: 116.426978, Lat: 40.020664},
		{Lng: 116.451124, Lat: 39.993252},
		{Lng: 116.498267, Lat: 39.959857},
		{Lng: 116.467797, Lat: 39.952114},
		{Lng: 116.439051, Lat: 39.975562},
		{Lng: 116.32838, Lat: 39.972907},
		{Lng: 116.315444, Lat: 39.968041},
		{Lng: 116.319181, Lat: 39.873749},
		{Lng: 116.353964, Lat: 39.854919},
		{Lng: 116.462335, Lat: 39.864888},
		{Lng: 116.485907, Lat: 39.846721},
		{Lng: 116.46291, Lat: 39.801281},
		{Lng: 116.408006, Lat: 39.838078},
		{Lng: 116.349652, Lat: 39.78487},
		{Lng: 116.299634, Lat: 39.836527},
		{Lng: 116.234956, Lat: 39.911836},
		{Lng: 116.302509, Lat: 39.939281},
	})
	return polygon
}

// http://10.96.112.48/polygon8.html
func getPolygon8() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.328667, Lat: 39.972907},
		{Lng: 116.362012, Lat: 39.949238},
		{Lng: 116.441063, Lat: 39.947246},
		{Lng: 116.457161, Lat: 39.970475},
		{Lng: 116.465785, Lat: 39.874413},
		{Lng: 116.436177, Lat: 39.910508},
		{Lng: 116.364887, Lat: 39.906301},
		{Lng: 116.322056, Lat: 39.873749},
		{Lng: 116.34534, Lat: 39.930871},
	})
	return polygon
}

// http://10.96.112.48/polygon9.html
func getPolygon9() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.403685, Lat: 39.909262},
		{Lng: 116.40461, Lat: 39.909255},
		{Lng: 116.40461, Lat: 39.908543},
		{Lng: 116.403676, Lat: 39.908543},
	})
	return polygon
}

// http://10.96.112.48/polygon10.html
func getPolygon10() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.363126, Lat: 39.913468},
		{Lng: 116.363162, Lat: 39.912777},
		{Lng: 116.442465, Lat: 39.914188},
		{Lng: 116.442213, Lat: 39.915046},
	})
	return polygon
}

// http://10.96.112.48/polygon11.html
func getPolygon11() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.211097, Lat: 39.913607},
		{Lng: 116.584217, Lat: 39.913607},
		{Lng: 116.385871, Lat: 40.07943},
	})
	return polygon
}

// http://10.96.112.48/polygon12.html
func getPolygon12() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.305383, Lat: 39.991262},
		{Lng: 116.452562, Lat: 39.991262},
		{Lng: 116.452562, Lat: 39.966493},
		{Lng: 116.385122, Lat: 39.966493},
		{Lng: 116.385122, Lat: 39.889031},
		{Lng: 116.380132, Lat: 39.889031},
		{Lng: 116.380122, Lat: 39.966493},
		{Lng: 116.305383, Lat: 39.966493},
	})
	return polygon
}

// 切多边形的特殊case
func getSpecialPolygon1() GeoPolygon {
	stp := customDeliveryGeoHashPrecision
	var points []GeoPoint
	points = append(points, GeoPoint{Lng: 116.315426, Lat: 40.012642})
	_, bRect := GeoHashEncode(39.998016, 116.322289, stp)
	p1 := GeoPoint{Lng: bRect.MaxLng, Lat: bRect.MaxLat}
	points = append(points, p1)
	p2 := GeoPoint{Lng: bRect.MaxLng + bRect.LngSpan(), Lat: bRect.MaxLat}
	points = append(points, p2)
	p3 := GeoPoint{Lng: p2.Lng + bRect.LngSpan(), Lat: p2.Lat + bRect.LatSpan()}
	points = append(points, p3)
	p4 := GeoPoint{Lng: p3.Lng + bRect.LngSpan(), Lat: bRect.MinLat}
	points = append(points, p4)
	p5 := GeoPoint{Lng: p4.Lng + 2*bRect.LngSpan(), Lat: p3.Lat}
	points = append(points, p5)
	p6 := GeoPoint{Lng: p4.Lng + 2*bRect.LngSpan(), Lat: p4.Lat}
	points = append(points, p6)
	p7 := GeoPoint{Lng: p5.Lng + 2*bRect.LngSpan(), Lat: p5.Lat}
	points = append(points, p7)
	p8 := GeoPoint{Lng: p6.Lng, Lat: p6.Lat - 8*bRect.LatSpan()}
	points = append(points, p8)
	p9 := GeoPoint{Lng: p6.Lng - 8*bRect.LngSpan(), Lat: p8.Lat}
	points = append(points, p9)
	p10 := GeoPoint{Lng: p9.Lng, Lat: bRect.MaxLat}
	points = append(points, p10)
	return MakeGeoPolygon(points)
}

// 切多边形的特殊case
func getSpecialPolygon2() GeoPolygon {
	stp := customDeliveryGeoHashPrecision
	var points []GeoPoint
	_, bRect := GeoHashEncode(39.998016, 116.322289, stp)
	b := bRect.Width() / bRect.Height()
	points = append(points, bRect.RightBottomPoint())
	points = append(points, GeoPoint{
		Lng: bRect.MaxLng,
		Lat: bRect.MaxLat + 10*bRect.LatSpan(),
	})
	points = append(points, GeoPoint{
		Lat: bRect.MinLat,
		Lng: bRect.MaxLng - 10*b*bRect.LngSpan(),
	})
	return MakeGeoPolygon(points)
}

// 切多边形的特殊case
func getSpecialPolygon3() GeoPolygon {
	stp := customDeliveryGeoHashPrecision
	var points []GeoPoint
	geos := GetNeighborsGeoCodes(39.998016, 116.322289, stp)

	rightBottom := GeoHashDecode(geos[8])
	rightUp := GeoHashDecode(geos[7])
	leftBottom := GeoHashDecode(geos[6])
	points = append(points,
		rightUp.RightUpPoint(),
		leftBottom.LeftBottomPoint(),
		rightBottom.RightBottomPoint(),
	)
	return MakeGeoPolygon(points)
}

// 切多边形的特殊case
func getSpecialPolygon4() GeoPolygon {
	stp := customDeliveryGeoHashPrecision
	var points []GeoPoint
	geos := GetNeighborsGeoCodes(39.998016, 116.322289, stp)
	//center := GeoHashDecode(geos[0])
	centerUp := GeoHashDecode(geos[1])
	centerBottom := GeoHashDecode(geos[2])
	//leftCenter:=GeoHashDecode(geos[3] )
	//rightCenter:=GeoHashDecode(geos[4] )
	leftUp := GeoHashDecode(geos[5])
	leftBottom := GeoHashDecode(geos[6])
	rightUp := GeoHashDecode(geos[7])
	rightBottom := GeoHashDecode(geos[8])

	points = append(points,
		centerUp.RightBottomPoint(),
		centerUp.RightUpPoint(),
		rightUp.RightUpPoint(),
		rightBottom.RightBottomPoint(),
		centerBottom.RightBottomPoint(),
		centerBottom.RightUpPoint(),
		centerBottom.LeftUpPoint(),
		centerBottom.LeftBottomPoint(),
		leftBottom.LeftBottomPoint(),
		leftUp.LeftUpPoint(),
		leftUp.RightUpPoint(),
		leftUp.RightBottomPoint(),
	)
	return MakeGeoPolygon(points)
}

// 切多边形的特殊case
func getSpecialPolygon5() GeoPolygon {
	stp := customDeliveryGeoHashPrecision
	var points []GeoPoint
	_, tRect := GeoHashEncode(39.998016, 116.322289, stp)
	points = append(points, GeoPoint{Lat: 39.9957275390625, Lng: 116.31774915309339})
	p2 := GeoPoint{Lat: tRect.MinLat - tRect.LatSpan(), Lng: tRect.MaxLng}
	points = append(points, p2)
	p3 := GeoPoint{Lat: tRect.MaxLat, Lng: tRect.MaxLng + tRect.LngSpan()}
	points = append(points, p3)
	p4 := GeoPoint{Lat: tRect.MinLat - 3*tRect.LatSpan(), Lng: p3.Lng}
	points = append(points, p4)
	p5 := GeoPoint{Lat: p4.Lat, Lng: p4.Lng - 6*tRect.LngSpan()}
	points = append(points, p5)
	p6 := GeoPoint{Lat: p3.Lat, Lng: p3.Lng - 6*tRect.LngSpan()}
	points = append(points, p6)
	p7 := GeoPoint{Lat: tRect.MinLat - tRect.LatSpan(), Lng: tRect.MinLng}
	points = append(points, p7)
	return MakeGeoPolygon(points)
}

// 切多边形的特殊case
func getSpecialPolygon6() GeoPolygon {
	var points []GeoPoint
	var stp uint8 = 14
	geos := GeoHashBitsNeighbors(39.998016, 116.322289, stp)
	maxLat := MIN_LATITUDE
	minLat := MAX_LATITUDE
	maxLng := MIN_LONGITUDE
	minLng := MAX_LONGITUDE
	for _, geo := range geos {
		rect := GeoHashBitsDecode(geo, stp)
		maxLat = math.Max(maxLat, rect.MaxLat)
		minLat = math.Min(minLat, rect.MinLat)
		maxLng = math.Max(maxLng, rect.MaxLng)
		minLng = math.Min(minLng, rect.MinLng)
	}
	rect := GeoRectangle{MaxLat: maxLat, MinLat: minLat, MaxLng: maxLng, MinLng: minLng}
	points = rect.GetRectVertex()
	return MakeGeoPolygon(points)
}

// 将多边形及切格子后的画在地图上
func splitGeoHashRect(
	polygon GeoPolygon, //多边形
	htmlName string, //生成的html文件识别名称
	level int, //百度地图显示的级别
	isRay bool, //是射线切割法、还是暴力遍历法
) float64 {
	//多边形的顶点、中心点坐标
	polygonRect := polygon.GetBoundsRect()
	midPoint := polygonRect.MidPoint()
	polygonPoints := polygon.GetPoints()

	//用不同的方法切格子
	st := time.Now().UnixNano()
	var inGrids, pGrids []string
	if isRay {
		inGrids, pGrids = polygon.RaySplitGeoHashRect(testSplitPolygonGeoHashPrecision)
	} else {
		inGrids, pGrids = polygon.ViolentSplitGeoHashRect(testSplitPolygonGeoHashPrecision)
	}
	et := time.Now().UnixNano()
	diff := pT(st, et)

	//收集所有的小格子的经纬度坐标信息
	var inRectList [][]GeoPoint
	var broderRectList [][]GeoPoint
	for _, grid := range inGrids {
		rect := GeneralGeoHashDecode(grid, testSplitPolygonGeoHashPrecision, polygon.GeoHashType)
		inRectList = append(inRectList, rect.GetRectVertex())
	}
	for _, grid := range pGrids {
		rect := GeneralGeoHashDecode(grid, testSplitPolygonGeoHashPrecision, polygon.GeoHashType)
		broderRectList = append(broderRectList, rect.GetRectVertex())
	}
	drawPolygonAndGridInMap(
		htmlName,
		midPoint,
		level,
		polygonPoints,
		inRectList,
		broderRectList,
	)

	return diff
}

// 处理请求的时间proc_time
// 时间要求time.Now().UnixNano()
func pT(st, et int64) float64 {
	var t1 int64 = 0
	var t2 int64 = 0
	if st > et {
		t1 = et
		t2 = st
	} else {
		t1 = st
		t2 = et
	}
	ret := float64(t2-t1) / float64(1000*1000*1000)
	return ret
}

// 在地图上画格子及多边形
func drawPolygonAndGridInMap(
	htmlName string,
	midPoint GeoPoint,
	level int,
	polygonPoints []GeoPoint,
	inRectList [][]GeoPoint,
	broderRectList [][]GeoPoint,
) {
	//输出到模板
	buf := bytes.Buffer{}
	tpl := template.New("polygon")
	_, err := tpl.Parse(geoPolygonHtmlTemplate)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(&buf, map[string]interface{}{
		"tplName":        htmlName,
		"midPoint":       midPoint,
		"mapLevel":       level,
		"polygonPoints":  polygonPoints,
		"inRectList":     inRectList,
		"borderRectList": broderRectList,
	})
	if err != nil {
		panic(err)
	}

	//输出html文件的目录
	outDir := "./log/drawPolygonAndGridInMap"
	os.MkdirAll(outDir, 0755)
	outHtmlFile := fmt.Sprintf("%s/%s.html", outDir, htmlName)
	htmlFP, err1 := os.Create(outHtmlFile) //创建文件
	if err1 != nil {
		panic(err1)
	}

	//写入到相应的html文件里
	fmt.Println(htmlName, "inRect:", len(inRectList), "\tborderRect:", len(broderRectList), "\toutHtmlFile：", outHtmlFile)
	_, err1 = io.WriteString(htmlFP, buf.String())
	if err1 != nil {
		fmt.Println(err1)
	}
}

// 画多边形的模板
var geoPolygonHtmlTemplate = `
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>切格子效果观察</title>
		<script type="text/javascript" src="http://api.map.baidu.com/api?v=1.2"></script>
		<script type="text/javascript" src="http://api.map.baidu.com/library/GeoUtils/1.2/src/GeoUtils_min.js"></script>
	</head>
	<body>
		<div style="width:100%;height:100%;border:1px solid gray" id="container_{{ .tplName }}"></div>
	</body>
</html>
<script type="text/javascript">
	(function(){
		{{/*设置地图相关属性*/}}
        var polygonMap = new BMap.Map("container_{{ .tplName }}");
        polygonMap.addControl(new BMap.NavigationControl());
        polygonMap.addControl(new BMap.ScaleControl());
        polygonMap.addControl(new BMap.OverviewMapControl());
        polygonMap.addControl(new BMap.CopyrightControl());
		polygonMap.enableContinuousZoom();

		{{/*地图视图的中心点经纬度及视图级别*/}}
        polygonMap.centerAndZoom(new BMap.Point({{ .midPoint.Lng}},{{ .midPoint.Lat }}), {{ .mapLevel }});

		{{/*绘制多边形*/}}
        var polygonPoints = [];
		{{ range .polygonPoints }}
			polygonPoints.push(new BMap.Point({{ .Lng }},{{ .Lat }}));
		{{ end }}
		{{/*设置多边形的显示属性，红边等*/}}
        var polygonObject = new BMap.Polygon(polygonPoints);
        polygonObject.setStrokeColor("red");
        polygonMap.addOverlay(polygonObject);

		{{/*开始绘制切出来的小格子*/}}
		var gridPoints = [];
		var gridPolygonObject;

		{{/*完全被多边形包围的小格子*/}}
		{{ range .inRectList }}
			gridPoints = [];
			{{ range .}}
				 gridPoints.push(new BMap.Point({{ .Lng }},{{ .Lat }}));
			{{ end }}
			{{/*完全被包围的小格子边框为实线*/}}
        	gridPolygonObject=new BMap.Polygon(gridPoints);
        	gridPolygonObject.setStrokeWeight('1');
        	polygonMap.addOverlay(gridPolygonObject);
		{{ end }}

		{{/*部分被多边形包围的小格子*/}}
		{{ range .borderRectList }}
			gridPoints = [];
			{{ range .}}
				gridPoints.push(new BMap.Point({{ .Lng }},{{ .Lat }}));
			{{ end }}
			{{/*半包围的小格子虚线边框、填充色较浅*/}}
        	gridPolygonObject=new BMap.Polygon(gridPoints);
        	gridPolygonObject.setStrokeWeight('1');
        	gridPolygonObject.setStrokeStyle('dashed');
        	gridPolygonObject.setFillColor('#F0F8FF');
        	polygonMap.addOverlay(gridPolygonObject);
		{{ end }}
	})();
</script>
`

func TestGeoRectangle_GetRandomGeoPoint(t *testing.T) {
	/*beijingRec := GeoRectangle{
		MinLat: 39.67068,
		MaxLat: 41.060816,
		MinLng: 115.423411,
		MaxLng: 116.67892,
	}*/

	testPoly := GetTestPolygon1()
	testRect := testPoly.GetBoundsRect()

	randomPoints := make([]GeoPoint, 0)
	for i := 0; i < 800; i++ {
		point := testRect.GetRandomGeoPoint()
		randomPoints = append(randomPoints, point)
	}
	drawRectangleAndPointsInMap("TestGeoRectangle_GetRandomGeoPoint", 14, testRect, randomPoints)
}

// 在地图上画格子及多边形
func drawRectangleAndPointsInMap(
	htmlName string, //生成的html文件识别名称
	level int, //百度地图显示的级别
	rectangle GeoRectangle, // 矩形
	randomPoints []GeoPoint, // 随机生成的所有点
) {
	midPoint := rectangle.MidPoint()
	polygonPoints := rectangle.GetRectVertex()

	//输出到模板
	buf := bytes.Buffer{}
	tpl := template.New("polygon")
	_, err := tpl.Parse(geoRectanglePointHtmlTemplate)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(&buf, map[string]interface{}{
		"tplName":       htmlName,
		"midPoint":      midPoint,
		"mapLevel":      level,
		"polygonPoints": polygonPoints,
		"randomPoints":  randomPoints,
	})
	if err != nil {
		panic(err)
	}

	//输出html文件的目录
	outDir := "./log/drawRectangleAndPointsInMap"
	os.MkdirAll(outDir, 0755)
	outHtmlFile := fmt.Sprintf("%s/%s.html", outDir, htmlName)
	htmlFP, err1 := os.Create(outHtmlFile) //创建文件
	if err1 != nil {
		panic(err1)
	}

	//写入到相应的html文件里
	fmt.Println(htmlName, "randomPoints num:", len(randomPoints), "\toutHtmlFile：", outHtmlFile)
	_, err1 = io.WriteString(htmlFP, buf.String())
	if err1 != nil {
		fmt.Println(err1)
	}
}

// 画多边形的模板
var geoRectanglePointHtmlTemplate = `
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>矩形随机点效果观察</title>
		<script type="text/javascript" src="http://api.map.baidu.com/api?v=1.2"></script>
		<script type="text/javascript" src="http://api.map.baidu.com/library/GeoUtils/1.2/src/GeoUtils_min.js"></script>
	</head>
	<body>
		<div style="width:100%;height:100%;border:1px solid gray" id="container_{{ .tplName }}"></div>
	</body>
</html>
<script type="text/javascript">
	(function(){
		{{/*设置地图相关属性*/}}
        var polygonMap = new BMap.Map("container_{{ .tplName }}");
        polygonMap.addControl(new BMap.NavigationControl());
        polygonMap.addControl(new BMap.ScaleControl());
        polygonMap.addControl(new BMap.OverviewMapControl());
        polygonMap.addControl(new BMap.CopyrightControl());
		polygonMap.enableContinuousZoom();
		polygonMap.enableScrollWheelZoom(true);

		{{/*地图视图的中心点经纬度及视图级别*/}}
        polygonMap.centerAndZoom(new BMap.Point({{ .midPoint.Lng}},{{ .midPoint.Lat }}), {{ .mapLevel }});

		{{/*绘制多边形*/}}
        var polygonPoints = [];
		{{ range .polygonPoints }}
			polygonPoints.push(new BMap.Point({{ .Lng }},{{ .Lat }}));
		{{ end }}
		{{/*设置多边形的显示属性，红边等*/}}
        var polygonObject = new BMap.Polygon(polygonPoints);
        polygonObject.setStrokeColor("red");
        polygonMap.addOverlay(polygonObject);


		{{/*画出所有的点*/}}
		{{ range .randomPoints}}
			var marker = new BMap.Marker(new BMap.Point({{ .Lng }}, {{ .Lat }}));
        	polygonMap.addOverlay(marker);
		{{ end }}
	})();
</script>
`

func TestGeoRectangle_DrawGeoPointAndPolygon(t *testing.T) {

	// your special polygon
	pString := []string{
		"30.56384531551,104.03984606266",
		"30.531321387182,104.04053270817",
		"30.534869981526,104.10645067692",
		"30.570348793497,104.10851061344",
	}
	testRect := MakeGeoPolygonByStr(pString)

	randomPoints := make([]GeoPoint, 0)
	// your special point
	point := GeoPoint{
		Lat: 30.56636,
		Lng: 104.06354,
	}
	randomPoints = append(randomPoints, point)

	drawPolygonAndPointsInMap("TestGeoRectangle_DrawGeoPointAndPolygon", 14, testRect, randomPoints)
}

// 在地图上画格子及多边形
func drawPolygonAndPointsInMap(
	htmlName string, //生成的html文件识别名称
	level int, //百度地图显示的级别
	polygon GeoPolygon, // 矩形
	randomPoints []GeoPoint, // 随机生成的所有点
) {

	polygonPoints := polygon.Points
	midPoint := polygon.Points[0]

	//输出到模板
	buf := bytes.Buffer{}
	tpl := template.New("polygon")
	_, err := tpl.Parse(geoRectanglePointHtmlTemplate)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(&buf, map[string]interface{}{
		"tplName":       htmlName,
		"midPoint":      midPoint,
		"mapLevel":      level,
		"polygonPoints": polygonPoints,
		"randomPoints":  randomPoints,
	})
	if err != nil {
		panic(err)
	}

	//输出html文件的目录
	outDir := "./log/drawPolygonAndPointsInMap"
	os.MkdirAll(outDir, 0755)
	outHtmlFile := fmt.Sprintf("%s/%s.html", outDir, htmlName)
	htmlFP, err1 := os.Create(outHtmlFile) //创建文件
	if err1 != nil {
		panic(err1)
	}

	//写入到相应的html文件里
	fmt.Println(htmlName, "randomPoints num:", len(randomPoints), "\toutHtmlFile：", outHtmlFile)
	_, err1 = io.WriteString(htmlFP, buf.String())
	if err1 != nil {
		fmt.Println(err1)
	}
}

func TestQieGeZi(t *testing.T) {

	// your special polygon
	pString := []string{
		"30.56384531551,104.03984606266",
		"30.531321387182,104.04053270817",
		"30.534869981526,104.10645067692",
		"30.570348793497,104.10851061344",
	}
	polygon := MakeGeoPolygonByStr(pString)

	//polygon.SetGeoHashType(splitPolygonRectGeoHashTypeForWendao)
	//splitGeoHashRect(polygon, "TestQieGeZi", 14, true)

	point := GeoPoint{
		Lat: 30.56636,
		Lng: 104.06354,

		//Lat:30.53886,
		//Lng:104.06853,
	}

	println(polygon.IsPointInPolygon(point))
}
