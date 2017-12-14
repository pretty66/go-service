package main

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

//自己的数据类
type Excel string

var console *log.Logger

type Out struct {
	Error int
	Data  interface{}
}

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	console = log.New(logFile, "[* === *]", log.LstdFlags)

	rpc.Register(new(Excel))
	listener, err := net.Listen("tcp", ":1314")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}

func (t *Excel) ReadExcel(file *FileName, reply *interface{}) error {
	*reply = read(file)
	return nil
}

type FileName struct {
	FileName string
	UserName string
}

// 读取Excel
func read(file *FileName) *Out {
	xl, err := excelize.OpenFile(file.FileName)
	if err != nil {
		console.Println("打开文件失败：[user:"+file.UserName+" => "+file.FileName+"]", err)
		return NewOut(10001, "文件不存在")
	}
	// 解析excel后 返回数据给调用方
	data := xl.GetRows("Sheet1")
	return NewOut(0, data)
}

// 统一数据返回方法
func NewOut(errMsg int, data interface{}) *Out {
	return &Out{errMsg, data}
}
