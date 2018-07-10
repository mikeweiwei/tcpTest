package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

//缓冲区大小
const BUFFSIZE = 1024

//定义消息缓冲区
var buff = make([]byte, BUFFSIZE)

func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("错误：%s\n", err.Error())
}

type Client_model struct {
	Ip     string `json:"ip"`
	Conn   *net.TCPConn
	Status bool `json:"status"`
}

var allClient map[string]Client_model

var buff2 = make([]byte, BUFFSIZE)

func Cmd() {

	for {

		reader := bufio.NewReader(os.Stdin)
		n, _ := reader.Read(buff2)

		if n > 0 {
			cmd := string(buff2[:n-1])
			cmds := strings.Split(cmd, " ")
			if len(cmds) != 0 {
				if cmds[0] == "status" {
					for k, v := range allClient {
						fmt.Println("name:" + k)
						bytes, _ := json.Marshal(v)
						fmt.Println(string(bytes[:]))
						fmt.Println("-------------------------------")
					}
				} else if cmds[0] == "open" {
					for i := 1; i < len(cmds); i++ {
						if model, ok := allClient[cmds[i]]; ok {
							if model.Status == false {
								model.Status = true
								model.Conn.Write([]byte("open"))
								allClient[cmds[i]] = model
								n, _ := model.Conn.Read(buff)
								if n > 0 {
									fmt.Printf("%s发来消息：%s", cmds[i], string(buff[:n]))
								}
							} else {
								fmt.Println(cmds[i] + "已开启！请勿重复操作！")
							}
						}
					}
					fmt.Println("-------------------------------")
				} else if cmds[0] == "stop" {
					for i := 1; i < len(cmds); i++ {
						for i := 1; i < len(cmds); i++ {
							if model, ok := allClient[cmds[i]]; ok {
								if model.Status == true {
									model.Status = false
									allClient[cmds[i]] = model
									n, _ := model.Conn.Read(buff)
									if n > 0 {
										fmt.Printf("%s发来消息：%s", cmds[i], string(buff[:n]))
									}
								} else {
									fmt.Println(cmds[i] + "已关闭！请勿重复操作！")
								}
							}
						}
					}
					fmt.Println("-------------------------------")
				} else {
					fmt.Println("请重新输入！")
				}
			}

		}
	}
}

func main() {
	port := "8080"
	allClient = make(map[string]Client_model)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:"+port)
	handleError(err)
	if err != nil {
		return
	}
	go Cmd()
	tcpListener, err := net.ListenTCP("tcp4", tcpAddr)
	handleError(err)
	if err != nil {
		return
	}

	fmt.Printf("启动监听，等待链接！\n")

	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		handleError(err)
		if err != nil {
			return
		}
		n, err := tcpConn.Read(buff)
		handleError(err)
		if n > 0 {
			if err != nil {
				tcpConn.Write([]byte("注册失败"))
			}
			tcpConn.Write([]byte("注册成功"))
			name := string(buff[:n])
			ip := tcpConn.RemoteAddr().String()
			allClient[name] = Client_model{Ip: ip, Conn: tcpConn, Status: false}
			fmt.Printf("客户端：%s 已连接！ip: %s \n", string(buff[:n]), tcpConn.RemoteAddr().String())
		}

		defer tcpConn.Close()

		//go handleConn(tcpConn, tcpConn.RemoteAddr().String())
	}

}
