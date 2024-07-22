package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	filename := "C:/Program Files (x86)/Path of Exile/logs/Client.txt" // 替换为你的文件名
	var offset int64 = 0

	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
	}
	offset = info.Size()

	for {
		readNewContent(file, &offset)
		time.Sleep(1 * time.Second) // 每隔 2 秒轮询一次
	}
}

// 读取新内容并更新偏移量
func readNewContent(file *os.File, offset *int64) {
	// 获取当前文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	// 如果文件大小小于当前偏移量，说明文件被截断了
	if fileInfo.Size() < *offset {
		*offset = 0
	}

	// 设置文件偏移量
	file.Seek(*offset, 0)

	// 读取新内容
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(line) // 处理新读取的内容
		*offset += int64(len(line))
	}
}
