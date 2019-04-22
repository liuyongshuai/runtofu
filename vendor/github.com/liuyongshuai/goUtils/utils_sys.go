// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-12-13 17:40

package goUtils

import (
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

//内存使用情况，字节
func MemoryGetUsage() uint64 {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	return stat.Alloc
}

//returnVar：0成功、1失败
func ExecCmd(command string, output *[]string, returnVar *int) string {
	r, _ := regexp.Compile(`[ ]+`)
	parts := r.Split(command, -1)
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}
	cmd := exec.Command(parts[0], args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		*returnVar = 1
		return ""
	} else {
		*returnVar = 0
	}
	*output = strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	if l := len(*output); l > 0 {
		return (*output)[l-1]
	}
	return ""
}
