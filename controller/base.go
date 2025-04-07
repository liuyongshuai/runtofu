// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        base.go
// @date        2025-04-07 12:00

package controller

// ajax响应总体返回的数据格式
type AjaxResponseResult struct {
	Errcode int64       `json:"error"` //返回的错误码
	ErrMsg  string      `json:"msg"`   //返回的错误信息
	Data    interface{} `json:"data"`
}
