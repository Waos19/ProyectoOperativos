package graphic

import (
	"fmt"
	"net"
	"proyoper/internal/client"
	"proyoper/internal/monitor"
	"time" // Necesitamos 'time'

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// 💡 CAMBIO: La firma ahora acepta un 'interval' de tipo int
func BuildMainClientLayout(conn net.Conn, interval int) fyne.CanvasObject {

	// 💡 CAMBIO: Convertimos el 'int' a 'time.Duration'
	// (Ej: 5 -> 5 * time.Second)
	reportInterval := time.Duration(interval) * time.Second

	outputChan := make(chan string)
	errorChan := make(chan string)

	// --- [GOROUTINE 1] Lector de Red (sin cambios) ---
	go client.StartClientReader(conn, outputChan, errorChan)

	// --- Creación de Widgets (Las 3 Cajas) ---
	commandBox := widget.NewEntry()
	commandBox.SetPlaceHolder("Escribe un comando y presiona Enter...")

	resultBinding := binding.NewString()
	resultBox := widget.NewMultiLineEntry()
	resultBox.Bind(resultBinding)
	resultBox.Disable()

	reportBinding := binding.NewString()
	reportBox := widget.NewMultiLineEntry()
	reportBox.Bind(reportBinding)
	reportBox.Disable()
	reportBinding.Set(fmt.Sprintf("Iniciando monitor de sistema (actualizando cada %d segundos)...", interval))

	// --- Lógica de envío de Comandos (sin cambios) ---
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

	// --- [GOROUTINE 2] Actualizador de GUI (sin cambios) ---
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

	// 💡 CAMBIO: [GOROUTINE 3] Monitor del Sistema Local
	go func() {
		for {
			reportString, err := monitor.GenerateReport()
			if err != nil {
				reportBinding.Set(fmt.Sprintf("Error generando reporte:\n%v", err))
			} else {
				reportBinding.Set(reportString)
			}

			// 💡 CAMBIO: Usamos la variable 'reportInterval' que recibimos
			time.Sleep(reportInterval)
		}
	}()

	// --- Creación del Layout (Splits - sin cambios) ---
	topSplit := container.NewHSplit(
		container.NewScroll(commandBox),
		container.NewScroll(resultBox),
	)
	topSplit.SetOffset(0.4)

	mainSplit := container.NewVSplit(
		topSplit,
		container.NewScroll(reportBox),
	)
	mainSplit.SetOffset(0.8)

	return mainSplit
}
