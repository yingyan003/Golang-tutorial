package main

import (
	"net/http"
	"fmt"
	"os"
	"golang.org/x/crypto/ssh"
	"time"
	"net"
	"log"
	"os/exec"
	"pingYCE/ping"
)

func main() {
	//todo go实现ping的4种方式
	//PingShellTest()
	http.HandleFunc("/pingWithSrc",PingWithSrc)
	http.HandleFunc("/pingShell",PingByShell)
	http.HandleFunc("/ping",PingFromLocal)
}

//需求场景：指定ping的源srcIP （src ping dst）
//实现思路：
//从本机登录到远程机器，在远程机器上执行ping命令
//友情提示：
//该需求是本人在容器云项目中遇到的，涉及kubernets相关知识，更多详情可以参考我的博客：https://blog.csdn.net/zxy_666/article/details/79958948
func PingWithSrc(w http.ResponseWriter, r *http.Request){
	//预设
	user:="root"
	pass:="yeepay.com"
	if r.Header.Get("user")!=""{
		user=r.Header.Get("user")
	}
	if r.Header.Get("pass")!=""{
		pass=r.Header.Get("pass")
	}

	srcIP:=r.Header.Get("srcIP")
	dstIP:=r.Header.Get("dstIP")
	count:=r.Header.Get("count")
	podName:=r.Header.Get("podName")
	fmt.Println("header:",srcIP,dstIP,count,user,pass,podName)

	//登录到远程
	session, err := connect(user, pass, srcIP, 22)
	if err != nil {
		fmt.Println(err)
	}
	//defer session.Close()
	fmt.Println("enter node:"+srcIP)

	session.Stdout = w
	session.Stderr = os.Stderr

	//err=session.Run("pwd ; ls ; cd zxyfile ; pwd ; ls")
	//err=session.Run("ping "+dstIP+" -c "+count)
	//fmt.Println("ping from node("+srcIP+") to node:",dstIP)
	//在远程机器执行命令
	err=session.Run("pwd; ls; kubectl exec -it "+podName+" -n configcenter bash; ls; ping -c 3 "+dstIP)
	if err!=nil{
		fmt.Println("err:",err)
	}
	//fmt.Println("kubectl exec done")
	//err=session.Run("pwd ; ls")
	//if err!=nil{
	//	fmt.Println("err:",err)
	//}
	//err=session.Run("ping "+dstIP+" -c "+count)
	//fmt.Println(podName+" ping "+dstIP+" -c "+count)
	if err!=nil{
		fmt.Println("err:",err)
		w.Write([]byte("失败："+srcIP+" ping "+dstIP+" 不通"))
	}else{
		w.Write([]byte("成功："+srcIP+" ping "+dstIP+" 通"))
	}
}

//打开模拟shell终端，在终端直接执行ping命令
//必须指定-c请求次数，否则程序一去不复返了
func PingByShell(w http.ResponseWriter, r *http.Request){
	check := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}

	//获取远程登录连接客户端
	client, err := ssh.Dial("tcp", "10.151.32.27:22", &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.Password("yeepay.c0m")},
		Timeout: 30 * time.Second,
		//需要验证服务端，不做验证返回nil就可以，点击HostKeyCallback看源码就知道了
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	check(err, "dial")

	session, err := client.NewSession()
	check(err, "new session")
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", 25, 100, modes)
	check(err, "request pty")

	//启动shell终端
	err = session.Shell()
	check(err, "start shell")

	err = session.Wait()
	check(err, "return")
}

//在本机直接发出ping（默认从本机发出）
func PingFromLocal(w http.ResponseWriter, r *http.Request){
	dstIP:=r.Header.Get("dstIP")
	cmd := exec.Command("ping","-c","3", dstIP)
	cmd.Stdout = w
	err:=cmd.Run()
	if err!=nil{
		w.Write([]byte(err.Error()))
	}
}

//直接实现ping命令原理（ping使用ICMP协议实现的）
//移步pingICMP.go

//远程登录连接
func connect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		//需要验证服务端，不做验证返回nil就可以，点击HostKeyCallback看源码就知道了
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

//----------------我是分割线-----------------

//上面的函数都是http服务器的角色，需要http请求访问才能查看效果
//这里提供可直接调用函数（对应上面的PingShell）
func PingShellTest(){
	check := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}

	client, err := ssh.Dial("tcp", "172.21.0.84:22", &ssh.ClientConfig{
		User: "远程机器登录账户名",
		Auth: []ssh.AuthMethod{ssh.Password("远程机器登录密码")},
		Timeout: 30 * time.Second,
		//需要验证服务端，不做验证返回nil就可以，点击HostKeyCallback看源码就知道了
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	check(err, "dial")

	session, err := client.NewSession()
	check(err, "new session")
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", 25, 100, modes)
	check(err, "request pty")

	err = session.Shell()
	check(err, "start shell")

	err = session.Wait()
	check(err, "return")
}

