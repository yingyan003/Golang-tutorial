package main

import (
	"bytes"
	"mime/multipart"
	"fmt"
	"os"
	"io"
	"Golang-tutorial/web/httpOper/structObj"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func main() {
	//filePath := "/Users/zxy/Desktop/3.png"
	//PostFile(filePath, "http://localhost:8001/upload")

	JsonTagTest1("http://localhost:8001/jsonTagTest1")
	JsonTagTest2("http://localhost:8001/jsonTagTest2")
	JsonTagTest3("http://localhost:8001/jsonTagTest2")

}

//上传文件
func PostFile(fname, targetUrl string) {
	bodyBuf, bodyWriter, err := GetReqBody(fname)
	if err != nil {
		fmt.Println(err)
		return
	}

	//TODO 惨痛的踩坑记
	//TODO 必须在发出请求前关闭bodyWriter，反则服务端ParseMultipartForm解析文件时，会报unexpected EOF错误
	//todo 切记是发出请求前关闭，所以这里不能使用 defer bodyWriter.Close()，那就是在发出请求前没关闭，同样会报错
	bodyWriter.Close()

	hc := structObj.NewHCImpl()

	//todo 坑！文件上传，Content-Type必须为"multipart/form-data；boundary=30个字母的随机串"
	//1.不设默认为空，报错
	//2.只设Req.Header.Set("Content-Type","multipart/form-data")也报错，缺少boundary
	//3.已在hc.NewRequest中默认赋值
	res, err := hc.NewRequest("POST", targetUrl, bodyBuf)
	if err != nil {
		fmt.Println("NewRequest error:", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Errorf("bad status: %s", res.Status)
	}

	//将io.Reader类型转换为string打印出来：
	//1.先将io.Reaser转化为[]byte，再将[]byte转化为string
	body, err := ioutil.ReadAll(res.Body)
	strbody := string(body)

	fmt.Println("header", res.Header, "\nbody", strbody)
}

//用于组装上传文件的http请求体
func GetReqBody(fname string) (*bytes.Buffer, *multipart.Writer, error) {
	//实现了io.Writer接口
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//第一个参数对应form表单的name字段
	fileWriter, err := bodyWriter.CreateFormFile("file", fname)
	if err != nil {
		fmt.Println("error writting to buffer")
		return nil, nil, err
	}

	f, err := os.Open(fname)
	if err != nil {
		fmt.Println("open file error")
		return nil, nil, err
	}
	defer f.Close()

	_, err = io.Copy(fileWriter, f)
	if err != nil {
		fmt.Println("copy error")
		return nil, nil, err
	}

	return bodyBuf, bodyWriter, nil
}

//测试1
//方式：客户端用于unmarshal的struct中Data的类型是指针：*FileSummary `json:"data"`
//json解析结果：成功解析出标签为`json:"data"`中的FileSummary结构体数据
//clientSatus1 0 ok color [green blue red]
func JsonTagTest1(targetUrl string) {
	hc := structObj.NewHCImpl()
	res, err := hc.NewRequest("GET", targetUrl, nil)
	if err != nil {
		fmt.Println("NewRequest error:", err)
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	stat1 := new(structObj.ClientStatus1)
	json.Unmarshal(body, stat1)
	fmt.Println("clientSatus1", stat1.Code, stat1.Message, stat1.Data.Bucket, stat1.Data.Files)
}

//测试2
//方式：服务器的struct中有组合对象（Obj），客户端同理设置，且用于unmarshal的struct中Data的类型是（非指针）：ObjSummary `json:"data"`
//json解析结果：成功解析组合对象（Obj），以及ObjSummary结构体对象
//clientSatus2 0 ok bucket {1 multipart/form-data 2018-01-07 21:41:03.619011841 +0800 CST}
//clientSatus2 1 multipart/form-data multipart/form-data
func JsonTagTest2(targetUrl string) {
	hc := structObj.NewHCImpl()
	res, err := hc.NewRequest("GET", targetUrl, nil)
	if err != nil {
		fmt.Println("NewRequest error:", err)
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	stat := new(structObj.ClientStatus2)
	json.Unmarshal(body, stat)
	fmt.Println("clientSatus2", stat.Code, stat.Message, stat.Data.Bucket, stat.Data.Obj)
	fmt.Println("clientSatus2", stat.Data.Obj.Size, stat.Data.Obj.ContentType, stat.Data.ContentType)
}

//测试3
//目的：当用于unmarshal的Data `json:"data"`是interface{}类型时，解析会如何
//结果：
//ClientStatus3中只有一个data变量，在json unmarshal解析时仍可解析，
//1.服务端返回的其他json（如ServerStatus中的`json:"code"`，ObjSummary中的`json:"bucket"`等），因在解析时找不到对应的tag，故忽略
//2.data `json:"data"`在解析时成功找到了Obj中的tag，故Obj解析成功.
//  todo 但由于客户端的data是interface{}类型，故在编译时（即编码时），无法逐一访问Obj中的变量，只能操作data本身
// （说白了就是打印data本身，不能单独操作Obj中的逐个变量）
func JsonTagTest3(targetUrl string) {
	hc := structObj.NewHCImpl()
	res, err := hc.NewRequest("GET", targetUrl, nil)
	if err != nil {
		fmt.Println("NewRequest error:", err)
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	stat := new(structObj.ClientStatus3)
	json.Unmarshal(body, stat)
	fmt.Println("clientSatus3", stat.Data)
}

//本函数是该httpOper例子的原型，是最直接也是最直观的go http请求的封装。
//该httpOper例子在此函数的基础上进行了封装。
//函数贴在这里供需要快速回顾时使用。
func OriginHttpReq() {
	filePath := "/Users/zxy/Desktop/3.png"

	//实现了io.Writer接口
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Println("error writting to buffer", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("open file error", err)
	}
	defer f.Close()

	_, err = io.Copy(fileWriter, f)
	if err != nil {
		fmt.Println("copy error", err)
	}

	//TODO 必须关闭，否则服务端ParseMultipartForm（）函数解析时将发生EOF错误
	bodyWriter.Close()

	req, err := http.NewRequest("POST", "http://localhost:8001/upload", bodyBuf)
	if err != nil {
		fmt.Println("NewRequest error", err)
	}

	//模拟表单上传的enctype="multipart/form-data"
	contentType := bodyWriter.FormDataContentType()
	req.Header.Set("Content-Type", contentType) //"multipart/form-data"
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do error", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		fmt.Println(err)
	}

	//将io.Reader类型转化为string打印出来
	//先将io.Reader转化为[]byte，再将[]byte转化为string
	body, err := ioutil.ReadAll(res.Body)
	str := string(body)

	fmt.Println("ok")
	fmt.Println("header", res.Header, "\nbody", str)
}
