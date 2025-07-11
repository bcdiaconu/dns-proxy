package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var config map[string]string

func main() {
	config = loadConfig("/etc/dns-proxy.conf")
	apiKey := config["API_KEY"]

	http.HandleFunc("/set_txt", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expected := "Bearer " + apiKey
		if authHeader != expected {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			Domain string `json:"domain"`
			Key    string `json:"key"`
			Value  string `json:"value"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.Domain == "" || req.Key == "" || req.Value == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = createCpanelTxtRecord(config, req.Domain, req.Key, req.Value)
		if err != nil {
			log.Println("cPanel error:", err)
			http.Error(w, "Failed to set TXT record", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TXT record set"))
	})

	log.Println("dns-proxy listening on :5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func loadConfig(path string) map[string]string {
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

func createCpanelTxtRecord(cfg map[string]string, domain, key, value string) error {
	// extract login
	apiURL := cfg["cpanel_url"]
	user := cfg["cpanel_user"]
	apiKey := cfg["cpanel_apikey"]

	if apiURL == "" || user == "" || apiKey == "" {
		return errors.New("config incomplete: missing url, user or apikey")
	}

	// prepare POST body
	data := url.Values{}
	data.Set("cpanel_jsonapi_user", user)
	data.Set("cpanel_jsonapi_apiversion", "2")
	data.Set("cpanel_jsonapi_module", "ZoneEdit")
	data.Set("cpanel_jsonapi_func", "add_zone_record")
	data.Set("domain", domain)
	data.Set("name", "_acme-challenge")
	data.Set("type", "TXT")
	data.Set("name", key)
	data.Set("txtdata", value)
	data.Set("ttl", "300")

	// create request
	fullURL := fmt.Sprintf("%s/json-api/cpanel", apiURL)
	req, err := http.NewRequest("POST", fullURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// set headers
	req.Header.Set("Authorization", fmt.Sprintf("cpanel %s:%s", user, apiKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// run request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// get answer
	body, _ := io.ReadAll(resp.Body)

	// verify response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, string(body))
	}

	// not implemented but json response can be parse for a better response
	return nil
}

