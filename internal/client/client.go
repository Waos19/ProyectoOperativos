package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// 💡 NUEVA FUNCIÓN: Conectar
// Esta función solo establece la conexión y la devuelve.
// La GUI la llamará al principio.
func Connect(ipServer string, port string) (net.Conn, error) {

	addressServ := ipServer + ":" + port

	dirTCP, err := net.ResolveTCPAddr("tcp4", addressServ)
	if err != nil {
		// Devolvemos el error para que la GUI pueda mostrarlo
		return nil, fmt.Errorf("error resolviendo dirección: %v", err)
	}

	socketC, err := net.DialTCP("tcp4", nil, dirTCP)
	if err != nil {
		// Devolvemos el error
		return nil, fmt.Errorf("error conectando al servidor: %v", err)
	}

	// Éxito, devolvemos la conexión
	return socketC, nil
}

// 💡 NUEVA FUNCIÓN: Lector de Red
// Esta función debe ejecutarse en una GOROUTINE.
// Escucha constantemente la red y envía los resultados a los canales.
func StartClientReader(conn net.Conn, outputChan chan<- string, errorChan chan<- string) {
	// Aseguramos que la conexión se cierre si esta goroutine termina
	// (por ejemplo, si el servidor se desconecta).
	defer conn.Close()

	serverReader := bufio.NewReader(conn)

	for {
		// Usamos un strings.Builder para ensamblar la respuesta
		// que puede tener múltiples líneas.
		var response strings.Builder

		// Bucle para leer hasta el delimitador __END__
		for {
			line, err := serverReader.ReadString('\n')
			if err != nil {
				// Si hay un error (ej. servidor desconectado),
				// lo enviamos al canal de error y terminamos la goroutine.
				errorChan <- fmt.Sprintf("Servidor desconectado o error: %v", err)
				return
			}

			lineTrimmed := strings.TrimSpace(line)
			if lineTrimmed == "__END__" {
				// Terminamos de leer esta respuesta
				break
			}

			// Añadimos la línea (con su \n) a la respuesta
			response.WriteString(line)
		}

		// Enviamos la respuesta completa (ya ensamblada) al canal de salida.
		// La GUI estará escuchando este canal.
		outputChan <- response.String()
	}
}

/*
// --- FUNCIONES ANTERIORES (OBSOLETAS PARA LA GUI) ---
//
// ClientStart y RunInteractiveSession ya no se necesitan,
// porque la GUI manejará el "bucle" y la entrada/salida.
//
// func ClientStart(ipServer string, port string) {
// 	 ...
// 	 RunInteractiveSession(socketC)
// }
//
// func RunInteractiveSession(conn net.Conn) {
// 	 ...
// }
*/
