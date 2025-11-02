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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func StartClientUI(myWindow fyne.Window) {
	// Configuración inicial de la ventana
	myWindow.SetTitle("Cliente - Login")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetFixedSize(false)
	myWindow.SetPadded(false)

	// Carga de configuración y usuarios
	cfg, err := server.LoadConfig("configs/server.conf")
	if err != nil {
		log.Fatalf("Error crítico: No se pudo cargar config local: %v", err)
	}
	users, err := auth.LoadUsers(cfg.Users_file)
	if err != nil {
		log.Fatalf("Error crítico: No se pudieron cargar usuarios: %v", err)
	}

	// Control de intentos de inicio de sesión
	var attempts int = 0
	maxAttempts := cfg.Max_attempts

	// Inicialización de widgets de login
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

	// Configuración del botón de login
	loginButton := widget.NewButton("Conectar y Autenticar", func() {
		statusBinding.Set("Verificando...")
		buttonDisabled.Set(true)

		// Obtención de datos del formulario
		ip := ipEntry.Text
		port := portEntry.Text
		username := userEntry.Text
		password := passEntry.Text
		interval, err := strconv.Atoi(intervalEntry.Text)

		// Validación del intervalo
		if err != nil || interval <= 0 {
			statusBinding.Set("Error: El intervalo debe ser un número positivo.")
			buttonDisabled.Set(false)
			return
		}

		// Verificación de login
		if !auth.VerifyLogin(username, password, users) {
			attempts++
			if attempts >= maxAttempts {
				statusBinding.Set(fmt.Sprintf("Demasiados intentos (%d/%d). La app se cerrará...",
					attempts, maxAttempts))
				go func() {
					time.Sleep(3 * time.Second)
					fyne.CurrentApp().Quit()
				}()
				return
			}

			wait := time.Duration(attempts*2) * time.Second
			statusBinding.Set(fmt.Sprintf("Intento %d/%d fallido. Espere %v...",
				attempts, maxAttempts, wait))
			go func() {
				time.Sleep(wait)
				statusBinding.Set("Intente de nuevo.")
				buttonDisabled.Set(false)
			}()
			return
		}

		// Intento de conexión
		statusBinding.Set("Autenticación exitosa. Conectando...")
		conn, err := client.Connect(ip, port)
		if err != nil {
			statusBinding.Set(fmt.Sprintf("Error de conexión: %v", err))
			buttonDisabled.Set(false)
			return
		}

		// Configuración de la ventana principal tras conexión exitosa
		myWindow.SetFixedSize(false)
		myWindow.SetTitle("Cliente - Conectado")
		myWindow.Resize(fyne.NewSize(800, 600))

		mainLayout := BuildMainClientLayout(conn, interval)
		myWindow.SetContent(mainLayout)
	})

	// Configuración del listener para el estado del botón
	buttonDisabled.AddListener(binding.NewDataListener(func() {
		if disabled, err := buttonDisabled.Get(); err == nil {
			if disabled {
				loginButton.Disable()
			} else {
				loginButton.Enable()
			}
		}
	}))

	// Construcción del formulario de login
	loginForm := widget.NewForm(
		widget.NewFormItem("IP", ipEntry),
		widget.NewFormItem("Puerto", portEntry),
		widget.NewFormItem("Usuario", userEntry),
		widget.NewFormItem("Contraseña", passEntry),
		widget.NewFormItem("Intervalo (s)", intervalEntry),
	)

	// Construcción del layout principal
	loginLayout := container.New(
		layout.NewVBoxLayout(),
		loginForm,
		container.NewCenter(statusLabel),
		container.NewCenter(loginButton),
	)

	myWindow.SetContent(loginLayout)
}
