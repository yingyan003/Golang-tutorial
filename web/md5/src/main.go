package main

import (
	"net/http"
	"fmt"
	"os"
	"Golang-tutorial/web/md5/src/md5Sum"
)

//对上传不熟悉的同学，请先移步upload demo
func upload(w http.ResponseWriter,r *http.Request){
	err:=r.ParseMultipartForm(100000)
	if err!=nil{
		fmt.Println(err)
	}
	f,_,err:=r.FormFile("uploadFile")

	fmd5:=md5Sum.NewFileMd5Sum("",f)
	fmd5.Sum()
	w.Write([]byte(fmd5.String()))
	fmt.Println("stream:",fmd5.String())
}

func main(){
	//测试本地文件
	f,err:=os.Open("/Users/zxy/Desktop/3.png")
	defer f.Close()
	if err!=nil{
		fmt.Println(err)
	}
	fmd5:=md5Sum.NewFileMd5Sum(f.Name(),nil)
	fmd5.Sum()
	fmt.Println(fmd5.String())

	//测试上传的文件流
	http.HandleFunc("/upload",upload)
	http.ListenAndServe(":8000",nil)

}