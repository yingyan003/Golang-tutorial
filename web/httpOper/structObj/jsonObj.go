package structObj

import "time"

type FileSummary struct {
	Bucket string   `json:"bucket"`
	Files  []string `json:"files"`
}

type Obj struct {
	Size           int64     `json:"size"`
	ContentType    string    `json:"contentType"`
	LastModifyTime time.Time `json:"lastModifyTime"`
}

type ObjSummary struct {
	//组合对象,此处不用打tag标签（对象内已打），解析json时会自动寻找匹配的标签赋值
	Obj
	Bucket string `json:"bucket"`
}

type ServerStatus struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	//服务端类型是interface{}，json tag内嵌json tag
	Data interface{} `json:"data"`
}

//服务端struct内嵌struct的情况
type ClientStatus1 struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	//todo 客户端类型是(指针与非指针均可): [*]FileSummary
	Data *FileSummary `json:"data"`
}

//服务端有组合对象（Obj）的情况
type ClientStatus2 struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	//todo 客户端类型是(指针与非指针均可): [*]ObjSummary
	Data ObjSummary `json:"data"`
}

//测试客户端是interface{}类型的接收情况
type ClientStatus3 struct {
	Data interface{} `json:"data"`
}
