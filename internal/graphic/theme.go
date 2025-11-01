package graphic

import (
	"image/color" // Necesario para el tema

	_ "embed" // ¡Importante! Necesario para 'go:embed'

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// 'go:embed' le dice a Go que incluya este archivo en el binario.
// Asegúrate de que el nombre coincida EXACTAMENTE con tu archivo.
//
//go:embed JetBrainsMono-Regular.ttf
var myFontData []byte

// Creamos un "recurso estático" para Fyne
var myFontResource = &fyne.StaticResource{
	StaticName:    "JetBrainsMono-Regular.ttf",
	StaticContent: myFontData,
}

// myTheme es nuestra implementación de tema personalizado
type myTheme struct{}

// Esto asegura que nuestro tema implementa la interfaz correcta
var _ fyne.Theme = (*myTheme)(nil)

// --- Sobrescribimos los métodos del tema ---

// Color: Usamos el tema por defecto (DefaultTheme) para que
// se ajuste al modo claro/oscuro de tu PC.
func (m *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	// 💡 CAMBIO: Si la app pide el color "Deshabilitado" (Disabled)
	// (que Fyne usa para el texto en logView.Disable())
	if name == theme.ColorNameDisabled {
		// Forzamos que sea blanco puro
		return color.White
	}

	// 💡 CAMBIO: (Opcional) Si el color del texto principal
	// (como el de la caja de input) también se ve gris
	if name == theme.ColorNameForeground {
		return color.White
	}

	// Para todo lo demás (fondo, botones, etc.), usa el tema del sistema
	return theme.DefaultTheme().Color(name, variant)
}

// Icon: Usamos los iconos del tema por defecto
func (m *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Font: ¡Esta es la parte clave!
// Sobrescribimos la fuente para TODOS los estilos.
func (mM *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Devolvemos siempre nuestra fuente personalizada
	return myFontResource
}

// Size: Usamos los tamaños del tema por defecto
func (m *myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
