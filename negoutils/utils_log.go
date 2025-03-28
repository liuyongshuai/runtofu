// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-29 11:34

package negoutils

import (
	"fmt"
)

// Printf 输出日志，无视当前设置的日志级别。
func LogPrintf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Debugf 输出 DEBUG 级别的日志信息。
func LogDebugf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Infof 输出 INFO 级别的日志信息。
func LogInfof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Warnf 输出 WARN 级别的日志信息。
func LogWarnf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Errorf 输出 ERROR 级别的日志信息。
func LogErrorf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Fatalf 输出 FATAL 级别的日志信息并采用 os.Exit 退出程序。
func LogFatalf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Panicf 输出 PANIC 级别的日志信息并采用 panic 抛出异常。
func LogPanicf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Print 输出日志，无视当前设置的日志级别。
func LogPrint(args ...interface{}) {
	fmt.Println(args...)
}

// Debug 输出 DEBUG 级别的日志信息。
func LogDebug(args ...interface{}) {
	fmt.Println(args...)
}

// Info 输出 INFO 级别的日志信息。
func LogInfo(args ...interface{}) {
	fmt.Println(args...)
}

// Warn 输出 WARN 级别的日志信息。
func LogWarn(args ...interface{}) {
	fmt.Println(args...)
}

// Error 输出 ERROR 级别的日志信息。
func LogError(args ...interface{}) {
	fmt.Println(args...)
}

// Fatal 输出 FATAL 级别的日志信息并采用 os.Exit 退出程序。
func LogFatal(args ...interface{}) {
	fmt.Println(args...)
}

// Panic 输出 PANIC 级别的日志信息并采用 panic 抛出异常。
func LogPanic(args ...interface{}) {
	fmt.Println(args...)
}
