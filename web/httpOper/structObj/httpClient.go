package structObj

import (
	"io"
	"net/http"
	"mime/multipart"
)

type HttpClient struct {
	Req *http.Request
}

func NewHCImpl() *HttpClient {
	return new(HttpClient)
}

//method:请求的方式，如POST/GET
//url:目标url
//body:POST方式的请求体，其他方式为nil即可
func (hc *HttpClient) NewRequest(method, url string, body io.Reader) (*http.Response, error) {
	//当请求方式是POST时，contentType赋默认值
	//todo 坑！文件上传时，Content-Type必须为"multipart/form-data；boundary=30个字母的随机串"
	//1.不设默认为空，报错
	//2.只设Req.Header.Set("Content-Type","multipart/form-data")也报错，缺少boundary
	if method == "POST" && hc.Req.Header.Get("contentType") == "" {
		w := new(multipart.Writer)
		contentType := w.FormDataContentType()
		hc.Req.Header.Set("Content-Type", contentType)
	}

	var err error
	hc.Req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	//发送请求
	res, err := client.Do(hc.Req)
	return res, err
}
