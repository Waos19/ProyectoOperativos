package graphic

import (
	"fmt"
	"proyoper/internal/auth"
	"proyoper/internal/client"
	"proyoper/internal/server"
	"strconv" // 💡 CAMBIO: Necesario para convertir el intervalo

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func StartClientUI() {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})
	myWindow := myApp.NewWindow("Cliente")
	myWindow.Resize(fyne.NewSize(400, 300))

	// --- 1. CREAMOS LOS WIDGETS DEL LOGIN ---
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("IP del Servidor (ej: 127.0.0.1)")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Puerto (ej: 9999)")

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Usuario")

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Contraseña")

	// 💡 CAMBIO: Añadimos la caja para el intervalo
	intervalEntry := widget.NewEntry()
	intervalEntry.SetPlaceHolder("Intervalo del monitor (segundos)")
	intervalEntry.SetText("5") // Valor por defecto

	statusLabel := widget.NewLabel("")

	var loginButton *widget.Button
	loginButton = widget.NewButton("Conectar y Autenticar", func() {
		statusLabel.SetText("Cargando configuración...")
		loginButton.Disable()

		// --- Obtenemos los datos ---
		ip := ipEntry.Text
		port := portEntry.Text
		username := userEntry.Text
		password := passEntry.Text

		// 💡 CAMBIO: Leemos y convertimos el intervalo
		interval, err := strconv.Atoi(intervalEntry.Text)
		if err != nil || interval <= 0 {
			// Si no es un número o es 0/negativo
			statusLabel.SetText("Error: El intervalo debe ser un número positivo.")
			loginButton.Enable()
			return
		}

		// --- Lógica de Autenticación ---
		cfg, err := server.LoadConfig("configs/server.conf")
		if err != nil {
			statusLabel.SetText("Error: No se pudo cargar config local.")
			loginButton.Enable()
			return
		}
		users, err := auth.LoadUsers(cfg.Users_file)
		if err != nil {
			statusLabel.SetText("Error: No se pudieron cargar usuarios.")
			loginButton.Enable()
			return
		}

		if !auth.VerifyLogin(username, password, users) {
			statusLabel.SetText("Error: Usuario o contraseña incorrectos.")
			loginButton.Enable()
			return
		}

		statusLabel.SetText("Autenticación exitosa. Conectando...")
		conn, err := client.Connect(ip, port)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Error de conexión: %v", err))
			loginButton.Enable()
			return
		}

		// 💡 CAMBIO: Pasamos el 'interval' a la función que construye la GUI
		mainLayout := BuildMainClientLayout(conn, interval)
		myWindow.SetContent(mainLayout)
		myWindow.Resize(fyne.NewSize(1024, 768))
	})

	// --- 2. LAYOUT DEL LOGIN ---
	loginForm := widget.NewForm(
		widget.NewFormItem("IP", ipEntry),
		widget.NewFormItem("Puerto", portEntry),
		widget.NewFormItem("Usuario", userEntry),
		widget.NewFormItem("Contraseña", passEntry),
		// 💡 CAMBIO: Añadimos el campo al formulario
		widget.NewFormItem("Intervalo (s)", intervalEntry),
	)

	loginLayout := container.New(
		layout.NewVBoxLayout(),
		loginForm,
		container.NewCenter(statusLabel),
		container.NewCenter(loginButton),
	)

	myWindow.SetContent(loginLayout)
	myWindow.ShowAndRun()
}
