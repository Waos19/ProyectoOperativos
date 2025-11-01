package server

import (
	"bufio"
	"fmt"
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

// 💡 CAMBIO: La función ahora acepta un canal de logs (logChan) y devuelve un error.
func StartServer(cfg Config, logChan chan<- string) error {
	if len(cfg.Allowed_ips) == 0 {
		// 💡 CAMBIO: En lugar de log.Fatal, enviamos al canal y devolvemos un error.
		logChan <- "Error: No hay IPs permitidas en la configuración. Saliendo."
		return fmt.Errorf("no hay IPs permitidas en la configuración")
	}

	address := fmt.Sprintf("%s:%d", "0.0.0.0", cfg.Port)
	addressTCP, err := net.ResolveTCPAddr("tcp4", address)

	if err != nil {
		// 💡 CAMBIO: Enviar al canal y devolver error.
		logChan <- fmt.Sprintf("Error fatal resolviendo dirección TCP: %v", err)
		return err
	}

	listener, err := net.ListenTCP("tcp4", addressTCP)
	if err != nil {
		// 💡 CAMBIO: Enviar al canal y devolver error.
		logChan <- fmt.Sprintf("Error fatal escuchando en TCP: %v", err)
		return err
	}

	// 💡 CAMBIO: Enviamos los logs de estado al canal.
	logChan <- fmt.Sprintf("Servidor escuchando en %s", addressTCP)
	logChan <- "Esperando conexiones..."

	for {
		socketServ, err := listener.Accept()
		if err != nil {
			// 💡 CAMBIO: Enviamos el error no-fatal al canal.
			logChan <- fmt.Sprintf("Error aceptando conexión: %v", err)
			continue
		}
		// 💡 CAMBIO: Pasamos el canal de logs a cada nueva conexión.
		go HandleConnection(socketServ, cfg, logChan)
	}

}

// 💡 CAMBIO: La función ahora acepta el canal de logs (logChan).
func HandleConnection(conn net.Conn, cfg Config, logChan chan<- string) {

	defer conn.Close()
	remoteIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	if !IsAllowed(remoteIP, cfg.Allowed_ips) {
		// 💡 CAMBIO: Log al canal.
		logChan <- fmt.Sprintf("Conexión rechazada desde IP no permitida: %s", remoteIP)
		return
	}

	// 💡 CAMBIO: Log al canal.
	logChan <- fmt.Sprintf("Cliente conectado desde: %s", remoteIP)

	reader := bufio.NewReader(conn)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 💡 CAMBIO: Log al canal.
		logChan <- fmt.Sprintf("Error crítico: No se pudo obtener el directorio home: %v", err)
		return
	}

	// 💡 BUGFIX CRÍTICO: NO cambiamos el directorio global (os.Chdir).
	// Solo rastreamos la ruta en una variable.
	currentDir := homeDir

	// Quitamos el 'os.Chdir(homeDir)' que era un bug.

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			// 💡 CAMBIO: Log al canal (informando qué cliente se fue).
			logChan <- fmt.Sprintf("Cliente %s desconectado (error leyendo comando): %v", remoteIP, err)
			return
		}

		// 💡 BUGFIX: Usamos TrimSpace para manejar \n y \r\n (Windows).
		command = strings.TrimSpace(command)
		if command == "" { // Ignorar entradas vacías
			continue
		}

		if command == "bye" {
			// 💡 CAMBIO: Log al canal.
			logChan <- fmt.Sprintf("Cliente %s cerró la sesión.", remoteIP)
			return
		}

		if strings.HasPrefix(command, "cd ") {
			target := strings.TrimSpace(strings.TrimPrefix(command, "cd "))
			if !filepath.IsAbs(target) {
				target = filepath.Join(currentDir, target)
			}

			// 💡 BUGFIX CRÍTICO: Validamos el directorio sin cambiarlo globalmente.
			// Verificamos si el directorio existe y es un directorio.
			if stat, err := os.Stat(target); err == nil && stat.IsDir() {
				currentDir = target // Actualizamos nuestro rastreador
				conn.Write([]byte("Directorio cambiado a: " + currentDir + "\n__END__\n"))
			} else {
				conn.Write([]byte("Error cambiando de directorio: " + err.Error() + "\n__END__\n"))
			}
			continue
		}

		if command == "repstats" {
			stats, err := monitor.GenerateReport() // Asumimos que esta versión no toma intervalo
			if err != nil {
				conn.Write([]byte("Error generando reporte: " + err.Error() + "\n__END__\n"))
				continue
			}
			conn.Write([]byte(string(stats) + "\n__END__\n"))
			continue
		}

		// 💡 CAMBIO: Log al canal (informando quién y qué).
		logChan <- fmt.Sprintf("[%s] Ejecutando comando: %s", remoteIP, command)

		// 💡 BUGFIX CRÍTICO: Pasamos el 'currentDir' a RunCommand.
		output, err := shell.RunCommand(command, currentDir)
		if err != nil {
			output += "\nError: " + err.Error()
		}

		conn.Write([]byte(output + "\n__END__\n"))
	}
}
