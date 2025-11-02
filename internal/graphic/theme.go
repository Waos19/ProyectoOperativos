package graphic

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// 💡 CAMBIO: Eliminamos 'embed', 'myFontData' y 'myFontResource'
// ya que no usaremos una fuente personalizada.

// MyTheme implementa la interfaz fyne.Theme para personalizar la apariencia
type MyTheme struct{}

// Verificación de implementación de la interfaz
var _ fyne.Theme = (*MyTheme)(nil)

// Color devuelve el color para el nombre de tema especificado
// Mantenemos esto para que el texto sea blanco
func (m *MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	if name == theme.ColorNameDisabled {
		return color.White
	}
	if name == theme.ColorNameForeground {
		return color.White
	}

	return theme.DefaultTheme().Color(name, variant)
}

// Icon devuelve el recurso del icono para el nombre especificado
// (Sin cambios)
func (m *MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// 💡 CAMBIO: Font ahora devuelve la fuente estándar del sistema
func (m *MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Devuelve la fuente por defecto del tema (la estándar de Fyne)
	return theme.DefaultTheme().Font(style)
}

// Size devuelve el tamaño para el nombre de tema especificado
// Mantenemos esto para que la fuente sea un 20% más grande
func (m *MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name) * 1.2
}
