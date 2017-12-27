package main

import (
	"time"
	"fmt"
)

func main(){
	//TimeFormatConvertor()
	//TimeAdd()
	IsTimeBefore()
}

//时间格式转换
func TimeFormatConvertor(){
	var t1 time.Time

	//返回类型：Time,
	//格式为：2017-12-25 18:51:38.54526049 +0800 CST
	t1=time.Now()

	//将Time类型转化为string
	//格式为：2017-12-25 18:51:38
	//参数必须为"2006-01-02 15:04:05",可简单记忆为123456（1月2号3点4分5秒6年）(go规定的写法,只有这样才转化为24制的日期时间格式)
	t2:=t1.Format("2006-01-02 15:04:05")

	//t2的类型是string
	fmt.Println(t2)
}

//时间相减
func TimeSub(){
	t1:=time.Now()

	for i:=0;i<10;i++{
		//模拟时间延迟
	}

	//sub函数:Now()-t1
	duration:=time.Now().Sub(t1)
	//将时间差的单位转换为秒/分/小时（还有纳秒）
	duration.Seconds()
	duration.Minutes()
	hour:=duration.Hours()


	//返回的是类型都是float64
	fmt.Println(hour)
}

//时间相加
func TimeAdd(){
	now:=time.Now()
	fmt.Println(now)

	//获取时间间隔duration，通过它与时间相加
	//参数的有效格式："ns", "us" (or "µs"), "ms", "s", "m", "h"
	//参数(string类型)可正可负，负数表示时间回退，正数表示增加
	//1分种后
	duration,_:=time.ParseDuration("1m")
	m1:=now.Add(duration)
	fmt.Println(m1)

	//1小时后
	duration,_=time.ParseDuration("1h")
	h1:=now.Add(duration)
	fmt.Println(h1)

	//1天前
	duration,_=time.ParseDuration("-24h")
	d1:=now.Add(duration)
	fmt.Println(d1)
}

//判断时间的先后
func IsTimeBefore(){
	t1:=time.Now()

	for i:=0;i<10;i++{
		//模拟时间延迟
	}

	t2:=time.Now()
	flag:=t2.After(t1)
	if flag==true{
		fmt.Println("时间t2",t2,"在t1",t1,"之后")
	}else{
		fmt.Println("时间t2",t2,"在t1",t1,"之前")
	}
}