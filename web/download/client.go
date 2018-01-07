package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"io"
)

var (
	url = "http://localhost:8000/d2"
)

func main() {
	//http get请求，go提供模拟表单请求的操作，不需要编写html也能实现
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("download fail")
	}
	//记得关闭response body，否则占用http连接资源
	defer res.Body.Close()
	//将下载的文件存储到系统的临时目录中，文件目录是os.TempDir()
	dir := os.TempDir()
	tmp, err := ioutil.TempFile(dir, "zxy-") //文件名前缀设置为"zxy"
	if err != nil {
		fmt.Println("create temp file fail")
	}
	io.Copy(tmp, res.Body)
}
