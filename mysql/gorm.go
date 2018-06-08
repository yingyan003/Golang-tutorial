package main

import (
	_ "github.com/go-sql-driver/mysql" //todo 必须import，否则连接出错
	"github.com/jinzhu/gorm"
	"fmt"
	"time"
	"strings"
	"encoding/json"
)

//不显示指定变量与表字段的映射时，默认变量首字母小写，其他字母大写时在前边加_对应数据库字段。如Name-name,OrgName-org_Name
// 查询时，字段自动映射到名字匹配的变量中，其他不匹配的变量忽略，保留变量0值（如：bool类型是false，int类型是0）
// 插入时，如果有一个变量与字段名不匹配，报错
// 若要指定数据库字段名，可以用'column:pool'
type Rbd struct {
	Id      int32
	Name    string
	OrgName string `gorm:"column:orgName"`
	Used    bool
	Size    int32 //单位：G
	Pool    string `json:"pool" gorm:"column:pool;default:rbd-zxy"`
	//变量的0值会写入表中，如果变量不被赋值，即时在建表时设置字段默认值，也会被变量的0值覆盖。
	//可以在这用gorm设置默认值，且建表时有默认值，此时0值不会覆盖字段的值，取表的默认值为准
	ReadOnly bool   `gorm:"column:readOnly;default:true"`
	FsType   string `gorm:"column:fsType;default:xfs"`
	//创建时，如果不指定gorm:"-"（即忽略该变量），且不给该变量赋值，创建时会报参数错误。
	//可以把类型变为*time.Time指针，但是插入的记录，时间字段为空，就算表本身设置了该字段的默认值为当前时间，也会被nill覆盖。
	//解决办法：
	//加上`gorm:"-"`即可，此时时间类型是否指针都无所谓.
	//将时间改为string类型，与time.Time类型一样，必须被赋值，否则报错：Incorrect datetime value: '' for column 'createtime' at row 1
	//string与time.Time都可以接收查询结果：
	// string接收结果：2018-05-14T12:06:34+08:00。
	// (*)time.Time接收结果：2018-05-14 12:06:34 +0800 CST
	Createtime time.Time `gorm:"-"`
	Updatetime time.Time `gorm:"-"`
	Operator   string
	DcIdList []int32
}

func main() {
	db := Conn()
	//QueryByName(db)
	//Insert(db)
	//QueryByOrgName(db)
	//fmt.Println("----------sort")
	//SortRbds(rbds)
	//RbdNameCheck()
	//getJson()
	UpdateByNameAndDcId(db)
}

func Conn() *gorm.DB {
	user := "root"
	pass := "root"
	host := "10.151.160.27"
	database := "yce"

	db, err := gorm.Open("mysql", user+":"+pass+"@tcp("+host+":3306)/"+database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("open mysql failed. err:", err)
	}
	return db
}

func QueryByName(conn *gorm.DB) {
	name := "rbd3"
	rbd := new(Rbd)
	db := conn.Table("rbd").Where("name = ?", name).Find(rbd)
	if db.Error != nil {
		fmt.Println("queryByName failed. err=", db.Error)
	}
	fmt.Println(*rbd)
	data,_:=json.Marshal(rbd)
	fmt.Println(string(data))
}

func Insert(conn *gorm.DB) {
	rbd := new(Rbd)
	rbd.Name = "rbd11"
	rbd.OrgName = "cf"
	rbd.Used = true
	rbd.Size = 20
	rbd.Operator = "zxy"
	//rbd.ReadOnly=true
	//不指定table时，根据Create的参数名+s(如：rbds)作为表名
	db := conn.Table("rbd").Create(&rbd)
	if db.Error != nil {
		fmt.Println("insert failed. err=", db.Error)
		//Error 1062: Duplicate entry 'rbd6' for key 'name'
	}
}

func QueryByOrgName(conn *gorm.DB) []*Rbd {
	orgName := "b"
	//[]*rbd和[]rbd类型均可
	rbds := make([]*Rbd, 0)
	//todo where里的 " ？ " 不能漏掉，否则报错
	db := conn.Table("rbd").Where("orgName=?", orgName).Find(&rbds)
	if db.Error != nil {
		fmt.Println("QueryByOrgName failed. err=", db.Error)
		//Error 1062: Duplicate entry 'rbd6' for key 'name'
	}
	fmt.Println("len rbds:\n", len(rbds))
	for i, r := range rbds {
		fmt.Println("r i", i, "： ", *r)
	}

	return rbds
}

func SortRbds(rbds []*Rbd) {
	r := new(Rbd)

	//冒泡排序1
	/*
	flag := true
	for i := 0; i < len(rbds)-1; i++ {
		if flag{
			for j := len(rbds) - 2; j >= i; j-- {
				flag=false
				if rbds[j].Createtime.Before(rbds[j+1].Createtime) {
					r = rbds[j]
					rbds[j] = rbds[j+1]
					rbds[j+1] = r
					flag=true
				}
			}
		}else{
			break
		}
	}*/

	//冒泡排序2（换个冒泡的方向）
	for i := len(rbds) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if rbds[j].Createtime.Before(rbds[j+1].Createtime) {
				r = rbds[j]
				rbds[j] = rbds[j+1]
				rbds[j+1] = r
			}
		}
	}

	fmt.Println("len rbds:\n", len(rbds))
	for i, rbd := range rbds {
		fmt.Println("r i", i, "： ", *rbd)
	}
}

func RbdNameCheck() {
	n := ""
	o := "yce"
	fmt.Println(strings.Split(n, "-"))
	//s:=strings.Split(n,"-")
	//if len(s)==1||!(strings.Split(n,"-")[0]==o){
	//	fmt.Println("failed")
	//}
	if strings.HasPrefix(n, o+"-")&&len(n)>len(o+"-"){
		fmt.Println("ok")
	}else{
		fmt.Println("failed",len(n),len(o+"-"))

	}
}

func getJson(){
	r:=new(Rbd)
	r.Name="rbd1"
	r.OrgName="yce"
	r.Size=20
	r.Operator="zxy"
	dcs:=[]int32{1,2}
	r.DcIdList=dcs
	data,_:=json.Marshal(r)
	fmt.Println(string(data))
}

func UpdateByNameAndDcId(conn *gorm.DB){
	name:="default-123"
	dcId:=216
	db:=conn.Table("rbd").Where("name = ? and dcId = ?",name,dcId).Update("used",true)
	if db.Error != nil {
		fmt.Println("UpdateByUsed failed. err=", db.Error)
	}
}