package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func ClientStart(ipServer string, port string) {

	addressServ := ipServer + ":" + port

	dirTCP, _ := net.ResolveTCPAddr("tcp4", addressServ)
	fmt.Println("Conectando...")
	socketC, _ := net.DialTCP("tcp4", nil, dirTCP)
	fmt.Println("Conectado al server...", socketC.RemoteAddr())
	RunInteractiveSession(socketC)
}

func RunInteractiveSession(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error leyendo entrada:", err)
			break
		}

		command := strings.TrimSpace(input)
		if command == "bye" {
			fmt.Println("Cerrando conexión...")
			conn.Write([]byte("bye\n"))
			break
		}

		_, err = conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println("Error enviando comando:", err)
			break
		}

		// Leer múltiples líneas hasta encontrar el delimitador __END__
		for {
			line, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Error leyendo respuesta:", err)
				break
			}

			line = strings.TrimSpace(line)
			if line == "__END__" {
				break
			}

			fmt.Println(line)
		}
	}
}
