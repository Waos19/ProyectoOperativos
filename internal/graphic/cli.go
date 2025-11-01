package graphic

import (
	"fmt"
	"log"
	"proyoper/internal/auth"
	"proyoper/internal/client"
	"proyoper/internal/server"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func StartClientUI() {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})
	myWindow := myApp.NewWindow("Cliente")
	myWindow.Resize(fyne.NewSize(400, 300))

	cfg, err := server.LoadConfig("configs/server.conf")
	if err != nil {
		log.Fatalf("Error crítico: No se pudo cargar config local: %v", err)
	}
	users, err := auth.LoadUsers(cfg.Users_file)
	if err != nil {
		log.Fatalf("Error crítico: No se pudieron cargar usuarios: %v", err)
	}

	var attempts int = 0
	maxAttempts := cfg.Max_attempts

	// --- 1. CREAMOS LOS WIDGETS DEL LOGIN ---
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("IP del Servidor (ej: 127.0.0.1)")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Puerto (ej: 9999)")

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Usuario")

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Contraseña")

	intervalEntry := widget.NewEntry()
	intervalEntry.SetPlaceHolder("Intervalo del monitor (segundos)")
	intervalEntry.SetText("5")

	statusBinding := binding.NewString()
	statusLabel := widget.NewLabelWithData(statusBinding)

	buttonDisabled := binding.NewBool()

	loginButton := widget.NewButton("Conectar y Autenticar", func() {
		// --- Inicia el proceso ---
		statusBinding.Set("Verificando...")
		buttonDisabled.Set(true) // Deshabilita el botón

		// --- Obtenemos los datos ---
		ip := ipEntry.Text
		port := portEntry.Text
		username := userEntry.Text
		password := passEntry.Text

		interval, err := strconv.Atoi(intervalEntry.Text)
		if err != nil || interval <= 0 {
			statusBinding.Set("Error: El intervalo debe ser un número positivo.")
			buttonDisabled.Set(false) // Rehabilita el botón
			return
		}

		if !auth.VerifyLogin(username, password, users) {
			// --- INTENTO FALLIDO ---
			attempts++

			if attempts >= maxAttempts {
				// --- BLOQUEADO ---
				statusBinding.Set(fmt.Sprintf("Demasiados intentos fallidos (%d/%d). Aplicación bloqueada.", attempts, maxAttempts))
				// No re-habilitamos el botón (se queda deshabilitado)
				return
			}

			// --- ESPERAR ---
			wait := time.Duration(attempts*2) * time.Second
			statusBinding.Set(fmt.Sprintf("Intento %d/%d fallido. Espere %v...", attempts, maxAttempts, wait))

			go func() {
				time.Sleep(wait)
				statusBinding.Set("Intente de nuevo.")
				buttonDisabled.Set(false) // Rehabilita el botón
			}()

			return
		}

		// --- INTENTO EXITOSO ---
		statusBinding.Set("Autenticación exitosa. Conectando...")
		conn, err := client.Connect(ip, port)
		if err != nil {
			statusBinding.Set(fmt.Sprintf("Error de conexión: %v", err))
			buttonDisabled.Set(false)
			return
		}

		// --- CONEXIÓN EXITOSA ---
		mainLayout := BuildMainClientLayout(conn, interval)
		myWindow.SetContent(mainLayout)
		myWindow.Resize(fyne.NewSize(1024, 768))
	})

	// 💡 CAMBIO: Esta es la forma correcta de enlazar el estado "deshabilitado"
	// Añadimos un "listener" al binding.
	buttonDisabled.AddListener(binding.NewDataListener(func() {
		// Esta función se ejecuta CADA VEZ que buttonDisabled.Set() es llamado.
		if disabled, err := buttonDisabled.Get(); err == nil {
			if disabled {
				loginButton.Disable()
			} else {
				loginButton.Enable()
			}
		}
	}))

	// --- 2. LAYOUT DEL LOGIN ---
	loginForm := widget.NewForm(
		widget.NewFormItem("IP", ipEntry),
		widget.NewFormItem("Puerto", portEntry),
		widget.NewFormItem("Usuario", userEntry),
		widget.NewFormItem("Contraseña", passEntry),
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
