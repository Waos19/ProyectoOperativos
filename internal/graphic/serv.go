package graphic

import (
	"fmt"
	"log"
	"proyoper/internal/server"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// StartServerUI inicializa y configura la interfaz gráfica del servidor
// Recibe una ventana principal donde se mostrarán los logs del servidor
func StartServerUI(myWindow fyne.Window) {
	// Configuración inicial de la ventana
	configureServerWindow(myWindow)

	// Inicialización de componentes UI
	logBinding := binding.NewString()
	logView := createLogView(logBinding)
	scrollContainer := container.NewScroll(logView)
	myWindow.SetContent(scrollContainer)

	// Canal para los logs del servidor
	logChan := make(chan string, 100)

	// Iniciar el manejador de logs
	go handleServerLogs(logBinding, logChan)

	// Cargar configuración e iniciar servidor
	startServerWithConfig(logChan)
}

// configureServerWindow configura los parámetros básicos de la ventana del servidor
func configureServerWindow(window fyne.Window) {
	window.SetTitle("Panel de Logs del Servidor")
	window.Resize(fyne.NewSize(800, 600))
	window.SetFixedSize(false)
	window.SetPadded(true)
}

// createLogView crea y configura el widget para mostrar los logs
func createLogView(binding binding.String) *widget.Entry {
	logView := widget.NewMultiLineEntry()
	logView.Bind(binding)
	logView.Disable()
	return logView
}

// handleServerLogs maneja la actualización de logs en la interfaz gráfica
func handleServerLogs(logBinding binding.String, logChan chan string) {
	for logLine := range logChan {
		currentLog, _ := logBinding.Get()
		logBinding.Set(currentLog + logLine + "\n")
	}
}

// startServerWithConfig carga la configuración e inicia el servidor
func startServerWithConfig(logChan chan string) {
	// Cargar configuración
	cfg, err := server.LoadConfig("configs/server.conf")
	if err != nil {
		errorMsg := fmt.Sprintf("ERROR FATAL: No se pudo cargar la configuración: %v", err)
		logChan <- errorMsg
		log.Println(errorMsg)
		return
	}

	// Iniciar servidor en una goroutine
	go func() {
		if err := server.StartServer(cfg, logChan); err != nil {
			logChan <- fmt.Sprintf("El servidor se ha detenido con un error: %v", err)
		}
	}()
}
