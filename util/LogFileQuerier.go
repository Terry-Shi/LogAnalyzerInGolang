package main

import (
	"os"
	//"bufio"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"bufio"
	"compress/gzip"
)

// ref:Golang-文件操作 http://www.nljb.net/default/Golang-%E6%96%87%E4%BB%B6%E6%93%8D%E4%BD%9C/

// 1. 函数参数为可变数量
// 2. 文件操作。 a.遍历某目录下所有文件 b. 读取文件内容
// 3. 利用正则表达式过滤文件
// 4. read gzip file directly
func main() {
	queryWithDateRange("C:\\WORK\\IDE\\IdeaProjects\\LogAnalyzerInGolang", "")
}

func queryWithDateRange(path string, logdate string) map[string]string {
	ret := make(map[string]string)

	//logpath,_ := os.Open(path)
	//getFilelist(path)
	fmt.Println(" ------ ")
	files, err := ListDir(path, "md")
	if (err != nil) {
		return ret
	}
	for _, file := range files {
	    fmt.Println(file)
		//打开文件，并进行相关处理
		f , err := os.Open(file)
		if err != nil{
			fmt.Printf("%v\n",err)
			os.Exit(1)
		}
		//文件关闭 TODO: how to avoid leak ? put the code which inside the loop into a function
		defer f.Close()
		//将文件作为一个io.Reader对象进行buffered I/O操作
		var br *bufio.Reader
		if strings.HasSuffix(file, ".gz") {
			br,_ = gzip.NewReader(f)
		} else {
		        br = bufio.NewReader(f)
		}
		for{

			line , err := br.ReadString('\n')
			if err == os.EOF {
				break
			}else{
				fmt.Printf("%v",line)
			}
		}

	}
	return ret
}

// 获取指定目录下的所有文件，不进入下级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		// debug
		//fmt.Println(fi.Name())
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

// 递归遍历给定目录下所有文件 （会进入下级子目录）
func getFilelist(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	//panic("not reached")
}
