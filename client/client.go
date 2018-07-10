package main

import (
	"net"
	"os"
	"fmt"
	"bufio"
	"flag"
)

type Data struct {
	ID     string
	Status string
	Data   string
}


const BUFFSIZE = 1024
var buff = make([]byte, BUFFSIZE)
var writebuff = make([]byte, BUFFSIZE)

func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("错误：%s\n", err.Error())
}

func main() {


	name := flag.String("name", "", "please use name")
	flag.Parse()
	port := "8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:" + port)
	handleError(err)
	tcpConn, err := net.DialTCP("tcp4", nil, tcpAddr)
	handleError(err)

	fmt.Println("Your Name is " + *name)

	//注册信息
	go func() {
		for  {
			n, err := tcpConn.Read(writebuff)
			handleError(err)
			if n > 0 {
				fmt.Println("服务端回复:", string(writebuff[:n]))
			}
			if string(writebuff[:n]) == "open"{
				tcpConn.Write([]byte("已开启！"))
			}
			if string(writebuff[:n]) == "stop"{
				tcpConn.Write([]byte("关闭！"))
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	continued := true
	inputStr := ""

	tcpConn.Write([]byte(*name))
	for(continued){
		n, err := reader.Read(buff)
		handleError(err)
		if n > 0 {
			k, _ := tcpConn.Write(buff[:n])
			if k > 0 {
				inputStr = string(buff[:n])
				fmt.Printf("发送消息：%s", inputStr)
				if inputStr == "exit\n" {
					continued = false
				}
			}
		}
	}

}