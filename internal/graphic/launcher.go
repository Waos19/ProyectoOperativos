package graphic

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// StartAppLauncher es el punto de entrada principal de la aplicación.
// Inicializa la interfaz gráfica y muestra la pantalla del lanzador.
func StartAppLauncher() {
	// Inicialización de la aplicación
	myApp := app.New()

	// Configuración del tema personalizado
	myApp.Settings().SetTheme(&MyTheme{})

	// Creación de la ventana principal
	myWindow := myApp.NewWindow("Lanzador")

	// Configuración de la pantalla inicial
	ShowLauncherScreen(myWindow)

	// Inicio de la aplicación
	myWindow.ShowAndRun()
}

// ShowLauncherScreen configura y muestra la pantalla principal del lanzador
func ShowLauncherScreen(myWindow fyne.Window) {
	// Configuración de la ventana
	myWindow.SetTitle("Lanzador")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetFixedSize(false)

	// Creación de botones
	serverButton := widget.NewButton(
		"Iniciar Servidor",
		func() {
			StartServerUI(myWindow)
		},
	)

	clientButton := widget.NewButton(
		"Iniciar Cliente",
		func() {
			StartClientUI(myWindow)
		},
	)

	// Construcción del layout
	content := container.NewVBox(
		layout.NewSpacer(), // Espacio superior
		serverButton,       // Botón del servidor
		clientButton,       // Botón del cliente
		layout.NewSpacer(), // Espacio inferior
	)

	// Aplicación del layout centrado
	myWindow.SetContent(
		container.NewCenter(content),
	)
}
