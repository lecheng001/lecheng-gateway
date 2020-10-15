package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func IO_GetRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func IO_FilePathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IO_CreateFolder(path string) {
	exist := IO_FilePathExists(path)
	if exist {
		return
	} else {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			LogError(err.Error())
		}
	}
}

func IO_ReadAllLine(infile string) []string {
	file, err := os.Open(infile)
	if err != nil {
		fmt.Println("file is no exist:" + infile)
		return nil
	}
	defer file.Close()

	br := bufio.NewReader(file)
	values := make([]string, 0)
	for {
		line, isPrefix, err := br.ReadLine()
		if err != nil {
			break
		}
		if isPrefix {
			fmt.Println("the line is to long")
			return nil
		}

		values = append(values, string(line))
	}

	fmt.Println(values)
	return values
}

func IO_GetFileContnet(infile string) string{
	content,err:= ioutil.ReadFile(infile)
	if err!=nil{
		return ""
	}
	return string(content)
}

func IO_PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//取得文件路径（根据不同系统取得）
func IO_GetFilePathFileName(file string) string {
	switch runtime.GOOS {
	case "darwin":
		panic("darwin不支持")
	case "windows":
		return IO_GetRootPath() + "\\" + file
	case "linux":
		return file
	default:
		return file
	}
	return ""
}

//写入文件
func IO_WriteStringToFile(filepath, content string) {
	//打开文件，没有则创建，有则append内容
	w1, error := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if error!=nil{
		return
	}

	_, err1 := w1.Write([]byte(content))
	if err1!=nil{
		return
	}

	errC := w1.Close()
	if errC!=nil{
		return
	}
}

//写入文件
func IO_WriteBytesToFile(filepath string, content []byte) {
	//打开文件，没有此文件则创建文件，将写入的内容append进去
	w1, error := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if error!=nil{
		return
	}

	_, err1 := w1.Write(content)
	if err1!=nil{
		return
	}

	errC := w1.Close()
	if errC!=nil{
		return
	}
}