package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Esta función no se toca, es para la entrada local (si la usas).
func GetInput() (string, error) {
	fmt.Print("-> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// 💡 CAMBIO: La firma ahora acepta un argumento 'dir' (directorio).
func RunCommand(input string, dir string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", fmt.Errorf("comando vacío")
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	// 💡 CAMBIO: Esta es la línea clave.
	// Le decimos al comando que se ejecute en el directorio
	// específico de ESE cliente ('currentDir'), no en el
	// directorio global del servidor.
	cmd.Dir = dir

	// CombinedOutput captura tanto stdout como stderr.
	output, err := cmd.CombinedOutput()
	return string(output), err
}
