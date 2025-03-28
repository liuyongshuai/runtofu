package negoutils

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestNewIDGenerator(t *testing.T) {
	testStart()

	b := "\t\t\t"
	b2 := "\t\t\t\t\t"
	d := "====================================="

	//第一个生成器
	gentor1, err := NewIDGenerator().SetWorkerId(100).Init()
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
	//第二个生成器
	gentor2, err := NewIDGenerator().
		SetTimeBitSize(48).
		SetSequenceBitSize(10).
		SetWorkerIdBitSize(5).
		SetWorkerId(30).Init()
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	fmt.Printf("%s%s%s\n", d, b, d)
	fmt.Printf("workerId=%d lastTimestamp=%d %s workerId=%d lastTimestamp=%d\n",
		gentor1.workerId, gentor1.lastMsTimestamp, b,
		gentor2.workerId, gentor2.lastMsTimestamp)
	fmt.Printf("sequenceBitSize=%d timeBitSize=%d %s sequenceBitSize=%d timeBitSize=%d\n",
		gentor1.sequenceBitSize, gentor1.timeBitSize, b,
		gentor2.sequenceBitSize, gentor2.timeBitSize)
	fmt.Printf("workerBitSize=%d sequenceBitSize=%d %s workerBitSize=%d sequenceBitSize=%d\n",
		gentor1.workerIdBitSize, gentor1.sequenceBitSize, b,
		gentor2.workerIdBitSize, gentor2.sequenceBitSize)
	fmt.Printf("%s%s%s\n", d, b, d)

	var ids []int64
	for i := 0; i < 100; i++ {
		id1, err := gentor1.NextId()
		if err != nil {
			fmt.Println(err)
			return
		}
		id2, err := gentor2.NextId()
		if err != nil {
			fmt.Println(err)
			return
		}
		ids = append(ids, id2)
		fmt.Printf("%d%s%d\n", id1, b2, id2)
	}

	//解析ID
	for _, id := range ids {
		ts, workerId, seq, err := gentor2.Parse(id)
		fmt.Printf("id=%d\ttimestamp=%d\tworkerId=%d\tsequence=%d\terr=%v\n",
			id, ts, workerId, seq, err)
	}
	testEnd()
}

// 多线程测试
func TestSnowFlakeIdGenerator_MultiThread(t *testing.T) {
	testStart()

	ff := "./snowflake.txt"
	//准备写入的文件
	fp, err := os.OpenFile(ff, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	//初始化ID生成器，采用默认参数
	gentor, err := NewIDGenerator().SetWorkerId(100).Init()
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	wg := new(sync.WaitGroup)

	//启动10个线程，出错就报出来
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 100; i++ {
				gid, err := gentor.NextId()
				if err != nil {
					panic(err)
				}
				n, err := fp.WriteString(fmt.Sprintf("%d\n", gid))
				if err != nil || n <= 0 {
					panic(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fp.Close()

	testEnd()
}
