package md5Sum

import (
	"io"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"os"
)

type FileMd5Sum struct {
	FileName string
	//默认值为""
	Md5Sum   string
	//默认值为nil
	Reader   io.Reader
}

//根据文件名创建对象。当然，文件名可以为空（此处的文件名要求是文件的路径名，即可以正常打开的文件路径）
func NewFileMd5Sum(fileName string, reader io.Reader) *FileMd5Sum {
	f := new(FileMd5Sum)
	f.FileName = fileName
	f.Reader = reader
	return f
}

//计算文件的md5，提供文件存在的两种方式：1.现有文件 2.文件流（如文件上传时FormFile（）函数返回的类型）
func (f *FileMd5Sum) Sum() []byte {
	//计算现有文件的md5码
	//当文件名不为空时，用于计算现有文件的MD5。
	var err error
	if f.FileName != "" {
		//打开现有文件，赋值给f.Reader（因为Open()函数的返回值类型*File实现Reader接口，故可赋值Reader类型的变量）
		f.Reader, err = os.Open(f.FileName)
	}
	if err != nil {
		fmt.Println("open file error: file=%s, err=%s", f.FileName, err)
		return []byte("")
	}

	//计算文件流stream的md5码
	//当文件名为空时，即计算文件流stream的md5码，此时该文件流通过创建FileMd5Sum实例对象时赋值
	if f.Reader == nil {
		fmt.Println("file stream can not be nil")
	}

	//以下4行代码就是就算MD5的关键代码
	//go sdk自身已经替我们封装了md5的算法，我们只需调用就OK
	hash := md5.New()
	io.Copy(hash, f.Reader)
	sum := hash.Sum([]byte(""))
	f.Md5Sum = hex.EncodeToString(sum)
	return sum

	//todo 此处注意一点(惨痛的踩坑史)
	//当执行io.Copy(w, r)函数时，w和r的文件指针都会移到文件末尾，当下次操作r引用的同一个文件时，
	//如打开或copy文件，会发现文件是空的。解决办法是将文件指针移回文件头。如r.Seek(0，0)
}

func (f *FileMd5Sum) String() string {
	return f.Md5Sum
}
