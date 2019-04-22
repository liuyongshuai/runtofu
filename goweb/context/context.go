package context

import (
	"fmt"
	"github.com/liuyongshuai/goUtils"
	"net/http"
)

var snowFlake *goUtils.SnowFlakeIdGenerator

func init() {
	snowFlake, _ = goUtils.
		NewIDGenerator().
		SetTimeBitSize(50).
		SetSequenceBitSize(8).
		SetWorkerIdBitSize(5).
		SetWorkerId(1).
		Init()
}

//返回请求上下文
func NewWeGoContext() *WeGoContext {
	WeGoCtx := &WeGoContext{
		Input:  NewInput(),
		Output: NewOutput(),
	}
	WeGoCtx.Output.Context = WeGoCtx
	WeGoCtx.Input.Context = WeGoCtx
	return WeGoCtx
}

//上下文的定义
type WeGoContext struct {
	Input          *WeGoInput          //收到的请求里相关信息，包括参数、方法、上传文件等
	Output         *WeGoOuput          //要发送给端的暂存用的数据
	Request        *http.Request       //请求原始对象指针
	ResponseWriter http.ResponseWriter //响应原始对象
	UniqueKey      string              //本次请求的唯一标识符
}

//重置本次请求的上下文
func (WeGoCtx *WeGoContext) Reset(rw *http.ResponseWriter, r *http.Request) {
	WeGoCtx.Request = r
	WeGoCtx.ResponseWriter = *rw
	WeGoCtx.Input.Reset(WeGoCtx)
	WeGoCtx.Output.Reset(WeGoCtx)
	var nextId int64 = 0
	var err error
	for {
		nextId, err = snowFlake.NextId()
		if err == nil {
			break
		}
	}
	WeGoCtx.UniqueKey = fmt.Sprintf("%x", nextId)
}

//跳转，状态码可选，默认301
func (WeGoCtx *WeGoContext) Redirect(locationUrl string, status ...int) {
	code := http.StatusTemporaryRedirect
	if len(status) > 0 {
		code = status[0]
	}
	http.Redirect(WeGoCtx.ResponseWriter, WeGoCtx.Request, locationUrl, code)
}

//刷新返回数据
func (WeGoCtx *WeGoContext) Flush() {
	if f, ok := WeGoCtx.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

//当客户端取消请求或连接断开时用
func (WeGoCtx *WeGoContext) CloseNotify() <-chan bool {
	if cn, ok := WeGoCtx.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}
