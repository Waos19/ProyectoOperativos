package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
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

func StartServer(cfg Config, logChan chan<- string) error {
	if len(cfg.Allowed_ips) == 0 {
		logChan <- "Error: No hay IPs permitidas en la configuración. Saliendo."
		return fmt.Errorf("no hay IPs permitidas en la configuración")
	}

	address := fmt.Sprintf("%s:%d", "0.0.0.0", cfg.Port)
	addressTCP, err := net.ResolveTCPAddr("tcp4", address)

	if err != nil {
		logChan <- fmt.Sprintf("Error fatal resolviendo dirección TCP: %v", err)
		return err
	}

	listener, err := net.ListenTCP("tcp4", addressTCP)
	if err != nil {
		logChan <- fmt.Sprintf("Error fatal escuchando en TCP: %v", err)
		return err
	}

	logChan <- fmt.Sprintf("Servidor escuchando en %s", addressTCP)
	logChan <- "Esperando conexiones..."

	for {
		socketServ, err := listener.Accept()
		if err != nil {
			logChan <- fmt.Sprintf("Error aceptando conexión: %v", err)
			continue
		}
		go HandleConnection(socketServ, cfg, logChan)
	}

}

func HandleConnection(conn net.Conn, cfg Config, logChan chan<- string) {

	defer conn.Close()
	remoteIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	if !IsAllowed(remoteIP, cfg.Allowed_ips) {
		logChan <- fmt.Sprintf("Conexión rechazada desde IP no permitida: %s", remoteIP)
		return
	}

	logChan <- fmt.Sprintf("Cliente conectado desde: %s", remoteIP)

	reader := bufio.NewReader(conn)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logChan <- fmt.Sprintf("Error crítico: No se pudo obtener el directorio home: %v", err)
		return
	}

	currentDir := homeDir

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			logChan <- fmt.Sprintf("Cliente %s desconectado (error leyendo comando): %v", remoteIP, err)
			return
		}

		command = strings.TrimSpace(command)
		if command == "" { // Ignorar entradas vacías
			continue
		}

		if command == "bye" {
			logChan <- fmt.Sprintf("Cliente %s cerró la sesión.", remoteIP)
			return
		}

		if strings.HasPrefix(command, "cd ") {
			target := strings.TrimSpace(strings.TrimPrefix(command, "cd "))
			if !filepath.IsAbs(target) {
				target = filepath.Join(currentDir, target)
			}

			if stat, err := os.Stat(target); err == nil && stat.IsDir() {
				currentDir = target // Actualizamos nuestro rastreador
				conn.Write([]byte("Directorio cambiado a: " + currentDir + "\n__END__\n"))
			} else {
				conn.Write([]byte("Error cambiando de directorio: " + err.Error() + "\n__END__\n"))
			}
			continue
		}

		logChan <- fmt.Sprintf("[%s] Ejecutando comando: %s", remoteIP, command)

		output, err := shell.RunCommand(command, currentDir)
		if err != nil {
			output += "\nError: " + err.Error()
		}

		conn.Write([]byte(output + "\n__END__\n"))
	}
}
