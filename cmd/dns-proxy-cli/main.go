package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"dns-proxy/internal/cpanel"
)

func loadCPanelConfig(path string) map[string]string {
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
	return cfg
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dns-proxy-cli set-txt --domain <domain> --key <key> --value <value>")
		os.Exit(1)
	}

	subcmd := os.Args[1]
	if subcmd != "set-txt" {
		fmt.Println("Unknown command:", subcmd)
		os.Exit(1)
	}

	setTxtCmd := flag.NewFlagSet("set-txt", flag.ExitOnError)
	domain := setTxtCmd.String("domain", "", "Domain name")
	key := setTxtCmd.String("key", "", "TXT record key")
	value := setTxtCmd.String("value", "", "TXT record value")

	setTxtCmd.Parse(os.Args[2:])

	if *domain == "" || *key == "" || *value == "" {
		fmt.Println("All arguments --domain, --key, and --value are required.")
		os.Exit(1)
	}

	cfg := loadCPanelConfig("/etc/dns-proxy-cli.conf")
	cpCfg, err := cpanel.NewCPanelConfig(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = cpCfg.CreateTxtRecord(*domain, *key, *value)
	if err != nil {
		log.Fatalf("Failed to set TXT record: %v", err)
	}

	fmt.Println("TXT record set successfully.")
}
