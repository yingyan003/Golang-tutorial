package mount

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
)

func WriteRbd(){
	//获取当前工作目录
	fmt.Println(os.Getwd())

	//判断挂载目录是否存在
	flag,err:=IsPathExist("/xy")
	if !flag{
		fmt.Println(err)
		//不存在时创建
		MkdirTest("/Users/zxy/xy",os.ModePerm)
	}

	//打开文件。如果文件不存在则创建（由os.O_CREATE指定），文件存在则以可读写（os.O_RDWR）的模式和指定权限（ModePerm）打开
	f1, err := os.OpenFile("/Users/zxy/xy/f1", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f1.Close()

	//写文件
	err=ioutil.WriteFile(f1.Name(),[]byte("hello word"),0644)
	if err!=nil{
		fmt.Println("write err:",err)
	}

	//读文件
	data,err:=ioutil.ReadFile(f1.Name())
	if err!=nil{
		fmt.Println("read data failed.",err)
	}
	fmt.Println("read from /Users/zxy/xy/f1:\n",data)

	fmt.Println("byte"," [104 101 108 108 111 32 119 111 114 100]")
	fmt.Println("string",string(data))
}

func GetRbd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//读文件
	data,err:=ioutil.ReadFile("/Users/zxy/xy/f1")
	if err!=nil{
		fmt.Println("read data failed.",err)
	}

	w.Write(data)
}

//判断文件或目录是否存在
func IsPathExist(fileName string) (bool, error) {
	//or os.Open(fileName)
	_, err := os.Stat(fileName)
	//存在
	if err == nil {
		fmt.Println("file exist. fileName:",fileName)
		return true, nil
	}
	//不存在
	if os.IsNotExist(err) {
		fmt.Println("file isn't exist. fileName:",fileName)
		return false, nil
	}
	//不确定
	return false,err
}

//创建单级目录
func MkdirTest(path string,perm os.FileMode){
	err:=os.Mkdir(path,perm)
	if err!=nil {
		fmt.Println("err from Mkdir:",path,"\nerr:\n",err)
		return
	}
	fmt.Println("Mkdir success. dir:",path)
}

func main(){
	//todo test rbd
	//mount.WriteRbd()
}