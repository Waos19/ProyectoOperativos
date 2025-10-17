package main

import (
	"bufio"
	"fmt"
	"os"
	"proyoper/internal/client"
	"proyoper/internal/server"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Bienvenido ===")
	fmt.Println("Elige modo: 'server' o 'client'")
	mode, _ := reader.ReadString('\n')
	mode = strings.TrimSpace(mode)

	fmt.Print("Ingresa IP: ")
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)

	fmt.Print("Ingresa puerto: ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)

	if mode == "server" {
		server.StarServer(ip, port)
	} else if mode == "client" {
		client.ClientStart(ip, port)
	} else {
		fmt.Println("Modo inválido. Debes elegir 'server' o 'client'.")
	}

}
