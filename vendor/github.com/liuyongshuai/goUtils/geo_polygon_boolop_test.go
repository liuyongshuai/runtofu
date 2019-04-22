// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-10-09 14:03

package goUtils

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"testing"
)

//交集及并集测试，要测试下面情况：
//目前已知：只要有边完全重合时，计算有问题
func TestGeoPolygon_BoolOperation(t *testing.T) {
	testStart()

	//一坨多边形的对对
	funcs := []func() ([2]GeoPolygon, string){
		getPolygonPairs1,      //两多边形完全没有交集
		getPolygonPairs2,      //两多边形只有部分顶点重合，但无交集
		getPolygonPairs3,      //两多边形有边重合情况，但无交集
		getPolygonPairs4,      //一个多边形完全包围另一个多边形，但边无相交情况
		getPolygonPairs5,      //一个多边形完全包围另一个多边形，有边重合情况
		getPolygonPairs6,      //两多边形有交集，但交集只有一个多边形，但边没有重合情况
		getPolygonPairs7,      //两多边形有交集，但交集只有一个多边形，边有重合的情况
		getPolygonPairs8,      //两多边形有交集，但交集是多个多边形，但边没有重合情况
		getPolygonPairs9,      //两多边形有交集，但交集是多个多边形，边有重合的情况
		getPolygonPairs10,     //两多边形完全相同
		getPolygonPairs11,     //两多边形有交集，但交集是多个多边形，边有部分重合的情况
		getPolygonPairs12,     //两个相邻但无交集的矩形
		getPolygonPairsRandom, //随机生成的多边形
	}
	titleFormat := "%s<br/><span style=\"color:red;font-weight:bold;\">（%s）</span>"
	for idx, f := range funcs {
		polys, title := f()
		poly1 := polys[0]
		poly2 := polys[1]
		fmt.Fprintf(os.Stdout, "start process polygon %d\n", idx+1)
		//如果多边形是空的
		if len(poly1.Points)*len(poly2.Points) <= 0 {
			fmt.Fprintf(os.Stdout, "skip index %d\n", idx+1)
			continue
		}
		//求出所有顶点的重心，计算地图显示的中心点
		var midPoint GeoPoint
		var latSum, lngSum float64
		for _, p := range poly1.Points {
			latSum += p.Lat
			lngSum += p.Lng
		}
		for _, p := range poly2.Points {
			latSum += p.Lat
			lngSum += p.Lng
		}
		pointNum := len(poly1.Points) + len(poly2.Points)
		midPoint.Lat = latSum / float64(pointNum)
		midPoint.Lng = lngSum / float64(pointNum)
		//交集多边形
		polygons := poly1.IntersectionWithPoly(poly2)
		drawPolygonBoolOperationInMap(
			fmt.Sprintf("%d_intersectionWithPolygon", idx+1),
			midPoint,
			13,
			poly1,
			poly2,
			polygons,
			fmt.Sprintf(titleFormat, title, "求交集结果"),
		)
		//并集多边形
		polygons = poly1.UnionWithPoly(poly2)
		drawPolygonBoolOperationInMap(
			fmt.Sprintf("%d_unionWithPolygon", idx+1),
			midPoint,
			13,
			poly1,
			poly2,
			polygons,
			fmt.Sprintf(titleFormat, title, "求并集结果"),
		)
		//差集多边形
		polygons = poly1.DifferenceWithPoly(poly2)
		drawPolygonBoolOperationInMap(
			fmt.Sprintf("%d_differenceWithPolygon", idx+1),
			midPoint,
			13,
			poly1,
			poly2,
			polygons,
			fmt.Sprintf(titleFormat, title, "求差集结果"),
		)
	}
	testEnd()
}

//两多边形完全没有交集
func getPolygonPairs1() (ret [2]GeoPolygon, title string) {
	title = "两多边形完全没有交集"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 115.860107, Lat: 40.11142},
		{Lng: 115.941745, Lat: 40.228742},
		{Lng: 116.137216, Lat: 40.153786},
		{Lng: 116.073976, Lat: 40.05931},
		{Lng: 116.034882, Lat: 40.141432},
		{Lng: 115.949794, Lat: 40.021306},
		{Lng: 115.830212, Lat: 40.04252},
		{Lng: 115.971641, Lat: 40.152903},
		{Lng: 115.860107, Lat: 40.11142},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.449971, Lat: 39.820325},
		{Lng: 116.52011, Lat: 39.985049},
		{Lng: 116.746627, Lat: 39.958508},
		{Lng: 116.802969, Lat: 39.822985},
		{Lng: 116.602898, Lat: 39.790171},
		{Lng: 116.617846, Lat: 39.898309},
		{Lng: 116.540807, Lat: 39.911593},
		{Lng: 116.45227, Lat: 39.820325},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两多边形只有部分顶点重合，但无交集
func getPolygonPairs2() (ret [2]GeoPolygon, title string) {
	title = "两多边形只有部分顶点重合，但无交集"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.277496, Lat: 39.965144},
		{Lng: 116.341311, Lat: 40.023958},
		{Lng: 116.470667, Lat: 39.998758}, //重合点
		{Lng: 116.379256, Lat: 39.980626},
		{Lng: 116.439047, Lat: 39.960278}, //重合点
		{Lng: 116.35281, Lat: 39.957623},
		{Lng: 116.38443, Lat: 39.922218}, //重合点
		{Lng: 116.279796, Lat: 39.900523},
		{Lng: 116.278646, Lat: 39.966471},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.548856, Lat: 40.006275},
		{Lng: 116.470667, Lat: 39.998758},
		{Lng: 116.525859, Lat: 39.958508},
		{Lng: 116.439047, Lat: 39.960278},
		{Lng: 116.514936, Lat: 39.915135},
		{Lng: 116.38443, Lat: 39.922218},
		{Lng: 116.445371, Lat: 39.847807},
		{Lng: 116.548281, Lat: 39.897423},
		{Lng: 116.555755, Lat: 39.955411},
		{Lng: 116.548281, Lat: 40.007601},
		{Lng: 116.548281, Lat: 40.007159},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两多边形有边重合情况，但无交集
func getPolygonPairs3() (ret [2]GeoPolygon, title string) {
	title = "两多边形有边重合情况，但无交集"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.314187, Lat: 39.981708},
		{Lng: 116.425146, Lat: 39.97817}, //以下6个点重合
		{Lng: 116.415372, Lat: 39.958262},
		{Lng: 116.353281, Lat: 39.956935},
		{Lng: 116.347532, Lat: 39.930383},
		{Lng: 116.403299, Lat: 39.920644},
		{Lng: 116.384902, Lat: 39.884334},
		{Lng: 116.307288, Lat: 39.891863},
		{Lng: 116.318786, Lat: 39.950298},
		{Lng: 116.261295, Lat: 39.980381},
		{Lng: 116.334884, Lat: 39.95959},
		{Lng: 116.314187, Lat: 39.98215},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.494136, Lat: 39.976843},
		{Lng: 116.425146, Lat: 39.97817},
		{Lng: 116.415372, Lat: 39.958262},
		{Lng: 116.353281, Lat: 39.956935},
		{Lng: 116.347532, Lat: 39.930383},
		{Lng: 116.403299, Lat: 39.920644},
		{Lng: 116.384902, Lat: 39.884334},
		{Lng: 116.475738, Lat: 39.871488},
		{Lng: 116.440094, Lat: 39.939677},
		{Lng: 116.493561, Lat: 39.977285},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//一个多边形完全包围另一个多边形，但边无相交情况
func getPolygonPairs4() (ret [2]GeoPolygon, title string) {
	title = "一个多边形完全包围另一个多边形，但边无相交情况"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.292915, Lat: 39.937907},
		{Lng: 116.330285, Lat: 39.9879},
		{Lng: 116.486087, Lat: 39.964014},
		{Lng: 116.441243, Lat: 39.853322},
		{Lng: 116.395825, Lat: 39.925956},
		{Lng: 116.347532, Lat: 39.869716},
		{Lng: 116.349832, Lat: 39.941005},
		{Lng: 116.269919, Lat: 39.890092},
		{Lng: 116.247497, Lat: 39.94543},
		{Lng: 116.29349, Lat: 39.937907},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.340058, Lat: 39.975073},
		{Lng: 116.292915, Lat: 39.925514},
		{Lng: 116.373978, Lat: 39.95251},
		{Lng: 116.36248, Lat: 39.907362},
		{Lng: 116.407323, Lat: 39.937907},
		{Lng: 116.434344, Lat: 39.891863},
		{Lng: 116.454467, Lat: 39.947643},
		{Lng: 116.406749, Lat: 39.969765},
		{Lng: 116.340633, Lat: 39.974189},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//一个多边形完全包围另一个多边形，有边重合情况
func getPolygonPairs5() (ret [2]GeoPolygon, title string) {
	title = "一个多边形完全包围另一个多边形，有边重合情况"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.292915, Lat: 39.937907}, //5
		{Lng: 116.330285, Lat: 39.9879},   //1
		{Lng: 116.486087, Lat: 39.964014}, //2
		{Lng: 116.441243, Lat: 39.853322},
		{Lng: 116.395825, Lat: 39.926399},
		{Lng: 116.347532, Lat: 39.869716},
		{Lng: 116.349832, Lat: 39.941005}, //4
		{Lng: 116.269919, Lat: 39.890092},
		{Lng: 116.247497, Lat: 39.94543}, //3
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.318212, Lat: 39.971535}, //1
		{Lng: 116.292915, Lat: 39.937907},
		{Lng: 116.247497, Lat: 39.94543},
		{Lng: 116.269919, Lat: 39.890092},
		{Lng: 116.349832, Lat: 39.941005},
		{Lng: 116.347532, Lat: 39.869716},
		{Lng: 116.395825, Lat: 39.926399},
		{Lng: 116.441243, Lat: 39.853322},
		{Lng: 116.486087, Lat: 39.964014}, //2
		{Lng: 116.31159316887657, Lat: 39.9628943614302},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两多边形有交集，但交集只有一个多边形，但边没有重合情况
func getPolygonPairs6() (ret [2]GeoPolygon, title string) {
	title = "两多边形有交集，但交集只有一个多边形，但边没有重合情况"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.275668, Lat: 39.954723},
		{Lng: 116.357306, Lat: 40.030343},
		{Lng: 116.504484, Lat: 39.992322},
		{Lng: 116.545303, Lat: 39.915332},
		{Lng: 116.442968, Lat: 39.985688},
		{Lng: 116.317062, Lat: 39.902934},
		{Lng: 116.276243, Lat: 39.954723},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.483787, Lat: 40.035204},
		{Lng: 116.428595, Lat: 39.894963},
		{Lng: 116.59877, Lat: 39.910904},
		{Lng: 116.570599, Lat: 40.00205},
		{Lng: 116.484362, Lat: 40.034762},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两多边形有交集，但交集只有一个多边形，边有重合的情况
func getPolygonPairs7() (ret [2]GeoPolygon, title string) {
	title = "两多边形有交集，但交集只有一个多边形，边有重合的情况"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.275668, Lat: 39.954723},
		{Lng: 116.357306, Lat: 40.030343},
		{Lng: 116.504484, Lat: 39.992322},
		{Lng: 116.545303, Lat: 39.915332},
		{Lng: 116.442968, Lat: 39.985688},
		{Lng: 116.317062, Lat: 39.902934},
		{Lng: 116.276243, Lat: 39.954723},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.317062, Lat: 39.902934},
		{Lng: 116.442968, Lat: 39.985688},
		{Lng: 116.450442, Lat: 39.944103},
		{Lng: 116.552777, Lat: 39.9764},
		{Lng: 116.442393, Lat: 39.883005},
		{Lng: 116.317637, Lat: 39.903377},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两多边形有交集，但交集是多个多边形，但边没有重合情况
func getPolygonPairs8() (ret [2]GeoPolygon, title string) {
	title = "两多边形有交集，但交集是多个多边形，但边没有重合情况"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.275668, Lat: 39.954723},
		{Lng: 116.357306, Lat: 40.030343},
		{Lng: 116.504484, Lat: 39.992322},
		{Lng: 116.545303, Lat: 39.915332},
		{Lng: 116.442968, Lat: 39.985688},
		{Lng: 116.317062, Lat: 39.902934},
		{Lng: 116.276243, Lat: 39.954723},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.261295, Lat: 40.00205},
		{Lng: 116.423996, Lat: 39.88832},
		{Lng: 116.314762, Lat: 40.036088},
		{Lng: 116.472864, Lat: 39.92817},
		{Lng: 116.334884, Lat: 40.048461},
		{Lng: 116.248072, Lat: 40.043159},
		{Lng: 116.26187, Lat: 40.000724},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两多边形有交集，但交集是多个多边形，边有重合的情况
func getPolygonPairs9() (ret [2]GeoPolygon, title string) {
	title = "两多边形有交集，但交集是多个多边形，边有重合的情况"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.275668, Lat: 39.954723},
		{Lng: 116.357306, Lat: 40.030343},
		{Lng: 116.504484, Lat: 39.992322},
		{Lng: 116.545303, Lat: 39.915332},
		{Lng: 116.442968, Lat: 39.985688},
		{Lng: 116.317062, Lat: 39.902934},
		{Lng: 116.276243, Lat: 39.954723},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.26302, Lat: 40.004261},
		{Lng: 116.317062, Lat: 39.902934},
		{Lng: 116.29694, Lat: 40.033436},
		{Lng: 116.442968, Lat: 39.985688},
		{Lng: 116.289466, Lat: 40.054205},
		{Lng: 116.26302, Lat: 40.004261},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两多边形完全相同
func getPolygonPairs10() (ret [2]GeoPolygon, title string) {
	title = "两多边形完全相同"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.275668, Lat: 39.954723},
		{Lng: 116.357306, Lat: 40.030343},
		{Lng: 116.504484, Lat: 39.992322},
		{Lng: 116.545303, Lat: 39.915332},
		{Lng: 116.442968, Lat: 39.985688},
		{Lng: 116.317062, Lat: 39.902934},
		{Lng: 116.276243, Lat: 39.954723},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.275668, Lat: 39.954723},
		{Lng: 116.357306, Lat: 40.030343},
		{Lng: 116.504484, Lat: 39.992322},
		{Lng: 116.545303, Lat: 39.915332},
		{Lng: 116.442968, Lat: 39.985688},
		{Lng: 116.317062, Lat: 39.902934},
		{Lng: 116.276243, Lat: 39.954723},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两多边形有交集，但交集是多个多边形，边有部分重合的情况
func getPolygonPairs11() (ret [2]GeoPolygon, title string) {
	title = "两多边形有交集，但交集是多个多边形，边有部分重合的情况"
	poly1 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.271643, Lat: 39.963129},
		{Lng: 116.429745, Lat: 40.045368},
		{Lng: 116.454467, Lat: 40.021945},
		{Lng: 116.403874, Lat: 39.988784},
		{Lng: 116.53323, Lat: 39.986131},
		{Lng: 116.382027, Lat: 39.979939},
		{Lng: 116.438369, Lat: 39.912233},
		{Lng: 116.346382, Lat: 39.969765},
		{Lng: 116.319936, Lat: 39.890535},
		{Lng: 116.319361, Lat: 39.958705},
		{Lng: 116.262445, Lat: 39.92994},
	})
	poly2 := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.272793, Lat: 39.962687},
		{Lng: 116.364205, Lat: 39.884777},
		{Lng: 116.304414, Lat: 39.962244},
		{Lng: 116.480913, Lat: 39.946315},
		{Lng: 116.478038, Lat: 40.014872},
		{Lng: 116.35064638850264, Lat: 40.00422373420368},
		{Lng: 116.271643, Lat: 39.963129},
	})
	ret = [2]GeoPolygon{poly1, poly2}
	return
}

//两个相邻没有交集的矩形，
func getPolygonPairs12() (ret [2]GeoPolygon, title string) {
	title = "两个相邻没有交集的矩形"
	rect1 := GeoRectangle{
		MinLat: 39.945034,
		MinLng: 116.308833,
		MaxLng: 116.410018,
		MaxLat: 39.999885,
	}
	//在经度上东延一点
	rect2 := GeoRectangle{
		MaxLat: rect1.MaxLat,
		MaxLng: rect1.MaxLng + 0.1,
		MinLng: rect1.MaxLng,
		MinLat: rect1.MinLat,
	}
	ret = [2]GeoPolygon{rect1.ToPolygon(), rect2.ToPolygon()}
	return
}

//随机生成的多边形
func getPolygonPairsRandom() (ret [2]GeoPolygon, title string) {
	randomPolygonList := GenPolygons(GeoRectangle{
		MaxLat: 40.033261,
		MinLat: 39.822564,
		MaxLng: 116.554322,
		MinLng: 116.190975,
	}, 2, 3, 20)
	title = "随机生成的多边形"
	ret = [2]GeoPolygon{randomPolygonList[0], randomPolygonList[1]}
	return
}

//在地图上画格子及多边形
func drawPolygonBoolOperationInMap(
	htmlName string,
	midPoint GeoPoint,
	level int,
	polygon1 GeoPolygon,
	polygon2 GeoPolygon,
	polygons []GeoPolygon,
	htmlTitle string,
) {
	if len(htmlTitle) <= 0 {
		htmlTitle = "多边形布尔操作效果观察"
	}
	//输出到模板
	buf := bytes.Buffer{}
	tpl := template.New("polygon")
	_, err := tpl.Parse(geoPolygonBoolOperationHtmlTemplate)
	if err != nil {
		panic(err)
	}
	point1 := polygon1.Points
	point2 := polygon2.Points
	var points [][]GeoPoint
	for _, p := range polygons {
		points = append(points, p.Points)
	}
	err = tpl.Execute(&buf, map[string]interface{}{
		"htmlTitle": htmlTitle,
		"tplName":   htmlName,
		"midPoint":  midPoint,
		"mapLevel":  level,
		"points1":   point1,
		"points2":   point2,
		"points3":   points,
	})
	if err != nil {
		panic(err)
	}

	//输出html文件的目录
	outDir := "./log/drawPolygonBoolOperationInMap"
	os.MkdirAll(outDir, 0755)
	outHtmlFile := fmt.Sprintf("%s/%s.html", outDir, htmlName)
	htmlFP, err1 := os.Create(outHtmlFile) //创建文件
	if err1 != nil {
		panic(err1)
	}

	//写入到相应的html文件里
	_, err1 = io.WriteString(htmlFP, buf.String())
	if err1 != nil {
		fmt.Println(err1)
	}
}

//画多边形的模板
var geoPolygonBoolOperationHtmlTemplate = `
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>{{ .htmlTitle }}</title>
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

		{{/*地图的可视区域*/}}
		var viewRect = polygonMap.getBounds();
		var southWest = viewRect.getSouthWest();
		var northEast = viewRect.getNorthEast();

		{{/*显示文字信息的label*/}}
		polygonMap.addOverlay(new BMap.Label(
			"{{ .htmlTitle }}",
			{ position: new BMap.Point(southWest.lng+(northEast.lng-southWest.lng)/15,northEast.lat)}
		));

		var polygonPoints = [];

		{{/*绘制多边形1*/}}
		polygonPoints = [];
		{{ range .points1 }}
			polygonPoints.push(new BMap.Point({{ .Lng }},{{ .Lat }}));
		{{ end }}
		{{/*设置多边形的显示属性，红边等*/}}
        var polygonObject = new BMap.Polygon(polygonPoints);
        polygonObject.setStrokeColor("red");
		polygonObject.setStrokeWeight('2');
		polygonObject.setFillColor('transparent');
        polygonMap.addOverlay(polygonObject);

		{{/*绘制多边形2*/}}
		polygonPoints = [];
		{{ range .points2 }}
			polygonPoints.push(new BMap.Point({{ .Lng }},{{ .Lat }}));
		{{ end }}
		{{/*设置多边形的显示属性，红边等*/}}
        var polygonObject = new BMap.Polygon(polygonPoints);
        polygonObject.setStrokeColor("blue");
		polygonObject.setStrokeWeight('2');
		polygonObject.setFillColor('transparent');
        polygonMap.addOverlay(polygonObject);


		{{/*绘制交集/并集多边形*/}}
		{{ range .points3 }}
			polygonPoints = [];
			{{ range .}}
				polygonPoints.push(new BMap.Point({{ .Lng }},{{ .Lat }}));
			{{ end }}
			var polygonObject = new BMap.Polygon(polygonPoints);
			polygonObject.setStrokeColor("green");
			polygonObject.setStrokeStyle('dashed');
			polygonObject.setFillColor('transparent');
			polygonMap.addOverlay(polygonObject);
		{{ end }}
	})();
</script>
`
