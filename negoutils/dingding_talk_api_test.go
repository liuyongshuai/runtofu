// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-26 14:38

package negoutils

import (
	"testing"
)

func TestDingTalkApi_Send(t *testing.T) {
	testStart()

	dt := NewDingTalkApi("abc", DING_TALK_MSG_TYPE_TEXT)
	dt.SetMsgTypeText("测试搜索数据重建监控报警的钉钉的机器人功能！请忽略")
	dt.IsAtAll(false)
	err := dt.Send()
	if err == nil {
		t.Fail()
	}

	testEnd()
}
