package meowlog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

// parseLogLevel 解析日志级别
func parseLogLevel(level string) (LogLevel, error) {
	level = strings.ToLower(level)
	switch {
	case level == "trace":
		return TraceLevel, nil
	case level == "debug":
		return DebugLevel, nil
	case level == "info":
		return InfoLevel, nil
	case level == "warn":
		return WarnLevel, nil
	case level == "error":
		return ErrorLevel, nil
	case level == "fatal":
		return FatalLevel, nil
	default:
		return 0, fmt.Errorf("日志级别错误: %s", level)
	}
}

// colorize 给文本着色
func colorize(text string, color string) string {
	colorCode := getEscapeCode(color)
	return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, text)
}

// resetColor 返回重置颜色的转义序列
func resetColor() string {
	return "\033[0m"
}

// getEscapeCode 获取颜色的转义序列
func getEscapeCode(color string) string {
	switch strings.ToLower(color) {
	case "black":
		return "30"
	case "red":
		return "31"
	case "green":
		return "32"
	case "yellow":
		return "33"
	case "blue":
		return "34"
	case "magenta":
		return "35"
	case "cyan":
		return "36"
	case "white":
		return "37"
	case "gray":
		return "90"
	default:
		return "0"
	}
}
func levelColor(level LogLevel) string {
	switch level {
	case TraceLevel:
		return colorize("【TRACE】", "blue")
	case DebugLevel:
		return colorize("【DEBUG】", "green")
	case InfoLevel:
		return colorize("【INFO】", "blue")
	case WarnLevel:
		return colorize("【WARN】", "yellow")
	case ErrorLevel:
		return colorize("【ERROR】", "red")
	case FatalLevel:
		return colorize("【FATAL】", "red")
	default:
		return resetColor()

	}
}

// getInfo 获取调用者信息
func getInfo(skip int) (funcName, fileName string, lineNum int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Print("runtime.Caller() failed")
		return

	}
	funcName = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
	fileName = path.Base(file)
	lineNum = line
	return
}

// checkDirExists 检查目录是否存在并且是一个目录
func checkDirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// checkFileExists 检查文件是否存在并且是一个文件
func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
