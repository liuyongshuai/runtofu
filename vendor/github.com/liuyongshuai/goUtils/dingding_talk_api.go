// 文档说明：
// 		https://open-doc.dingtalk.com/docs/doc.htm?spm=a219a.7629140.0.0.karFPe&treeId=257&articleId=105735&docType=1
//
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-08 13:37

package goUtils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	DING_TALK_API_HOST = "https://oapi.dingtalk.com/robot/send?access_token=%s"
)

//构造
func NewDingTalkApi(accessToken string, t string) *DingTalkApi {
	ret := &DingTalkApi{
		AccessToken: accessToken,
		Msg: DingTalkMsgInfo{
			MsgType: t,
			At: DingTalkAt{
				AtMobiles: []string{},
				IsAtAll:   false,
			},
		},
	}
	return ret
}

//api
type DingTalkApi struct {
	AccessToken string          //access_token
	Msg         DingTalkMsgInfo //消息
}

//at所有人
func (dt *DingTalkApi) IsAtAll(b bool) *DingTalkApi {
	dt.Msg.At.IsAtAll = b
	return dt
}

//at指定人
func (dt *DingTalkApi) AtMobiles(mobiles []string) *DingTalkApi {
	dt.Msg.At.AtMobiles = mobiles
	return dt
}

//text类型的content
func (dt *DingTalkApi) SetMsgTypeText(content string) *DingTalkApi {
	if dt.Msg.MsgType != DING_TALK_MSG_TYPE_TEXT {
		fmt.Fprintf(os.Stderr, "消息类型 %s 只能调用方法\tSetMsgType%v\n", dt.Msg.MsgType, dt.Msg.MsgType)
		return dt
	}
	dt.Msg.Text = &DingTalkMsgTypeText{
		Content: content,
	}
	return dt
}

//link类型
func (dt *DingTalkApi) SetMsgTypeLink(title, text, msgUrl, picUrl string) *DingTalkApi {
	if dt.Msg.MsgType != DING_TALK_MSG_TYPE_LINK {
		fmt.Fprintf(os.Stderr, "消息类型 %s 只能调用方法\tSetMsgType%v\n", dt.Msg.MsgType, dt.Msg.MsgType)
		return dt
	}
	dt.Msg.Link = &DingTalkMsgTypeLink{
		Text:   text,
		Title:  title,
		MsgUrl: msgUrl,
		PicUrl: picUrl,
	}
	return dt
}

//markdown类型
func (dt *DingTalkApi) SetMsgTypeMarkdown(title, text string) *DingTalkApi {
	if dt.Msg.MsgType != DING_TALK_MSG_TYPE_MARKDOWN {
		fmt.Fprintf(os.Stderr, "消息类型 %s 只能调用方法\tSetMsgType%v\n", dt.Msg.MsgType, dt.Msg.MsgType)
		return dt
	}
	dt.Msg.MarkDown = &DingTalkMsgTypeMarkdown{
		Text:  text,
		Title: title,
	}
	return dt
}

//发送消息
func (dt *DingTalkApi) Send() error {
	url := fmt.Sprintf(DING_TALK_API_HOST, dt.AccessToken)
	ct := "application/json;charset=utf-8"
	httpClient := NewHttpClient(url, context.Background())
	httpClient.SetContentType(ct)
	httpClient.SetTimeout(time.Duration(10 * time.Second))
	c, err := json.Marshal(dt.Msg)
	if err != nil {
		return err
	}
	httpClient.SetRawRequestBody(c)
	resp, err := httpClient.Post()
	if err != nil {
		return err
	}
	var ret DingTalkApiResp
	err = json.Unmarshal(resp.GetBody(), &ret)
	if err != nil {
		return err
	}
	if ret.ErrCode != 0 {
		errMsg := fmt.Sprintf("data=%s||content-Type=%s||url=%s||errCode=%v||errMsg=%v", string(c), ct, url, ret.ErrCode, ret.ErrMsg)
		return fmt.Errorf(errMsg)
	}
	return nil
}

//纯文本类型消息结构体
type DingTalkMsgInfo struct {
	MsgType  string                   `json:"msgtype"`
	Text     *DingTalkMsgTypeText     `json:"text,omitempty"`
	Link     *DingTalkMsgTypeLink     `json:"link,omitempty"`
	MarkDown *DingTalkMsgTypeMarkdown `json:"markdown,omitempty"`
	At       DingTalkAt               `json:"at,omitempty"`
}

const (
	DING_TALK_MSG_TYPE_TEXT     = "text"
	DING_TALK_MSG_TYPE_MARKDOWN = "markdown"
	DING_TALK_MSG_TYPE_LINK     = "link"
)

//纯文本类型
type DingTalkMsgTypeText struct {
	Content string `json:"content"`
}

//MarkDown类型
type DingTalkMsgTypeMarkdown struct {
	Title string `json:"title"` //[必选]首屏会话透出的展示内容
	Text  string `json:"text"`  //[必选]markdown格式的消息
}

//link类型
type DingTalkMsgTypeLink struct {
	Title  string `json:"title"`            //[必选]消息标题
	Text   string `json:"text"`             //[必选]消息内容。如果太长只会部分展示
	MsgUrl string `json:"messageUrl"`       //[必选]点击消息跳转的URL
	PicUrl string `json:"picUrl,omitempty"` //[可选]图片URL
}

//AT信息
type DingTalkAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"` //被@人的手机号
	IsAtAll   bool     `json:"isAtAll"`             //@所有人时:true,否则为:false
}

//响应值
type DingTalkApiResp struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
