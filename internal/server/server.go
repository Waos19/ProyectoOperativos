package server

import (
	"fmt"
	"log"
	"net"
)

func StarServer(ipServer string, port string) {
	addressTCP, err := net.ResolveTCPAddr("tcp4", ipServer+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP("tcp4", addressTCP)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Servidor escuchando en", addressTCP)
	fmt.Println("Esperando conexiones...")
	socketServ, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cliente conectado desde: ", socketServ.RemoteAddr())
}
