package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sort"
)

/*
 ref:Golang-文件操作 http://www.nljb.net/default/Golang-%E6%96%87%E4%BB%B6%E6%93%8D%E4%BD%9C/

OK 1. 函数参数为可变数量
OK 2. 文件操作。 a.遍历某目录下所有文件 b. 读取文件内容
 3. 利用正则表达式过滤文件
OK 4. read gzip file directly
 5 错误、异常处理
 6 对map按值排序 http://golanghome.com/post/423
   https://gist.github.com/ikbear/4038654
   http://blog.csdn.net/slvher/article/details/44779081
 7 best practice: clean code
*/

// 对map按值排序 (升序)
type SortByValueMap struct {
    m map[string]float32
    s []string
}
func (sm *SortByValueMap) Len() int {
    return len(sm.m)
}
func (sm *SortByValueMap) Less(i, j int) bool {
    return sm.m[sm.s[i]] < sm.m[sm.s[j]]
}
func (sm *SortByValueMap) Swap(i, j int) {
    sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}
func sortedKeys(m map[string]float32) []string {
    sm := new(SortByValueMap)
    sm.m = m
    sm.s = make([]string, len(m))
    i := 0
    for key, _ := range m {
        sm.s[i] = key
        i++
    }
    sort.Sort(sm)
    return sm.s
}


func main() {
	ret, err := queryWithDateRange("C:\\WORK\\IDE\\IdeaProjects\\LogAnalyzerInGolang\\logfile", "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ret)

		//// 根据key排序
		//sorted_keys := make([]string, 0)
		//for k, _ := range ret {
		//	sorted_keys = append(sorted_keys, k)
		//}
		//// sort 'string' key in increasing order
		//sort.Strings(sorted_keys)
		//for _, k := range sorted_keys {
		//	fmt.Printf("k=%v, v=%v\n", k, ret[k])
		//}

        // 计算访问量占比
        var totalReqest int = 0
        for _, value := range ret {
            totalReqest = totalReqest + value
        }
        frequencyMap := make(map[string]float32)
        for key, value := range ret {
            frequencyMap[key] = float32(value)/ float32(totalReqest) * 100
        }

		// TODO: 自定义数据结构排序, 以map的value从小到大排序
		// http://www.kancloud.cn/itfanr/go-by-example/81648
		//实现了sort接口的Len，Less和Swap方法这样我们就可以使用sort包的通用方法Sort
        keysAfterSort := sortedKeys(frequencyMap)
        for _, key := range keysAfterSort {
            // 格式化 http://gobyexample.everyx.in/string-formatting/
            fmt.Printf("%s %6.4f \n", key, frequencyMap[key])
        }
	}
}

/**
 path: log file path
 logdate: which date's log you want to query
 */
func queryWithDateRange(path string, logdate ...string) (map[string]int, error) {
	ret := make(map[string]int) // key: eprid, value: request times

	files, err := ListDir(path, "gz")
	if err != nil {
		return ret, err
	}
	for _, fileFullPathName := range files {
		fmt.Println(fileFullPathName)

		func() {
			//打开文件，并进行相关处理
			file, err := os.Open(fileFullPathName)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			//文件关闭
			// DONE: how to avoid leak ? put the code which inside the loop into a function。用匿名函数解决
			defer file.Close()

			var line string = ""
			if strings.HasSuffix(fileFullPathName, ".gz") {
				//将文件作为一个io.Reader对象进行buffered I/O操作
				reader, _ := gzip.NewReader(file)
				br := bufio.NewReader(reader)
				for {
					line, err = br.ReadString('\n')
					if err != nil {
						break
					} else {
						eprid := GetEprid(line)
                        if eprid != "" {
                            ret[eprid] = ret[eprid] + 1
                        }
					}
				}
			} else {
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line = scanner.Text()
					//fmt.Printf("%v",line)
					eprid := GetEprid(line)
                    if eprid != "" {
                        ret[eprid] = ret[eprid] + 1
                    }
				}
			}
		}()
	}
	return ret, err
}

func GetEprid(oneline string) (eprid string) {
	idx := strings.Index(oneline, "/fez/")
	if idx != -1 {
		return oneline[idx+5 : idx+5+6]
	} else {
		return ""
	}
}

// 获取指定目录下的所有文件，不进入下级目录搜索，可以通过匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10) // TODO: 含义:allocates a slice of length 0 and capacity 10.
	dir, err := ioutil.ReadDir(dirPth) // TODO: 熟悉常用lib的功能 ioutil, os, fmt,  etc
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
