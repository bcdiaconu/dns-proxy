package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadConfig(path string) map[string]string {
	cfg := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			cfg[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	return cfg
}
