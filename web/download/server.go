package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"io"
)

func main() {
	http.HandleFunc("/d1", Download1)
	http.HandleFunc("/d2", Download2)
	http.HandleFunc("/upload", Upload)
	http.ListenAndServe(":8000", nil)
}

//将本地文件返回客户端，方式1
func Download1(w http.ResponseWriter, r *http.Request) {
	var filePath = "/Users/zxy/Desktop/trans.txt"
	//返回[]byte类型的文件数据
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("readFile error")
	}
	w.Header().Set("test", "readFile：type is []byte")

	//写response body:方法1
	//参数为[]byte类型，返回文件大小(类型:int,单位：字节）
	w.Write(f)
}

//将本地文件返回客户端，方式2
func Download2(w http.ResponseWriter, r *http.Request) {
	var filePath = "/Users/zxy/Desktop/3.png"
	//返回*File(os包)类型，File是个struct，实现了io.Reader接口
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		fmt.Println("readFile error")
	}
	w.Header().Set("test", "open-reader")

	//写response body:方法2
	//补充：io.Write()接口只有Write()方法，只要实现了该方法的结构体struct
	// 	   或拥有该方法的接口（如http.ResponseWriter），都可以作为io.Writer接口类型的参数。
	//	   io.Reader接口类型的参数同理，因为http.Request类型的结构体实现了Read()方法，故也可作为参数。
	//todo 本教程多次提到了io.Copy()方法，只要执行该方法，2个参数对应的文件指针(前提参数是个文件指针)都会移到文件末尾，
	//todo 下次再直接操作文件时，很有可能会得到空文件.需要调用"文件指针.Seek(0，0)"将指针移回文件头
	io.Copy(w, f)

}

//将上传的文件返回客户端
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseMultipartForm(10000)
		if err != nil {
			fmt.Println(err)
			return
		}
		f, header, err := r.FormFile("uploadFile")
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("uploadHeader", header.Filename)
		io.Copy(w, f)
	}
}
