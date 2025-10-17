package client

import (
	"fmt"
	"net"
)

func ClientStart(ipServer string, port string) {

	addressServ := ipServer + ":" + port

	dirTCP, _ := net.ResolveTCPAddr("tcp4", addressServ)
	fmt.Println("Conectando...")
	socketC, _ := net.DialTCP("tcp4", nil, dirTCP)
	fmt.Println("Conectado al server...", socketC.RemoteAddr())

}
