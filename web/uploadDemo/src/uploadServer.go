package main

import (
	"net/http"
	"fmt"
	"os"
	"io"
)

const (
	//单位：byte
	maxMemory = 100000
)

//服务器保存上传文件的目录
var uploadDir string

//分隔符
var sep = "/"

func main() {
	//当访问的路径pattern是"/upload",请求则由Upload处理器（handler）处理
	http.HandleFunc("/upload", Upload)

	//监听本机8000端口，handler为nil时由默认处理器处理，go最终将根据访问路径的pattern选择相应具体的handler来处理
	http.ListenAndServe("localhost:8000", nil)
}

func Upload(w http.ResponseWriter, r *http.Request) {

	//文件通过表单以POST方式提交
	if r.Method == "POST" {

		//处理文件上传需要调用ParseMultipartForm,参数maxMemory表示上传的文件存储在maxMemory大小的内存中，
		//如果文件大小超过maxMemory，剩余的部分存储到系统的临时文件中
		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			fmt.Println(err)
			return
		}

		//通过FormFile获取上面文件的信息，参数key对应表单input标签中type="file"的name项（同一个name可能对应多个上传的文件），
		// 且根据给定的key只返回第一个文件的信息（句柄）。最后使用io.Copy()存储文件
		fileTemp, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			fmt.Println(err)
			return
		}

		//知道包含该defer语句的函数（这里是Upload）执行完毕时，defer后的函数（file.Close()）才会被执行。
		// 无论Upload是return正常返回，还是panic异常退出。保证资源被释放
		defer fileTemp.Close()

		//将header信息返回客户端
		fmt.Fprintf(w, "fileHeader.Header:\n%v", fileHeader.Header)

		//获取当前目录的绝对路径。一般通过开发工具（比如Gogland-EAP）提供的运行按钮运行时，
		// Getwd()通常是项目根目录（这里是uploadDemo），但如果通过命令行运行时，获取的是
		//当前工作目录（linux下与pwd获取的一致）
		path, err := os.Getwd()

		//uploadDemo/doc/uploadFiles
		uploadDir = path + sep + "doc" + sep + "uploadFiles"

		//判断目录是否存在
		y, err := IsPathExists(uploadDir)
		if err != nil {
			fmt.Println(err)
			return
		}

		//目录不存在
		if !y {
			//创建目录。根据参数path，MkdirAll创建多级目录（如"/a/b",a可存在可不存在，不存在则创建），
			// Mkdir创建单级目录(如"/b"，假定"/a/b"，若a不存在则不创建，返回err)
			err = os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		//打开文件。如果文件不存在则创建（由os.O_CREATE指定），文件存在则以可读写（os.O_RDWR）的模式和指定权限（ModePerm）打开
		f2, err := os.OpenFile(uploadDir+sep+fileHeader.Filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f2.Close()

		//将FormFile返回的保存在内存中的文件复制到指定文件中
		io.Copy(f2, fileTemp)
	}
}

//判断文件或目录是否存在。通过os.Stas（）函数返回的错误值进行判断
//1.返回的错误信息为nil，说明文件或目录存在
//2.返回的错误类型使用os.IsNotExist()判断为true，说明文件或目录不存在
//3.返回的错误信息为其它类型，则不确定是否存在
func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	//存在
	if err == nil {
		return true, nil
	}
	//不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	//不确定
	return false, err
}
