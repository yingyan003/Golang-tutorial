package main

import (
	"net/http"
	"io"
)

func main() {

	//todo pattern的说明（这应该与语言无关，是http协议的规则）
	//1.pattern为"/"时，无论什么url都匹配（如"/","/zxy","/a/b"等），都会执行handler（todo 注意"/"）
	//2.不为/时，如"/a"，那只有"/a"才会匹配
	//3.第一次访问：无论什么url都执行handler，往浏览器里设置cookie，
	//  第二次访问：只有url（r.URL.Path）为"/zxy"才能访问cookie，因为该cookie的path是"/zxy"
	http.HandleFunc("/", Cookie)
	http.ListenAndServe(":8080", nil)
}

func Cookie(w http.ResponseWriter, r *http.Request) {

	ck := &http.Cookie{
		Name:   "mycookie",
		Value:  "hello",
		Path:   "/zxy",
		Domain: "localhost",
		//有效时间，单位：秒，负数表示删除cookie，0表示不设置cookie
		MaxAge: 1200,
	}

	//往响应设置cookie
	http.SetCookie(w, ck)

	//r.Cookie("mycookie")读取当前url下的key对应cookie值。
	//所以在设置cookie的内容时，cookie路径（Path: "/zxy"）的设置很关键，
	//如果url与cookie的path不匹配就获取不到
	//fmt.Println(r.URL.Path)

	//读取请求request的cookie
	ck2, err := r.Cookie("mycookie")
	if err != nil {
		//err.Error()把error类型转换为string类型
		io.WriteString(w, "cookie没找到"+err.Error())
		return
	}

	//以下注释中"go:"表示该注释是针对go语言特性的说明
	//todo go：对nil取值会产生panic:invalid memory address or nil pointer dereference
	//如果ck为nil的话，取value会panic
	//fmt.Println(ck)

	//todo go：执行写响应io.WriteString()就直接返回了,不再往下执行了
	//ck2.Value把Cookie类型转化为string类型
	io.WriteString(w, ck2.Value)
}

//该方法测试cookie的value中有中文的形式
func Cookie2(w http.ResponseWriter, r *http.Request) {
	ck := &http.Cookie{
		Name:   "mycookie",
		Value:  "hello go语言lang", //注意这里有中文
		Path:   "/zxy",
		Domain: "localhost",
		//有效时间，单位：秒，负数表示删除cookie，0表示不设置cookie
		MaxAge: 1200,
	}

	//1.该方法执行时，如果cookie的value中含有中文时，
	//  go会把中文视为非法字符，并自动忽略中文部分,只打印英文部分。（中文解决办法有时间再研究吧）
	//2.无论通过http.SetCookie(w,ck)还是w.Header().Set("Set-Cookie",ck.String())
	//  设置cookie,http协议都会自动在浏览器中设置cookie
	//3.ck.String()将Cookie类型转换为字符串
	w.Header().Set("Set-Cookie", ck.String())

	ck2, err := r.Cookie("mycookie")
	if err != nil {
		io.WriteString(w, "cookie没找到，err："+err.Error())
		return
	}

	io.WriteString(w, ck2.Value)
}
