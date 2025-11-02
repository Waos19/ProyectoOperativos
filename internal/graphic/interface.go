package graphic

import (
	"fmt"
	"net"
	"proyoper/internal/client"
	"proyoper/internal/monitor"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// BuildMainClientLayout construye la interfaz principal del cliente
// Recibe una conexión de red y un intervalo para el monitor del sistema
func BuildMainClientLayout(conn net.Conn, interval int) fyne.CanvasObject {
	// Configuración inicial
	reportInterval := time.Duration(interval) * time.Second
	outputChan := make(chan string)
	errorChan := make(chan string)

	// Iniciar el lector de red en una goroutine
	go client.StartClientReader(conn, outputChan, errorChan)

	// Inicialización de widgets
	// 1. Caja de comandos
	commandBox := widget.NewEntry()
	commandBox.SetPlaceHolder("Escribe un comando y presiona Enter...")

	// 2. Caja de resultados
	resultBinding := binding.NewString()
	resultBox := widget.NewMultiLineEntry()
	resultBox.Bind(resultBinding)
	resultBox.Disable()

	// 3. Caja de reportes
	reportBinding := binding.NewString()
	reportBox := widget.NewMultiLineEntry()
	reportBox.Bind(reportBinding)
	reportBox.Disable()
	reportBinding.Set(fmt.Sprintf("Iniciando monitor de sistema (actualizando cada %d segundos)...",
		interval))

	// Configuración del manejador de comandos
	commandBox.OnSubmitted = func(cmd string) {
		if cmd == "" {
			return
		}

		_, err := conn.Write([]byte(cmd + "\n"))
		if err != nil {
			errorChan <- fmt.Sprintf("Error enviando comando: %v", err)
		}

		if cmd == "bye" {
			fyne.CurrentApp().Quit()
		}

		commandBox.SetText("")
	}

	// Goroutine para actualizar la interfaz gráfica
	go func() {
		for {
			select {
			case newOutput := <-outputChan:
				resultBinding.Set(newOutput)
			case errMsg := <-errorChan:
				resultBinding.Set(errMsg)
			}
		}
	}()

	// Goroutine para el monitor del sistema
	go func() {
		for {
			reportString, err := monitor.GenerateReport()
			if err != nil {
				reportBinding.Set(fmt.Sprintf("Error generando reporte:\n%v", err))
			} else {
				reportBinding.Set(reportString)
			}
			time.Sleep(reportInterval)
		}
	}()

	// Construcción del layout
	// 1. Split superior (comandos y resultados)
	topSplit := container.NewHSplit(
		container.NewScroll(commandBox),
		container.NewScroll(resultBox),
	)
	topSplit.SetOffset(0.4)

	// 2. Split principal (superior y reportes)
	mainSplit := container.NewVSplit(
		topSplit,
		container.NewScroll(reportBox),
	)
	mainSplit.SetOffset(0.8)

	return mainSplit
}
