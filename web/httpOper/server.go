package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"io"
	"Golang-tutorial/web/httpOper/structObj"
	"encoding/json"
	"time"
)

func main() {
	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/jsonTagTest1", JsonTagServer1)
	http.HandleFunc("/jsonTagTest2", JsonTagServer2)
	http.ListenAndServe(":8001", nil)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Content-Type:", r.Header.Get("Content-Type"))

	err := r.ParseMultipartForm(100000)
	if err != nil {
		fmt.Println("ParseMultipartForm error:", err) //unexpected EOF
		return
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("FormFile error:", err)
		return
	}
	defer f.Close()

	dir := os.TempDir()

	//todo 惨痛的踩坑记
	//1、h.Filename=/Users/zxy/Desktop/3.png277875241（本以为返回的仅仅的纯文件名，忽略了人家返回的是文件全路径名）
	//"zxy-"+h.Filename=zxy-/Users/zxy/Desktop/3.png277875241
	//非法名称，无法创建这样的文件，f2=nil,如果f2操作就会报空指针异常
	//2、只要是返回的err的地方，一定得判断err是否为空，不空得打印err,便于排查错误
	//笔者在这里漏打了，结果这里出错了，半天没找到原因。
	//3、如果err!=nil，打印err后最好就return了，不然下面对空指针的操作会报一堆空指针异常
	//影响最初错误的排查,良好习惯一定要坚守，否则大把的时间浪费在找问题上实在是太尴尬了
	//4.千万不要以你平时正确的操作来认为同样的操作在你的程序不可能产生error，那就太天真了
	f2, err := ioutil.TempFile(dir, "zxy-")
	if err != nil {
		fmt.Println("TempFile error", err)
		return
	}
	written, err := io.Copy(f2, f)
	if err != nil {
		fmt.Printf("copy error=%s,written=%d\n", err, written)
		return
	}
	_, err = w.Write([]byte("ok"))
}

//json的tag标签 测试1
//目的：
//当服务端返回的数据类型是：ServerStatus + FileSummary
//客户端如何定义才能正确解析数据？
//(如struct内的数据定义：x变量 x类型 `json:"data"`，json的tag标签指的是`json:"data"`)
//结果：
//客户端通过ClientStatus1接受json解析，其中 Data *FileSummary `json:"data"`定义中的FileSummary，
//todo 可以是指针，也可以不是指针类型均可解析出Data中的FileSummary类型数据。
//(tag标签在服务器marshal与客户端unmarshal的过程中，通过tag解析.所以在编程中要注意解析时tag的一致性，不对应的话解析会失败的。
//如果报空指针异常，可能你其他地方写错了，或者tag不一致等，笔者就在这里栽跟头，一直以为是是否指针的问题，后来才发现是其他地方报的空指针）

func JsonTagServer1(w http.ResponseWriter, r *http.Request) {
	fs := new(structObj.FileSummary)
	fs.Bucket = "color"
	fs.Files = append(fs.Files, "green")
	fs.Files = append(fs.Files, "blue", "red")

	stat := new(structObj.ServerStatus)
	stat.Code = 0
	stat.Message = "ok"
	stat.Data = fs

	w.Write(Bytes(stat))
}

//json的tag标签 测试2
//目的：
//当服务端返回的数据类型是：ServerStatus + ObjSummary（TODO 其中有组合对象Obj）
//客户端如何定义才能正确解析数据？
//结果：
//客户端通过ClientStatus1接受json解析，其中 Data *ObjSummary `json:"data"`ObjSummary，
//todo 可以是指针，也可以不是指针类型均可解析出Data中的ObjSummary类型数据
func JsonTagServer2(w http.ResponseWriter, r *http.Request) {
	o := new(structObj.Obj)
	o.Size = 1
	o.ContentType = "multipart/form-data"
	o.LastModifyTime = time.Now()

	os := new(structObj.ObjSummary)
	os.Obj = *o
	os.Bucket = "bucket"

	stat := new(structObj.ServerStatus)
	stat.Code = 0
	stat.Message = "ok"
	stat.Data = os

	w.Write(Bytes(stat))
}

func Bytes(s *structObj.ServerStatus) []byte {
	//将数据转化成json
	data, _ := json.Marshal(s)
	return data
}
