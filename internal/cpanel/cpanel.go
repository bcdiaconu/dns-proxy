package cpanel

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type CPanelConfig struct {
	URL    string
	User   string
	APIKey string
}

func NewCPanelConfig(cfg map[string]string) (*CPanelConfig, error) {
	url := cfg["cpanel_url"]
	user := cfg["cpanel_user"]
	apikey := cfg["cpanel_apikey"]
	if url == "" || user == "" || apikey == "" {
		return nil, errors.New("config incomplete: missing url, user or apikey")
	}
	return &CPanelConfig{URL: url, User: user, APIKey: apikey}, nil
}

func (c *CPanelConfig) CreateTxtRecord(domain, key, value string) error {
	data := url.Values{}
	data.Set("cpanel_jsonapi_user", c.User)
	data.Set("cpanel_jsonapi_apiversion", "2")
	data.Set("cpanel_jsonapi_module", "ZoneEdit")
	data.Set("cpanel_jsonapi_func", "add_zone_record")
	data.Set("domain", domain)
	data.Set("name", key)
	data.Set("type", "TXT")
	data.Set("txtdata", value)
	data.Set("ttl", "300")

	fullURL := fmt.Sprintf("%s/json-api/cpanel", c.URL)
	req, err := http.NewRequest("POST", fullURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("cpanel %s:%s", c.User, c.APIKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *CPanelConfig) DeleteTxtRecord(domain, key, value string) error {
	data := url.Values{}
	data.Set("cpanel_jsonapi_user", c.User)
	data.Set("cpanel_jsonapi_apiversion", "2")
	data.Set("cpanel_jsonapi_module", "ZoneEdit")
	data.Set("cpanel_jsonapi_func", "remove_zone_record")
	data.Set("domain", domain)
	data.Set("name", key)
	data.Set("type", "TXT")
	data.Set("txtdata", value)

	fullURL := fmt.Sprintf("%s/json-api/cpanel", c.URL)
	req, err := http.NewRequest("POST", fullURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("cpanel %s:%s", c.User, c.APIKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
