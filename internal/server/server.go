package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"proyoper/internal/shell"
)

func IsAllowed(ip string, allowed []string) bool {
	for _, allowedip := range allowed {
		if ip == allowedip {
			return true
		}
	}
	return false
}

func StarServer(cfg Config) {
	if len(cfg.allowed_ips) == 0 {
		log.Fatal("No hay IPs permitidas en la configuración")
	}

	address := fmt.Sprintf("%s:%d", "0.0.0.0", cfg.port)
	addressTCP, err := net.ResolveTCPAddr("tcp4", address)

	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp4", addressTCP)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Servidor escuchando en", addressTCP)
	fmt.Println("Esperando conexiones...")

	for {
		socketServ, err := listener.Accept()
		if err != nil {
			log.Printf("Error aceptando conexión: %v", err)
			continue
		}
		go HandleConnection(socketServ, cfg)
	}
}

func HandleConnection(conn net.Conn, cfg Config) {

	defer conn.Close()
	remoteIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	if !IsAllowed(remoteIP, cfg.allowed_ips) {
		fmt.Printf("Conexión rechazada desde IP no permitida: %s\n", remoteIP)
		return
	}

	fmt.Println("Cliente conectado desde:", remoteIP)

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error leyendo comando:", err)
			return
		}

		command = command[:len(command)-1] // quitar salto de línea
		if command == "bye" {
			fmt.Println("Cliente cerró la sesión.")
			return
		}

		fmt.Println("Ejecutando comando:", command)
		output, err := shell.RunCommand(command)
		if err != nil {
			output += "\nError: " + err.Error()
		}

		conn.Write([]byte(output + "\n__END__\n"))
	}
}
