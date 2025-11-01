// 💡 CAMBIO: El paquete ahora es 'graphic'
package graphic

import (
	"fmt"
	"log"
	"proyoper/internal/server" // Necesita el paquete 'server'

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// 💡 CAMBIO: Nombre en mayúscula para exportar
func StartServerUI() {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})
	myWindow := myApp.NewWindow("Panel de Logs del Servidor")
	myWindow.Resize(fyne.NewSize(800, 600))

	logBinding := binding.NewString()

	logView := widget.NewMultiLineEntry()
	logView.Bind(logBinding)
	logView.Disable()

	scrollContainer := container.NewScroll(logView)
	myWindow.SetContent(scrollContainer)

	logChan := make(chan string, 100)

	// [GOROUTINE 1] Escuchador de GUI
	go func() {
		for logLine := range logChan {
			currentLog, _ := logBinding.Get()
			// 💡 CORRECCIÓN MENOR: Asegurar que SetString es thread-safe
			// Fyne prefiere que las actualizaciones de binding se hagan así
			logBinding.Set(currentLog + logLine + "\n")
		}
	}()

	cfg, err := server.LoadConfig("configs/server.conf")
	if err != nil {
		errorMsg := fmt.Sprintf("ERROR FATAL: No se pudo cargar la configuración: %v", err)
		logChan <- errorMsg
		log.Println(errorMsg)
		myWindow.ShowAndRun()
		return
	}

	// [GOROUTINE 2] Servidor
	go func() {
		err := server.StartServer(cfg, logChan)
		if err != nil {
			logChan <- fmt.Sprintf("El servidor se ha detenido con un error: %v", err)
		}
	}()

	myWindow.ShowAndRun()
}
