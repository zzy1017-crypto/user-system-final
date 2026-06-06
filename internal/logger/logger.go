package logger

import (
	"log"
)

// 定义Info函数，接受一个字符串参数msg，使用log包的Println函数输出一条日志消息，日志级别为INFO，格式为"[INFO] <msg>"，用于记录程序中的信息性事件，例如操作成功、状态变化等
func Info(msg string) {
	log.Println("[INFO]", msg)
}

// 定义Error函数，接受一个字符串参数msg，使用log包的Println函数输出一条日志消息，日志级别为ERROR，格式为"[ERROR] <msg>"，用于记录程序中的错误事件，例如操作失败、异常情况等
func Error(msg string) {
	log.Println("[ERROR]", msg)
}
