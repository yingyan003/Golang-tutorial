package main

import (
	"net/http"
	"fmt"
	"os"
	"io"
)
const(
	maxMemory =100000
)

//var path="/"

func main(){
	http.HandleFunc("/upload",Upload)
	http.ListenAndServe("localhost:8000",nil)
}

func Upload(w http.ResponseWriter,r * http.Request){
	fmt.Println("method:",r.Method)
	fmt.Println("header:",r.Header)
	if r.Method=="POST"{
		r.ParseMultipartForm(maxMemory)
		file,fileHeader,err:=r.FormFile("uploadFile")
		if err!=nil{
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w,"fileHeader.Header:\n%v",fileHeader.Header)
		fmt.Println("fileHeader.Filename",fileHeader.Filename)
		path,err:=os.Getwd()
		fmt.Println("path:",path)
		//创建文件夹
		err=os.MkdirAll("./doc/uploadFiles",os.ModePerm)
		if err!=nil{
			fmt.Println(err)
		}
		//当找不到该文件所在的目录时，竟然会返回nil
		fcreate,err:=os.Create("./doc/uploadFiles/"+fileHeader.Filename)
		if err!=nil{
			fmt.Println(err)
		}
		fmt.Println(fcreate)
		f,err:=os.OpenFile("./doc/uploadFiles/"+fileHeader.Filename,os.O_WRONLY|os.O_CREATE,0666)
		if err!=nil{
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f,file)
	}
}