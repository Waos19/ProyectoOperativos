package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"proyoper/internal/graphic"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Bienvenido ===")
	fmt.Println("Elige modo: 'server' o 'client'")
	fmt.Print("-> ")
	mode, _ := reader.ReadString('\n')
	mode = strings.TrimSpace(mode)

	switch mode {
	case "server":

		graphic.StartServerUI()

	case "client":

		graphic.StartClientUI()

	default:
		fmt.Println("Modo inválido. Debes elegir 'server' o 'client'.")
	}
}
