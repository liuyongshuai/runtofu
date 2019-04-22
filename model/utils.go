/**
 * @author      Liu Yongshuai
 * @package     model
 * @date        2018-03-22 16:50
 */
package model

import (
	"sync"
	"time"
)

const (
	//要减去的基准时间点
	START_TIMESTAMP = 1514736000 //2018-01-01 00:00:00

	//用户ID中序号部分位数
	USER_SEQUENCE_BIT_SIZE = 8
)

var (
	//用户锁
	userSyncLock = new(sync.RWMutex)

	//产生用户ID上一个时间戳
	user_pre_tm int64 = 0

	//产生用户ID时上一个序号
	user_pre_seq int64 = 0

	//用户ID时的序号最大值
	max_user_seq int64 = -1 ^ (-1 << USER_SEQUENCE_BIT_SIZE)
)

//产生一个uid
func GenUid() int64 {
	userSyncLock.Lock()
	defer userSyncLock.Unlock()
	var curSeq int64 = 0
	for {
		curTM := time.Now().Unix()
		if curTM < user_pre_tm {
			panic("invalid timestamp")
		}
		//时间相等时得判断序号
		if curTM == user_pre_tm {
			if user_pre_seq >= max_user_seq {
				continue
			} else {
				curSeq = user_pre_seq + 1
				user_pre_seq = curSeq
			}
		}
		user_pre_tm = curTM
		curTM = curTM - START_TIMESTAMP
		curTM = ((curTM << USER_SEQUENCE_BIT_SIZE) | curSeq) - START_TIMESTAMP
		return curTM
	}
}
