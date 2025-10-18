package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetInput() (string, error) {
	fmt.Print("-> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func RunCommand(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", fmt.Errorf("comando vacío")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
