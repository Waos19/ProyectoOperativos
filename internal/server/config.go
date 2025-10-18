package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	port         int
	allowed_ips  []string
	max_attempts int
	users_file   string
}

func LoadConfig(path string) (Config, error) {
	//lee el archivo
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	//separa las lineas de texto del archivo
	lines := strings.Split(string(data), "\n")
	cfg := Config{}

	for _, line := range lines {
		if line == "" || !strings.Contains(line, "=") {
			continue
		}
		//Separa el texto con = y en 2 slices
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "port":
			cfg.port, _ = strconv.Atoi(value)
		case "allowed_ips":
			raw := strings.Split(value, ",")
			for _, r := range raw {
				ip := strings.TrimSpace(r)
				if ip != "" {
					cfg.allowed_ips = append(cfg.allowed_ips, ip)
				}
			}
		case "max_attempts":
			cfg.max_attempts, _ = strconv.Atoi(value)
		case "users_file":
			cfg.users_file = value

		}
	}

	fmt.Printf("Config cargado:\n%+v\n", cfg)

	return cfg, err
}
