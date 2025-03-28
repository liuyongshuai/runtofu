// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        redis_client_test.go
// @date        2025-03-25 10:56

package negoutils

import (
	"fmt"
	"testing"
	"time"
)

func TestRedisClient(t *testing.T) {
	r, e := InitRedisClient("100.69.239.173:3000", 60*time.Second)
	if e != nil {
		fmt.Println(e)
		return
	}
	v, e := r.Get("vector_20220620:raw_1080125649252777986")
	fmt.Println(v, e)
}
