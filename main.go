package main

import (
	"log"
	"net/http"

	"dns-proxy/internal/api"
	"dns-proxy/internal/config"
	"dns-proxy/internal/cpanel"
)

func main() {
	cfg := config.LoadConfig("/etc/dns-proxy.conf")
	apiKey := cfg["API_KEY"]

	cpCfg, err := cpanel.NewCPanelConfig(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	http.HandleFunc("/set_txt", api.SetTxtHandler(apiKey, cpCfg))

	log.Println("dns-proxy listening on :5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
