// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-10-30 11:14

package goUtils

import (
	"context"
	"fmt"
	"github.com/kr/pretty"
	"net/http"
	"testing"
)

var (
	testUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
)

//纯GET请求
func TestHttpClient_Get(t *testing.T) {
	testStart()

	url := "http://10.96.114.84/add1.php?a=1&b=3"
	client := NewHttpClient(url, context.Background())
	client.AddCookie("c1", "v1")
	client.SetReferer("http://microsoft.com")
	client.AddHeader("myHeaderKey", "myHeaderValue")
	client.SetHost("test.wendao.com")
	client.SetKeepAlive(true)
	resp, _ := client.Get()
	fmt.Println(string(resp.GetBody()))
	resp, _ = client.Get()
	fmt.Println(string(resp.GetBody()))
	resp, _ = client.Post()
	fmt.Println(string(resp.GetBody()))

	testEnd()
}

//POST几个字段信息
func TestHttpClient_Post(t *testing.T) {
	testStart()

	url := "http://10.96.114.84/add1.php"
	client := NewHttpClient(url, context.Background())
	client.SetField("singleValue", "value")
	client.AddField("multiValue", "v1")
	client.AddField("multiValue", "v2")
	client.AddField("multiValue", "v3")
	client.AddCookie("c1", "v1")
	client.SetUserAgent(testUserAgent)
	client.SetReferer("http://microsoft.com")
	client.AddHeader("myHeaderKey", "myHeaderValue")
	client.SetHost("test.wendao.com")
	resp, _ := client.Post()
	fmt.Println(string(resp.GetBody()))

	testEnd()
}

//POST几个字段信息并上传文件
func TestHttpClient_PostUploadFiles(t *testing.T) {
	testStart()

	url := "http://10.96.114.84/add1.php"
	client := NewHttpClient(url, context.Background())
	client.SetField("singleValue", "value")
	client.AddField("multiValue", "v1")
	client.AddField("multiValue", "v2")
	client.AddField("multiValue", "v3")
	client.AddCookie("c1", "v1")
	client.AddFile("abc", "./http_client_test.go", "my.cnf")
	client.SetUserAgent(testUserAgent)
	client.SetReferer("http://microsoft.com")
	client.SetHost("test.wendao.com")
	client.AddHeader("myHeaderKey", "myHeaderValue")
	resp, _ := client.Post()
	fmt.Println(string(resp.GetBody()))
	resp, _ = client.Post()
	fmt.Println(string(resp.GetBody()))
	client.SetUrl("http://10.96.114.84/add1.php")
	client.SetHost("phpmyadmin.wendao.com")
	resp, _ = client.Post()
	fmt.Println(string(resp.GetBody()))

	testEnd()
}

//纯上传文件
func TestHttpClient_UploadFiles(t *testing.T) {
	testStart()

	url := "http://10.96.114.84/add1.php"
	client := NewHttpClient(url, context.Background())
	client.AddCookie("c1", "v1")
	client.AddFile("abc", "./http_client_test.go", "my.cnf")
	client.SetUserAgent(testUserAgent)
	client.SetReferer("http://microsoft.com")
	client.AddHeader("myHeaderKey", "myHeaderValue")
	client.SetHost("test.wendao.com")
	resp, _ := client.Post()
	fmt.Println(string(resp.GetBody()))
	resp, _ = client.Post()
	fmt.Println(string(resp.GetBody()))
	client.SetUrl("http://10.96.114.84/add1.php")
	client.SetHost("phpmyadmin.wendao.com")
	resp, _ = client.Post()
	fmt.Println(string(resp.GetBody()))

	testEnd()
}

//直接设置请求的POST的body信息，没有字段，没有文件
func TestHttpClient_SetRawPostBody(t *testing.T) {
	testStart()

	url := "http://10.96.114.84/add1.php"
	client := NewHttpClient(url, context.Background())
	client.AddCookie("c1", "v1")
	client.SetUserAgent(testUserAgent)
	client.SetReferer("http://microsoft.com")
	client.AddHeader("myHeaderKey", "myHeaderValue")
	client.AddHeader("Content-Type", "application/json")
	client.SetHost("test.wendao.com")
	//设置原始的请求信息
	client.GetBuffer().Write([]byte("laskdfjalksd;fjals;djal;dsfjalds;kjadlsf;jk"))
	resp, err := client.Post()
	fmt.Println(err)
	fmt.Println(string(resp.GetBody()))
	resp, err = client.Post()
	fmt.Println(err)
	fmt.Println(string(resp.GetBody()))
	resp, err = client.Post()
	fmt.Println(err)
	fmt.Println(string(resp.GetBody()))

	testEnd()
}

func TestHttpClient_CheckRedirect(t *testing.T) {
	testStart()
	url := "http://10.96.114.84/redirect.php"
	client := NewHttpClient(url, context.Background())
	client.SetCheckRedirectFunc(checkRedirectFunc)
	resp, err := client.Get()
	fmt.Println(err)
	fmt.Println(string(resp.GetBody()))
	testEnd()
}

//校验跳转的函数
func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	if len(via) >= 2 {
		//fmt.Printf("%# v\n", pretty.Formatter(req))
		fmt.Printf("%# v\n", pretty.Formatter(via))
		return fmt.Errorf("retry too times")
	}
	return nil
}
