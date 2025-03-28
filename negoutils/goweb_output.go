// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        goweb_output.go
// @date        2018-02-02 17:51

package negoutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

// 输出结构体定义
type RuntofuOuput struct {
	Context *RuntofuContext //本次请求相关的上下文
	Status  int             //手动设置的响应的http code码
	Body    []byte          //暂存要发送的响应body的信息
	Cookies []string        //暂存要发送响应的cookie信息
	Started bool            //是否开始发送请求了
}

// 获取一个输出实例
func NewOutput() *RuntofuOuput {
	output := &RuntofuOuput{}
	return output
}

// 重置输出信息
func (output *RuntofuOuput) Reset(RuntofuCtx *RuntofuContext) {
	output.Context = RuntofuCtx
	output.Status = 0
	output.Body = []byte{}
	output.Cookies = []string{}
	output.Started = false

	//默认添加的几个公共头信息
	output.AddHeader("Referrer-Policy", "origin-when-cross-origin") //默认的referrer策略
	output.AddHeader("X-Powered-By", "RuntofuGo")
	output.AddHeader("Server", "RuntofuGo")
	output.AddHeader("Github", "https://github.com/liuyongshuai/Runtofugo")
	output.AddHeader("Cache-Control", "no-cache, must-revalidate")
	output.AddHeader("Expires", "0")
}

// 设置要输出的header信息
func (output *RuntofuOuput) AddHeader(key, val string) {
	output.Context.ResponseWriter.Header().Set(key, val)
}

// 设置输出的body
func (output *RuntofuOuput) SetBody(body []byte) {
	output.Body = body
}

/*
*
设置cookie，name/value是必需的，其他参数顺序如下：

	maxAgeTime：相对过期时间，单位秒，删除cookie时传负值
	path：指定的路径信息
	domain：指定的域名，默认为创建cookie的网页所属域名
	secure：只对HTTPS请求可见，对HTTP请求不可见
	httponly：对浏览器端的javascript中的document对象不可见
*/
func (output *RuntofuOuput) AddCookie(name string, value string, others ...interface{}) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", sanitizeName(name), sanitizeValue(value))
	if len(others) > 0 {
		var maxAge int64
		switch v := others[0].(type) {
		case int:
			maxAge = int64(v)
		case int32:
			maxAge = int64(v)
		case int64:
			maxAge = v
		}

		switch {
		case maxAge > 0:
			fmt.Fprintf(&b, "; Expires=%s; Max-Age=%d",
				time.Now().Add(time.Duration(maxAge)*time.Second).UTC().Format(time.RFC1123),
				maxAge)
		case maxAge < 0:
			fmt.Fprintf(&b, "; Max-Age=0")
		}
	}

	//path信息
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Path=%s", sanitizeValue(v))
		}
	} else {
		fmt.Fprintf(&b, "; Path=%s", "/")
	}

	//domain信息
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Domain=%s", sanitizeValue(v))
		}
	}

	//Secure
	if len(others) > 3 {
		var secure bool
		switch v := others[3].(type) {
		case bool:
			secure = v
		default:
			if others[3] != nil {
				secure = true
			}
		}
		if secure {
			fmt.Fprintf(&b, "; Secure")
		}
	}

	//httponly
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			fmt.Fprintf(&b, "; HttpOnly")
		}
	}
	output.Cookies = append(output.Cookies, b.String())
}

// 格式化cookie的键值名称
var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

var cookieValueSanitizer = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func sanitizeValue(v string) string {
	return cookieValueSanitizer.Replace(v)
}

// 返回json数据
func (output *RuntofuOuput) RenderJson(data interface{}) error {
	output.AddHeader("Content-Type", "application/json; charset=utf-8")
	var err error
	content, err := json.Marshal(data)
	if err != nil {
		output.Started = true
		http.Error(output.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return err
	}
	output.SetBody(content)
	return nil
}

// 响应jsonp数据，要求传个callback参数
func (output *RuntofuOuput) RenderJsonp(data interface{}, callback ...string) error {
	var content []byte
	var err error
	content, err = json.Marshal(data)
	if err != nil {
		output.Started = true
		http.Error(output.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return err
	}
	var cb string
	if len(callback) > 0 {
		cb = callback[0]
	} else {
		cb = output.Context.Input.Query("callback")
	}
	if cb == "" {
		output.Started = true
		return errors.New(`"callback" parameter required`)
	}
	output.AddHeader("Content-Type", "application/javascript; charset=utf-8")
	cb = template.JSEscapeString(cb)
	callbackContent := bytes.NewBufferString(cb)
	callbackContent.WriteString("(")
	callbackContent.Write(content)
	callbackContent.WriteString(");\r\n")
	output.SetBody(callbackContent.Bytes())
	return nil
}

// 设置响应的状态值
func (output *RuntofuOuput) SetStatus(status int) {
	output.Status = status
}

// 输出所有的内容，包括header、body、cookie等
func (output *RuntofuOuput) Send() {
	//输出cookie信息
	if len(output.Cookies) > 0 {
		for _, c := range output.Cookies {
			output.Context.ResponseWriter.Header().Add("Set-Cookie", c)
		}
	}
	if output.Started == true {
		return
	}

	if output.Status > 0 {
		output.Context.ResponseWriter.WriteHeader(output.Status)
		output.Status = 0
	}

	//输出body
	if len(output.Body) > 0 {
		buf := new(bytes.Buffer)
		buf.Write(output.Body)
		io.Copy(output.Context.ResponseWriter, buf)
	}

	output.Started = true
}
