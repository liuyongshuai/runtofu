package context

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

//请求输入的相关信息结构体
type WeGoInput struct {
	Context     *WeGoContext      //上下文的指针
	Args        map[string]string //所有的参数
	RequestBody []byte
	Controller  reflect.Type //相关的控制层
}

//新建输入结构体
func NewInput() *WeGoInput {
	return &WeGoInput{
		Args:       make(map[string]string),
		Controller: nil,
		Context:    nil,
	}
}

//重置输入数据
func (input *WeGoInput) Reset(WeGoCtx *WeGoContext) {
	input.Context = WeGoCtx
	input.Args = make(map[string]string)
	input.RequestBody = []byte{}
	input.Controller = nil
}

//提取请求时用的协议，如"HTTP/1.1"
func (input *WeGoInput) Protocol() string {
	return input.Context.Request.Proto
}

//获取请求的URI信息
func (input *WeGoInput) URI() string {
	return input.Context.Request.RequestURI
}

//获取请求的URL信息
func (input *WeGoInput) URL() string {
	return input.Context.Request.URL.Path
}

//请求的站点信息，如scheme://domain
func (input *WeGoInput) Site() string {
	return input.Scheme() + "://" + input.Domain()
}

//请求协议，一般为“http”、“https”
func (input *WeGoInput) Scheme() string {
	if scheme := input.Header("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if input.Context.Request.URL.Scheme != "" {
		return input.Context.Request.URL.Scheme
	}
	if input.Context.Request.TLS == nil {
		return "http"
	}
	return "https"
}

//域名信息
func (input *WeGoInput) Domain() string {
	return input.Host()
}

//域名信息
func (input *WeGoInput) GetCookie(key string) string {
	return input.Cookie(key)
}

//域名信息
func (input *WeGoInput) Host() string {
	if input.Context.Request.Host != "" {
		hostParts := strings.Split(input.Context.Request.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return input.Context.Request.Host
	}
	return "localhost"
}

//请求方法名，GET/POST.....
func (input *WeGoInput) Method() string {
	return input.Context.Request.Method
}

//判断请求是否为某个方法
func (input *WeGoInput) Is(method string) bool {
	return input.Method() == method
}

//是否为GET
func (input *WeGoInput) IsGet() bool {
	return input.Is("GET")
}

//是否为POST
func (input *WeGoInput) IsPost() bool {
	return input.Is("POST")
}

//是否为DELETE
func (input *WeGoInput) IsDelete() bool {
	return input.Is("DELETE")
}

//是否为PUT
func (input *WeGoInput) IsPut() bool {
	return input.Is("PUT")
}

//是否为PATCH
func (input *WeGoInput) IsPatch() bool {
	return input.Is("PATCH")
}

//是否为Ajax请求
func (input *WeGoInput) IsAjax() bool {
	return input.Header("X-Requested-With") == "XMLHttpRequest"
}

//是否为https
func (input *WeGoInput) IsSecure() bool {
	return input.Scheme() == "https"
}

//是否为上传文件的请求
func (input *WeGoInput) IsUpload() bool {
	return strings.Contains(input.Header("Content-Type"), "multipart/form-data")
}

//客户端
func (input *WeGoInput) IP() string {
	ips := input.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		return rip[0]
	}
	ip := strings.Split(input.Context.Request.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

// Proxy returns proxy client ips slice.
func (input *WeGoInput) Proxy() []string {
	if ips := input.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

//返回Referer信息
func (input *WeGoInput) Referer() string {
	return input.Header("Referer")
}

//返回Referer信息
func (input *WeGoInput) Refer() string {
	return input.Referer()
}

//返回UserAgent
func (input *WeGoInput) UserAgent() string {
	return input.Header("User-Agent")
}

//参数长度
func (input *WeGoInput) ParamsLen() int {
	return len(input.Args)
}

//提取某个参数
func (input *WeGoInput) Param(key string) string {
	ret, ok := input.Args[key]
	if !ok {
		return ""
	}
	return ret
}

//所有的参数
func (input *WeGoInput) Params() map[string]string {
	return input.Args
}

//设置某个参数的值
func (input *WeGoInput) SetParam(key, val string) {
	input.Args[key] = val
}

//清除所有的参数
func (input *WeGoInput) ResetParams() {
	input.Args = make(map[string]string)
}

//提取一个参数值，包括/POST
func (input *WeGoInput) Query(key string) string {
	if val := input.Param(key); val != "" {
		return val
	}
	if input.Context.Request.Form == nil {
		input.Context.Request.ParseForm()
	}
	return input.Context.Request.Form.Get(key)
}

//提取头信息里的信息
func (input *WeGoInput) Header(key string) string {
	return input.Context.Request.Header.Get(key)
}

//提取某个cookie值
func (input *WeGoInput) Cookie(key string) string {
	ck, err := input.Context.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

//以字节切片的形式返回原始的请求body信息
func (input *WeGoInput) CopyBody(MaxMemory int64) []byte {
	if input.Context.Request.Body == nil {
		return []byte{}
	}

	var requestbody []byte
	safe := &io.LimitedReader{R: input.Context.Request.Body, N: MaxMemory}
	if input.Header("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(safe)
		if err != nil {
			return nil
		}
		requestbody, _ = ioutil.ReadAll(reader)
	} else {
		requestbody, _ = ioutil.ReadAll(safe)
	}

	input.Context.Request.Body.Close()
	bf := bytes.NewBuffer(requestbody)
	input.Context.Request.Body = http.MaxBytesReader(input.Context.ResponseWriter, ioutil.NopCloser(bf), MaxMemory)
	input.RequestBody = requestbody
	return requestbody
}

//解析请求的表单
func (input *WeGoInput) ParseFormOrMulitForm(maxMemory int64) error {
	if strings.Contains(input.Header("Content-Type"), "multipart/form-data") {
		if err := input.Context.Request.ParseMultipartForm(maxMemory); err != nil {
			return errors.New("Error parsing request body:" + err.Error())
		}
	} else if err := input.Context.Request.ParseForm(); err != nil {
		return errors.New("Error parsing request body:" + err.Error())
	}
	return nil
}
