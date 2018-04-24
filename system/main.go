package main
//
//import (
//	"net/http"
//	"troubleshooting/ping"
//	"github.com/maxwell92/gokits/log"
//	"github.com/julienschmidt/httprouter"
//	"troubleshooting/netstat"
//	"github.com/kataras/iris/middleware/logger"
//	"fmt"
//)
//
//func main() {
//	router:=httprouter.New()
//	//router.GET("/ping/:dstIP/:count",Ping)
//	router.GET("/netstat/:port",Netstat)
//
//	fmt.Println("troubleshooting listen at port: 8090")
//	http.ListenAndServe(":8090", router)
//}
