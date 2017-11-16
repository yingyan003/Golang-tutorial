package main

import (
	"fmt"
	"os"
	//"io/ioutil"
	//"io"
)
//var path="/"

func main(){
	//系统的临时目录：/var/folders/v_/66t4cmys0sg90gdkll8wbvpm0000gn/T/
	fmt.Println(os.TempDir())

	//获取当前的目录:os.Getwd: /Users/zxy/techAbout/GoWorkspace/web/uploadDemo
	wd,err:=os.Getwd()
	fmt.Println("os.Getwd:",wd)

	//todo 在当前目录下创建文件夹,如果文件存在，返回err！=nil
	//.：指项目目录下
	err=os.Mkdir("./doc/uploadFiles",os.ModePerm)
	if err!=nil{
		fmt.Println("create os.MKdir:",err)
	}

	//创建文件,根据文件名（路径名：包括路径与文件名）。若路径不存在，返回nil
	fcreate,err:=os.Create("./doc/uploadFiles/2.txt")
	fmt.Println("fceate:",fcreate)

	/*
	//todo 创建目录,dir:指定文件目录所在的路径（必须存在，否则返回nil）
	dir,err:=ioutil.TempDir("./","test2")
	if err!=nil{
		fmt.Println("create dir:",err)
	}
	fmt.Println("dir",dir)
	//根据指定文件全路径名，创建文件
	//todo 只创建文件，不创建目录，当创建的文件名包含目录名，而该目录不存在时返回nil
	f,err:=os.Create(dir+"/1.txt")

	fmt.Println("f:",f)
	file1,err:=os.Open("./test")
	defer file1.Close()

	//os.O_RDONLY以只读的形式打开一个文件，perm：文件权限相关的操作
	//file2,err:=os.OpenFile("./test",os.O_RDONLY,0666)
	//defer file2.Close()

	if err!=nil{
		fmt.Println("OPEN FILE:",err)
	}
	//当文件不存在时，对文件的具体的操作将导致panic
	fmt.Println(file1.Name(),file1)
	fmt.Fprintf(os.Stdout,"%q\n",*file1)
	fmt.Println(1)
	*/
}
