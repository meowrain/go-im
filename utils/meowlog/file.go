package meowlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type MeowFileLog struct {
	LogLevel           //日志级别
	file     *os.File  //日志文件
	lastTime time.Time //上一次写入时间
}

// newMeowFileLog 创建一个新的日志对象
func newMeowFileLog(level string, dir string) *MeowFileLog {
	loglevel, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	dirExists := checkDirExists("logs")
	if !dirExists {
		err := os.Mkdir("logs", 0777)
		if err != nil {
			panic(err)
		}
	}
	currentDate := time.Now().Format("2006-01-02")
	filePath := fmt.Sprintf("%s/%s.log", dir, currentDate)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return &MeowFileLog{loglevel, file, time.Now()}
}
func (m *MeowFileLog) enable(level LogLevel) bool {
	return level <= m.LogLevel
}

// log 记录日志消息到文件
func (m *MeowFileLog) log(level LogLevel, format string, args ...any) {
	if m.enable(level) {
		msg := fmt.Sprintf(format, args...)
		now := time.Now()
		// 切换到新的日志文件
		if m.lastTime.Day() != now.Day() {
			m.closeAndOpenFile(now)
		}
		funcName, fileName, lineNum := getInfo(3)
		logMsg := fmt.Sprintf("%s %s [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05"), levelColor(level), funcName, fileName, lineNum, msg)
		_, err := m.file.WriteString(logMsg)
		if err != nil {
			fmt.Printf("Failed to write log to file: %v\n", err)
		}
		// 更新最后一次写入时间
		m.lastTime = now
	}
}

// closeAndOpenFile 关闭当前日志文件并打开新的日志文件
func (m *MeowFileLog) closeAndOpenFile(now time.Time) {
	if m.file != nil {
		m.file.Close()
	}

	// 重新打开新一天的日志文件
	currentDate := now.Format("2006-01-02")
	filePath := filepath.Join("logs", fmt.Sprintf("%s.log", currentDate))
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
	} else {
		m.file = file
		m.lastTime = now
	}
}

// SetLogLevel 设置日志级别
func (m *MeowFileLog) SetLogLevel(level string) {
	loglevel, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	m.LogLevel = loglevel
}

// GetLogLevel 获取日志级别
func (m *MeowFileLog) GetLogLevel(level LogLevel) string {
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

// Debug 输出调试日志
func (m *MeowFileLog) Debug(format string, args ...any) {
	m.log(DebugLevel, format, args...)
}

// Trace 输出跟踪日志
func (m *MeowFileLog) Trace(format string, args ...any) {
	m.log(TraceLevel, format, args...)
}

// Info 输出信息日志
func (m *MeowFileLog) Info(format string, args ...any) {
	m.log(InfoLevel, format, args...)
}

// Warn 输出警告日志
func (m *MeowFileLog) Warn(format string, args ...any) {
	m.log(WarnLevel, format, args...)
}

// Error 输出错误日志
func (m *MeowFileLog) Error(format string, args ...any) {
	m.log(ErrorLevel, format, args...)
}

// Fatal 输出致命错误日志
func (m *MeowFileLog) Fatal(format string, args ...any) {
	m.log(FatalLevel, format, args...)
	os.Exit(1)
}
