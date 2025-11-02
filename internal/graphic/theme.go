package graphic

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// 💡 CAMBIO: Incrustamos la fuente JetBrains Mono directamente en el binario
var myFontData []byte

// myFontResource es el recurso estático de la fuente para Fyne
var myFontResource = &fyne.StaticResource{
	StaticName:    "JetBrainsMono-Regular.ttf",
	StaticContent: myFontData,
}

// MyTheme implementa la interfaz fyne.Theme para personalizar la apariencia de la aplicación
type MyTheme struct{}

// Verificación de implementación de la interfaz en tiempo de compilación
var _ fyne.Theme = (*MyTheme)(nil)

// Color devuelve el color para el nombre de tema especificado
// Personaliza los colores de texto deshabilitado y de primer plano a blanco
func (m *MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameDisabled, theme.ColorNameForeground:
		return color.White
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Icon devuelve el recurso del icono para el nombre especificado
// Utiliza los iconos predeterminados del tema
func (m *MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Font devuelve el recurso de la fuente personalizada
// Usa JetBrains Mono como fuente predeterminada
func (m *MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return myFontResource
}

// Size devuelve el tamaño para el nombre de tema especificado
// Aumenta el tamaño predeterminado en un 20%
func (m *MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name) * 1.2
}
