package server

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port         int
	Allowed_ips  []string
	Max_attempts int
	Users_file   string
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
		case "Port":
			cfg.Port, _ = strconv.Atoi(value)
		case "Allowed_ips":
			raw := strings.Split(value, ",")
			for _, r := range raw {
				ip := strings.TrimSpace(r)
				if ip != "" {
					cfg.Allowed_ips = append(cfg.Allowed_ips, ip)
				}
			}
		case "Max_attempts":
			cfg.Max_attempts, _ = strconv.Atoi(value)
		case "Users_file":
			cfg.Users_file = value

		}
	}

	//fmt.Printf("Config cargado:\n%+v\n", cfg)

	return cfg, err
}
