package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func Connect(ipServer string, port string) (net.Conn, error) {

	addressServ := ipServer + ":" + port

	dirTCP, err := net.ResolveTCPAddr("tcp4", addressServ)
	if err != nil {
		return nil, fmt.Errorf("error resolviendo dirección: %v", err)
	}

	socketC, err := net.DialTCP("tcp4", nil, dirTCP)
	if err != nil {
		return nil, fmt.Errorf("error conectando al servidor: %v", err)
	}

	return socketC, nil
}

func StartClientReader(conn net.Conn, outputChan chan<- string, errorChan chan<- string) {
	defer conn.Close()

	serverReader := bufio.NewReader(conn)

	for {
		var response strings.Builder

		for {
			line, err := serverReader.ReadString('\n')
			if err != nil {
				errorChan <- fmt.Sprintf("Servidor desconectado o error: %v", err)
				return
			}

			lineTrimmed := strings.TrimSpace(line)
			if lineTrimmed == "__END__" {

				break
			}
			response.WriteString(line)
		}

		outputChan <- response.String()
	}
}
