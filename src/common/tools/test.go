// package tools
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const FILE_NAME = "go-example.txt"

func main() {

	// 设置信号处理
	SetupCloseHandler()

	// 运行我们的程序，也就是创建一个文件并在睡眠 10 秒后继续执行
	CreateFile()
	for {
		fmt.Println("- Sleeping")
		time.Sleep(10 * time.Second)
	}
}

// SetupCloseHandler 在一个新的 goroutine 上创建一个监听器。
// 如果接收到了一个 interrupt 信号，就会立即通知程序，做一些清理工作并退出
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		DeleteFiles()
		os.Exit(0)
	}()
}

// 模仿退出时程序的清理工作。
// 因为这只是一个示例，我们并没有对错误进行任何处理
func DeleteFiles() {
	fmt.Println("- Run Clean Up - Delete Our Example File")
	_ = os.Remove(FILE_NAME)
	fmt.Println("- Good bye!")
}

// 创建一个文件，目的是模仿退出时的清理工作
func CreateFile() {
	fmt.Println("- Create Our Example File")
	file, _ := os.Create(FILE_NAME)
	defer file.Close()
}
