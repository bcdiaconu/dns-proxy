package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"dns-proxy/internal/commands"
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
		fmt.Println("Usage: dns-proxy-cli <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  set-txt --domain <domain> --key <key> --value <value>")
		fmt.Println("  delete-txt --domain <domain> --key <key> --value <value>")
		os.Exit(1)
	}

	subcmd := os.Args[1]

	// Create command factory and get command
	factory := commands.NewCommandFactory()
	cmd, err := factory.CreateCommand(subcmd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Parse arguments based on command
	args := parseCommandArgs(subcmd, os.Args[2:])

	// Validate arguments
	if err := cmd.ValidateArgs(args); err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Usage: %s\n", cmd.Usage())
		os.Exit(1)
	}

	// Load cPanel config
	cfg := loadCPanelConfig("/etc/dns-proxy-cli.conf")
	cpCfg, err := cpanel.NewCPanelConfig(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Execute command
	if err := cmd.Execute(cpCfg, args); err != nil {
		log.Fatalf("%v", err)
	}
}

func parseCommandArgs(subcmd string, args []string) map[string]string {
	var cmdFlags *flag.FlagSet

	switch subcmd {
	case "set-txt", "delete-txt":
		cmdFlags = flag.NewFlagSet(subcmd, flag.ExitOnError)
	default:
		return nil
	}

	domain := cmdFlags.String("domain", "", "Domain name")
	key := cmdFlags.String("key", "", "TXT record key")
	value := cmdFlags.String("value", "", "TXT record value")

	cmdFlags.Parse(args)

	return map[string]string{
		"domain": *domain,
		"key":    *key,
		"value":  *value,
	}
}
