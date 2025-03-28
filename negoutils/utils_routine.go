// 同一类任务多个goroutin的简单封装
// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-12-06 23:32

package negoutils

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 构造
func NewRoutineTool(routineLimit, argQueueSize int, fn RoutineFunc) *RoutineTool {
	ret := &RoutineTool{
		routineLimit: routineLimit,
		retryTimes:   0,
		workFunc:     fn,
		argQueue:     make(chan routineArgs, argQueueSize),
		wg:           new(sync.WaitGroup),
	}
	ret.wg.Add(routineLimit)
	for i := 0; i < routineLimit; i++ {
		go ret.doWork()
	}
	return ret
}

type RoutineFunc func(arg interface{}, retryTime int, commonArg interface{}) error

// 内部使用的参数
type routineArgs struct {
	argData   interface{}
	retryTime int
}

// 管理用的结构体
type RoutineTool struct {
	routineLimit   int              //总的工作线程数量
	workFunc       RoutineFunc      //工作函数
	retryTimes     int              //执行失败后的重试次数
	argQueue       chan routineArgs //传递参数的通道
	argQueueSize   int              //参数通道大小
	waitRoutineNum int32            //等待提取数据的线程数量
	commonArg      interface{}      //通用的参数
	wg             *sync.WaitGroup  //控制各线程用的
}

// 设置重试次数
func (rp *RoutineTool) SetRetryTimes(retryTime int) {
	rp.retryTimes = retryTime
}

// 设置通用参数
func (rp *RoutineTool) SetCommonArg(arg interface{}) {
	rp.commonArg = arg
}

// 将文件的每一行添加到参数里
func (rp *RoutineTool) AddFileLines(files []string) {
	for _, f := range files {
		if !FileExists(f) {
			msg := fmt.Sprintf("file[%s] not exists", f)
			fmt.Println(msg)
			continue
		}
		iter := NewFileIterator().SetFile(f)
		_, err := iter.Init()
		if err != nil {
			fmt.Println(f, err)
			continue
		}
		iter.IterLine(func(line string) {
			rp.AddArg(line)
		})
	}
	return
}

// 添加参数
func (rp *RoutineTool) AddArg(arg interface{}) {
	rp.setArg(routineArgs{argData: arg, retryTime: -1})
}

// 添加参数
func (rp *RoutineTool) setArg(arg routineArgs) {
	arg.retryTime++
	rp.argQueue <- arg
}

// 添加重试次数
func (rp *RoutineTool) Wait() {
	for {
		if int(rp.waitRoutineNum) >= rp.routineLimit {
			close(rp.argQueue)
			break
		}
		time.Sleep(time.Duration(10 * time.Millisecond))
	}

	rp.wg.Wait()
}

// 运行
func (rp *RoutineTool) doWork() {
	for {
		atomic.AddInt32(&rp.waitRoutineNum, 1)
		arg, ok := <-rp.argQueue
		atomic.AddInt32(&rp.waitRoutineNum, -1)
		if !ok {
			rp.wg.Done()
			return
		}
		if arg.retryTime > rp.retryTimes {
			continue
		}
		if rp.workFunc(arg.argData, arg.retryTime, rp.commonArg) != nil {
			rp.setArg(arg)
		}
	}
}
