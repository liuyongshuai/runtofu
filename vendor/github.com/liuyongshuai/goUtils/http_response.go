// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-22 18:38

package goUtils

import "net/http"

//响应结构体，在response基础上封装
type HttpResponse struct {
	body             []byte            //响应的body
	status           string            //响应码的描述信息，"200 OK"
	statusCode       int               //响应码"200"
	proto            string            //所用协议"HTTP/1.1"
	header           map[string]string //头信息
	contentLen       int64             //返回内容长度
	setCookie        []http.Cookie     //要设置的cookie信息
	location         string            //当statusCode为3XX如301时重定向的链接
	err              error             //请求的出错信息
	transferEncoding []string          //所用的编码信息
}

//提取body信息
func (tfr HttpResponse) GetBody() []byte {
	return tfr.body
}

//返回响应的body的字符串格式
func (tfr HttpResponse) GetBodyString() string {
	return string(tfr.body)
}

//提取状态描述信息，如"200 OK"
func (tfr HttpResponse) GetStatus() string {
	return tfr.status
}

//提取状态码，如200
func (tfr HttpResponse) GetStatusCode() int {
	return tfr.statusCode
}

//所用协议"HTTP/1.1"
func (tfr HttpResponse) GetProto() string {
	return tfr.proto
}

//提取响应的头信息
func (tfr HttpResponse) GetHeader() map[string]string {
	return tfr.header
}

//获取响应信息的长度
func (tfr HttpResponse) GetContentLen() int64 {
	return tfr.contentLen
}

//获取响应头里面set-cookie设置cookie信息的列表
func (tfr HttpResponse) GetSetCookie() []http.Cookie {
	return tfr.setCookie
}

//获取重定向后的地址信息
func (tfr HttpResponse) GetLocation() string {
	return tfr.location
}

//提取响应的编码信息
func (tfr HttpResponse) GetTransferEncoding() []string {
	return tfr.transferEncoding
}

//返回出错的信息
func (tfr HttpResponse) Error() error {
	return tfr.err
}
