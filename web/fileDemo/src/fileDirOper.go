package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

func main() {

}

//获取系统用于保存临时文件的默认目录
func TempDirTest() {
	fmt.Println(os.TempDir())
}

//获取当前工作目录的根路径
func GetwdTest() {
	fmt.Println(os.Getwd())
}

//判断文件或目录是否存在
func IsPathExist(fileName string) (bool, error) {
	//or os.Open(fileName)
	_, err := os.Stat(fileName)
	//存在
	if err == nil {
		fmt.Println("file exist")
		return true, nil
	}
	//不存在
	if os.IsNotExist(err) {
		fmt.Println("file isn't exist")
		return false, nil
	}
	//不确定
	return false,err
}

//创建单级目录
func MkdirTest(path string,perm os.FileMode){
	err:=os.Mkdir(path,perm)
	if err!=nil {
		fmt.Println("err from Mkdir:",err)
	}
	fmt.Println("Mkdir success")
}

//创建多级目录
func MkdirAllTest(path string,perm os.FileMode){
	err:=os.MkdirAll(path,perm)
	if err!=nil{
		fmt.Println("err from MkdirAll:",err)
	}
	fmt.Println("MkdirAll success")
}

//删除文件或目录
func RemoveTest(path string){
	err:=os.Remove(path)
	if err!=nil{
		fmt.Println("err from Remove:",err)
	}
	fmt.Println("Remove success")
}

//删除文件或目录
func RemoveAll(path string){
	err:=os.RemoveAll(path)
	if err!=nil{
		fmt.Println("err from RemoveAll:",err)
	}
	fmt.Println("RemoveAll success")
}

//创建文件
func CreateTest(name string){
	f,err:=os.Create(name)
	if err!=nil{
		fmt.Println("err from Create",err)
	}
	fmt.Println("Create success",f.Name())
}

//打开文件
func OpenTest(name string){
	f,err:=os.Open(name)
	if err!=nil{
		fmt.Println("err from Open",err)
	}
	fmt.Println("Open success",f.Name())
}

//打开文件，文件不存在时创建（由os.O_CREATE指定），并以指定模式打开（这里是读写模式0666，不可执行）
func OpenFileTest(name string){
	f,err:=os.OpenFile(name,os.O_RDWR|os.O_CREATE,0666)
	if err!=nil{
		fmt.Println("err from OpenFile",err)
	}
	fmt.Println("OpenFile success",f.Name())

}

//用指定前缀建立临时文件夹，在dir目录下建立，dir为空时使用os.TempDir返回的默认目录
func TempDirPreTest(dir,prefix string){
	name,err:=ioutil.TempDir(dir,prefix)
	if err!=nil{
		fmt.Println("err from ioutil.TempDir",err)
	}
	fmt.Println("dir's name:",name)
}

//用指定前缀建立临时文件，在dir目录下建立，dir为空时使用os.TempDir返回的默认目录
func TempFileTest(dir,prefix string){
	n,err:=ioutil.TempFile(dir,prefix)
	if err!=nil{
		fmt.Println("err from ioutil.TempFile",err)
	}
	fmt.Println("file's name:",n)

}