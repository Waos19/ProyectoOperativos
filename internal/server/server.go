package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"proyoper/internal/monitor"
	"proyoper/internal/shell"
	"strings"
)

func IsAllowed(ip string, allowed []string) bool {
	for _, allowedip := range allowed {
		if ip == allowedip {
			return true
		}
	}
	return false
}

func StartServer(cfg Config) {
	if len(cfg.Allowed_ips) == 0 {
		log.Fatal("No hay IPs permitidas en la configuración")
	}

	address := fmt.Sprintf("%s:%d", "0.0.0.0", cfg.Port)
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
	if !IsAllowed(remoteIP, cfg.Allowed_ips) {
		fmt.Printf("Conexión rechazada desde IP no permitida: %s\n", remoteIP)
		return
	}

	fmt.Println("Cliente conectado desde:", remoteIP)

	reader := bufio.NewReader(conn)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("No se pudo obtener el directorio home:", err)
		return
	}

	if err := os.Chdir(homeDir); err != nil {
		fmt.Println("Error cambiando al directorio home:", err)
		return
	}

	currentDir := homeDir

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

		if strings.HasPrefix(command, "cd ") {
			target := strings.TrimSpace(strings.TrimPrefix(command, "cd "))
			if !filepath.IsAbs(target) {
				target = filepath.Join(currentDir, target)
			}
			if err := os.Chdir(target); err != nil {
				conn.Write([]byte("Error cambiando de directorio: " + err.Error() + "\n__END__\n"))
			} else {
				currentDir = target
				conn.Write([]byte("Directorio cambiado a: " + currentDir + "\n__END__\n"))
			}
			continue
		}

		if command == "repstats" {
			stats, err := monitor.GenerateReport()
			if err != nil {
				conn.Write([]byte("Error generando reporte: " + err.Error() + "\n__END__\n"))
				continue
			}
			conn.Write([]byte(string(stats) + "\n__END__\n"))
			continue
		}

		fmt.Println("Ejecutando comando:", command)
		output, err := shell.RunCommand(command)
		if err != nil {
			output += "\nError: " + err.Error()
		}

		conn.Write([]byte(output + "\n__END__\n"))
	}
}
