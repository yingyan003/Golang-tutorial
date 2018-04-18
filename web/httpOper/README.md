# 1.用go模拟http请求 
# 2.json标签在http服务器与客户端间的使用
------------------------------------------------

> 前言：本demo信息量较大，读者可分模块或按需阅读

## 模块1：用go模拟http请求 

#### 最简单直观的demo：client.go中的OriginHttpReq（）方法
该方法直观展示了用go如何封装http的请求。读者只看这个就够了。其他的是笔者自己为了扩展而在此基础上封装的。

#### 笔者自己封装后的demo

structObj包中的httpClient封装了go的http请求.操作主要是2步：<br>

1.建立请求
```
http.NewRequest(method, url, body)
```
2.发出请求
```
client.Do(hc.Req)
```
<b>踩坑记</b><br>
文件上传
* Content-Type必须为"multipart/form-data；boundary=30个字母的随机串"。
* 不设默认为空，报错
* 只设Req.Header.Set("Content-Type","multipart/form-data")也报错，因缺少boundary

笔者在这里为POST方式默认设置Content-Type（可覆盖）
```
if method == "POST" && hc.Req.Header.Get("contentType") == "" {
		w := new(multipart.Writer)
		contentType := w.FormDataContentType()
		hc.Req.Header.Set("Content-Type", contentType)
	}
```

这里只测试了上传，关键代码在client.go中，由GetReqBody方法模拟form表单组装了http请求体，然后由PostFile方法发出请求并处理结果。<br>

## 模块2：json标签在http服务器与客户端间的使用

所谓json标签形如下面👇的`json:"bucket"`
```
type FileSummary struct {
	Bucket string   `json:"bucket"`
	Files  []string `json:"files"`
}
```
本模块的意图为了测试json标签如何在http的服务器与客户端的数据传递间使用。<br>

structObj包中的jsonObj.go中定义了服务器与客户端交互所涉及的所有struct对象。client.go中用JsonTagTest1-3测试了特定struct下json解析的情况，具体情况请参看代码。