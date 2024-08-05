package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	sourceFile := "D:\\go\\socket\\iocopy/source.txt"
	destinationFile := "D:\\go\\socket\\iocopy/destination.txt"

	// 打开源文件
	src, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println("无法打开源文件:", err)
		return
	}
	// 调用Close方法
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(destinationFile)
	if err != nil {
		fmt.Println("无法创建目标文件:", err)
		return
	}
	// 调用Close 方法
	defer dst.Close()

	// 执行文件复制
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println("复制文件出错:", err)
		return
	}

	fmt.Println("文件复制成功!")
}
