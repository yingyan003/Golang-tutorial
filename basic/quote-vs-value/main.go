package main

import (
	"net/http"
	"fmt"
	"Golang-tutorial/basic/quote-vs-value/getUpload"
	"strings"
)

type Stu struct {
	Name string
	Age  int
}

func main() {
	Test1()
	//Test2()
}

func Test1() {
	//todo 测试1
	//目的:测试在函数调用中，参数类型为值传递和指针传递（引用）时，若修改了参数是否对原变量产生影响。并总结go常用类型中哪些会引用传递
	//结论:
	//1.基本类型都是值传递。基本类型包括整型，浮点型，复数，字符串，常量
	//2.复合类型多是指针传递。如数组，切片，Map,都是指针传递，结构体struct可值可指针，取决于传递的是结构体本身（值：struct）还是结构体指针（指针：*struct）
	//补充：
	//接口类型作为函数参数的情况见test2()

	//值传递，对函数参数的改变不会影响原变量
	str := "string"
	ToUpper(str)
	fmt.Println(str)

	//指针传递，会改变原变量
	sl := []string{"hello", "world"}
	ToUpper(sl)
	fmt.Println(sl)

	//指针传递，会改变原变量
	stu := new(Stu)
	stu.Name = "name"
	stu.Age = 11
	ToUpper(stu)
	fmt.Println(stu.Name, stu.Age)

	//值传递，不会改变原变量
	var st Stu = Stu{
		Name: "n",
		Age:  22,
	}
	ToUpper(st)
	fmt.Println(st.Name, st.Age)
}

func ToUpper(it interface{}) {
	//it是string类型
	if a, ok := it.(string); ok {
		a = strings.ToUpper(a)
	}

	//it是[]string切片类型
	if a, ok := it.([]string); ok {
		for j := 0; j < len(a); j++ {
			a[j] = strings.ToUpper(a[j])
		}
	}

	//it是*Stu类型
	if a, ok := it.(*Stu); ok {
		a.Name = strings.ToUpper(a.Name)
		a.Age = 0
	}

	//it是Stu类型
	if a, ok := it.(Stu); ok {
		a.Name = "HAHA"
		a.Age = -1
	}
}

func ToLower(slice []string) {
	//若s:=range slice，此时，s会自动赋值为下标，是int类型，并不是迭代的内容值，需要_,s:=range slice获取切片或数组的值
	for _, s := range slice {
		strings.ToLower(s)
	}
}

func Test2() {
	http.HandleFunc("/upload", Upload)
	http.ListenAndServe(":8000", nil)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		fmt.Println(err)
	}

	//返回multipart.File类型的文件流
	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		fmt.Println(err)
	}

	if header == nil {
		fmt.Println("解析文件失败")
	}

	//todo 测试2
	//目的：文件以流的形式传输时（作为函数调用的参数等），验证文件流的完整性（即文件流有没有丢失数据或被修改等）
	//方式：
	//1.直接传r *http.Request方式
	//2.传io.Reader(接口)方式
	//结果：
	//1.无论被调用与调用的函数与是否在同一个package，并不对调用参数参数产生影响。
	//  而实际可能产生影响的地方在于，传递的参数是值类型，引用类型，还是接口类型。
	//2.值传递，不会影响原来的变量；
	//  引用(指针)传递，会影响原变量，因为传递的是原变量的地址；
	//  接口传递，至少这里是影响的。
	//3.*http.Request是指针传递，在被调用函数中，引用了仍然是同一个http请求，对request的修改全局生效。
	//4.FormFile()函数返回multipart.File接口类型的文件流，因为该接口含有io.Reader接口，故可以传递给io.Reader类型的参数。
	//  且对该io.Reader类型参数的操作会直接影响原文件流。当然，接口作为被调用参数类型情况时（如这里的io.Reader），
	//  调用参数类型（如这里的multipart.File）不非得是接口类型，也可以是struct对象，只要该对象实现了被调用参数类型的接口
	// （如这里的io.Reader），也一样会影响原来的数据。
	//5.补充一句，接口没有指针的说法，指针只是针对实体对象（如struct对象），取其地址。
	getUpload.Upload(r, file, header.Filename)
}
