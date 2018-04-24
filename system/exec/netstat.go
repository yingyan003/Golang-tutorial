package main

import (
	"fmt"
	"strconv"
	"bytes"
	"os/exec"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"time"
)

func main() {
	//提供rest api的第三方包
	router:=httprouter.New()
	router.GET("/netstat/:port",Netstat)

	fmt.Println("troubleshooting listen at port: 8090")
	http.ListenAndServe(":8090", router)
}

func Netstat(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("---------enter Netstat")
	grepPort := ""
	//获取api中参数，就是":xx"中的xx
	port := params.ByName("port")

	//校验port是否合法（即是否是数字的形式）
	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Errorf("port=%s 参数错误\n", port)
		w.Write([]byte("port参数错误"))
		return
	}

	//约定port=999表示用户不指定端口（这种方式纯粹为了省事）
	//背景：
	//1.因为httprouter包的api没有提供可选参数（就是可填可不填），
	//如"/:port"，那么port必须填写，不允许为空，否则访问不了该api。
	//2.设计需求是用户可选择是否指定该参数，鉴于httprouter的情况，
	//可以在header请求头中添加该参数，就避免了在api中这个尴尬情况。
	if portInt != 999 {
		grepPort = " | grep " + port
	}

	fmt.Println("------port", grepPort)

	//将命令执行结果写如缓存中，以便灵活操作
	var buf bytes.Buffer
	//mac对netstat的支持与Linux不完全一致
	//mac不支持-p参数，如果指定了，那么cmd.run启动会报错，且可能使程序陷入漫长等待
	//todo key code
	cmd := exec.Command("/bin/bash", "-c", "netstat -ant"+grepPort)
	cmd.Stdout = &buf
	//用于测试netstat命令是否成功启动
	s := make([]string, 0)
	go run(cmd, &s)

	time.Sleep(50 * time.Millisecond)
	if len(s) != 0 {
		w.Write([]byte("连接失败，请重试"))
		fmt.Println("return")
		return
	}

	//设置时钟，当倒计时结束时返回
	//防止因错误命令而导致程序长时间执行，无法返回（比如，在我的mac执行-p时，程序卡死了）
	timer := time.NewTimer(10 * time.Second)
	timeout := timer.C
	for {
		//select：
		//1.case的条件参数必须是管道
		//2.在一次循环中，从前到后，选择第一个管道中有数据输出的case，并执行。
		//3.2中如果没有那样的case，就执行default（前提是有default），故可以把一些业务逻辑的实现放在default中
		select {
		//倒计时结束时，返回
		case <-timeout:
			{
				w.Write([]byte("连接超时，请重试"))
				return
			}
		default:
			{
				//todo
				//当成功执行命令时，执行过程就开始写入执行结果，
				//如果不稍加等待就返回，往往执行没结束，数据没写完就返回了
				if buf.String() != "" {
					//等待数据完成写入，等待时间依据系统环境而定，mac相对Linux慢
					time.Sleep(1 * time.Second)
					w.Write([]byte(buf.String()))
					fmt.Println(buf.String())
					return
				}
			}
		}
	}
}

//todo
//因为直接用协程启动，无法获取返回参数，如：go cmd.Run()
//故引入函数，在函数中获取并处理错误信息
func run(cmd *exec.Cmd, s *[]string) {
	err := cmd.Run()
	if err != nil {
		fmt.Errorf("netstat failed. err = %s\n", err)
		*s = append(*s, "run netstat failed")
	}
}
