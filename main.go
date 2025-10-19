package main

import (
	"bufio"
	"fmt"
	"os"
	"proyoper/internal/auth"
	"proyoper/internal/client"
	"proyoper/internal/server"
	"strings"

	"golang.org/x/term"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Bienvenido ===")
	fmt.Println("Elige modo: 'server' o 'client'")
	mode, _ := reader.ReadString('\n')
	mode = strings.TrimSpace(mode)

	if mode == "server" {
		cfg, err := server.LoadConfig("configs/server.conf")
		if err != nil {
			fmt.Println("Error cargando configuración:", err)
			return
		}
		server.StarServer(cfg)

	} else if mode == "client" {
		fmt.Print("Ingresa IP: ")
		ip, _ := reader.ReadString('\n')
		ip = strings.TrimSpace(ip)

		fmt.Print("Ingresa puerto: ")
		port, _ := reader.ReadString('\n')
		port = strings.TrimSpace(port)

		cfg, err := server.LoadConfig("configs/server.conf")
		if err != nil {
			fmt.Println("Error cargando configuración:", err)
			return
		}

		users, err := auth.LoadUsers("configs/users.db")
		if err != nil {
			fmt.Println("Error cargando usuarios:", err)
			return
		}

		attempts := 0
		for {
			fmt.Print("Usuario: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			fmt.Print("Contraseña: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				fmt.Println("Error leyendo la contraseña:", err)
				return
			}
			password := strings.TrimSpace(string(bytePassword))

			if auth.VerifyLogin(username, password, users) {
				fmt.Println("Login exitoso.")
				break
			}

			attempts++
			if !auth.HandleFailedAttempt(attempts, cfg.Max_attempts) {
				fmt.Println("Demasiados intentos fallidos. Saliendo.")
				return
			}
		}

		var interval int
		fmt.Print("Ingresa intervalo en segundos para el monitor del sistema: ")
		fmt.Scanln(&interval)

		client.ClientStart(ip, port, interval)

	} else {
		fmt.Println("Modo inválido. Debes elegir 'server' o 'client'.")
	}
}
