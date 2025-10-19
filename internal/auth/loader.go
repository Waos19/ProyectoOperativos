package auth

import (
	"fmt"
	"os"
	"strings"
)

func LoadUsers(path string) (map[string]string, error) {
	users := make(map[string]string)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de usuarios: %v", err)
	}

	lines := strings.SplitSeq(string(data), "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		user := strings.TrimSpace(parts[0])
		hash := strings.TrimSpace(parts[1])
		users[user] = hash
	}

	return users, nil
}
