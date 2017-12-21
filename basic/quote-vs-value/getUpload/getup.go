package getUpload

import (
	//"mime/multipart"
	"os"
	"fmt"
	"io"
	"net/http"
)

//（r *http.Request）传r的指针方式可以
func Upload(r *http.Request,reader io.Reader,fileName string){
	err:=os.Mkdir("./temp",os.ModePerm)
	if err!=nil{
		fmt.Println(err)
	}
	//创建并打开文件（可读写的模式），如果文件存在则返回原文件指针
	f,err:=os.Create("./temp/"+fileName)
	if err!=nil{
		fmt.Println(err)
	}
	defer f.Close()

	//open以只读方式打开
	//f1,err:=os.Open(f.Name())

	//f1,err:=os.OpenFile(f.Name(),os.O_RDWR,os.ModePerm)
	//defer f1.Close()
	io.Copy(f,reader)
}

