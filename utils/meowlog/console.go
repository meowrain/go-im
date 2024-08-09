package meowlog

import (
	"fmt"
	"os"
	"time"
)

type MeowConsoleLog struct {
	LogLevel
}

// NewMeowConsoleLog 创建一个MeowConsoleLog实例
func newMeowConsoleLog(level string) *MeowConsoleLog {
	loglevel, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	return &MeowConsoleLog{loglevel}
}

// SetLogLevel 设置日志级别
func (m *MeowConsoleLog) SetLogLevel(level string) {
	loglevel, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	m.LogLevel = loglevel
}

func (m *MeowConsoleLog) GetLogLevel(level LogLevel) string {
	switch level {
	case TraceLevel:
		return "Trace"
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case WarnLevel:
		return "Warn"
	case ErrorLevel:
		return "Error"
	case FatalLevel:
		return "Fatal"
	default:
		return "Unknown"
	}
}

// enable enable 函数的目的是为了检查当前的日志级别是否允许记录某个特定级别的日志消息。如果当前的日志级别设置为 WarnLevel，那么只有 WarnLevel 及以上的日志级别会被记录。这意味着 DebugLevel 和 TraceLevel 的日志消息不会被记录，因为它们的级别低于 WarnLevel。
func (m *MeowConsoleLog) enable(level LogLevel) bool {
	return level <= m.LogLevel
}

// log 函数的目的是根据日志级别和格式化字符串，生成日志消息，并打印到控制台。
func (m *MeowConsoleLog) log(level LogLevel, format string, args ...any) {
	if m.enable(level) {
		msg := fmt.Sprintf(format, args...)
		now := time.Now()
		funcName, fileName, lineNum := getInfo(3)
		fmt.Printf("%s %s [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05"), levelColor(level), funcName, fileName, lineNum, msg)
	}
}

// Debug 打印 Debug 级别的日志消息
func (m *MeowConsoleLog) Debug(format string, args ...any) {
	m.log(DebugLevel, format, args...)
}

// Trace 打印 Trace 级别的日志消息
func (m *MeowConsoleLog) Trace(format string, args ...any) {
	m.log(TraceLevel, format, args...)
}

// Info 打印 Info 级别的日志消息
func (m *MeowConsoleLog) Info(format string, args ...any) {
	m.log(InfoLevel, format, args...)
}

// Warn 打印 Warn 级别的日志消息
func (m *MeowConsoleLog) Warn(format string, args ...any) {
	m.log(WarnLevel, format, args...)
}

// Error 打印 Error 级别的日志消息
func (m *MeowConsoleLog) Error(format string, args ...any) {
	m.log(ErrorLevel, format, args...)
}

// Fatal 打印 Fatal 级别的日志消息，并调用 os.Exit(1) 退出程序
func (m *MeowConsoleLog) Fatal(format string, args ...any) {
	m.log(FatalLevel, format, args...)
	os.Exit(1)
}
