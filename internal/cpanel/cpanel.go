package cpanel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	// Extract the actual zone and record name
	zone, recordName := extractZoneAndName(domain)
	if recordName != "" {
		// If we have a subdomain, prepend the key to the record name
		recordName = key + "." + recordName
	} else {
		// If no subdomain, just use the key
		recordName = key
	}

	fmt.Printf("DEBUG: Creating TXT record - zone='%s', recordName='%s', value='%s'\n", zone, recordName, value)

	data := url.Values{}
	data.Set("cpanel_jsonapi_user", c.User)
	data.Set("cpanel_jsonapi_apiversion", "2")
	data.Set("cpanel_jsonapi_module", "ZoneEdit")
	data.Set("cpanel_jsonapi_func", "add_zone_record")
	data.Set("domain", zone)     // Use the extracted zone
	data.Set("name", recordName) // Use the extracted record name
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
	// Extract the actual zone and record name
	zone, recordName := extractZoneAndName(domain)
	if recordName != "" {
		// If we have a subdomain, prepend the key to the record name
		recordName = key + "." + recordName
	} else {
		// If no subdomain, just use the key
		recordName = key
	}

	fmt.Printf("DEBUG: Using zone='%s', recordName='%s'\n", zone, recordName)

	// 1. Fetch all zone records using cPanel API v2
	fetchData := url.Values{}
	fetchData.Set("cpanel_jsonapi_user", c.User)
	fetchData.Set("cpanel_jsonapi_apiversion", "2")
	fetchData.Set("cpanel_jsonapi_module", "ZoneEdit")
	fetchData.Set("cpanel_jsonapi_func", "fetchzone")
	fetchData.Set("domain", zone)    // Use the extracted zone
	fetchData.Set("customonly", "0") // Return all records, not just non-essential ones

	fullURL := fmt.Sprintf("%s/json-api/cpanel", c.URL)
	req, err := http.NewRequest("POST", fullURL, bytes.NewBufferString(fetchData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create fetch request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("cpanel %s:%s", c.User, c.APIKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, string(body))
	}

	// Debug: log the fetch response
	fmt.Printf("DEBUG: fetchzone response: %s\n", string(body))

	// 2. Parse cPanel API v2 response and find the record
	var fetchResp struct {
		CPanelResult struct {
			Data []struct {
				Record []struct {
					Line    int    `json:"Line"` // Capital L as per API docs
					Name    string `json:"name"`
					Type    string `json:"type"`
					TxtData string `json:"txtdata"`
				} `json:"record"`
			} `json:"data"`
		} `json:"cpanelresult"`
	}
	if err := json.Unmarshal(body, &fetchResp); err != nil {
		return fmt.Errorf("failed to parse fetchzone response: %w", err)
	}

	// Debug: log what we're searching for
	fmt.Printf("DEBUG: Looking for TXT record with name='%s' and txtdata='%s'\n", recordName+"."+zone+".", value)

	var foundID *int
	for _, data := range fetchResp.CPanelResult.Data {
		for _, rec := range data.Record {
			fmt.Printf("DEBUG: Found record - Line: %d, Name: '%s', Type: '%s', TxtData: '%s'\n",
				rec.Line, rec.Name, rec.Type, rec.TxtData)

			// Check if this is our TXT record
			if rec.Type == "TXT" && rec.Name == recordName+"."+zone+"." && rec.TxtData == value {
				id := rec.Line
				foundID = &id
				break
			}
		}
		if foundID != nil {
			break
		}
	}
	if foundID == nil {
		return fmt.Errorf("TXT record not found for deletion")
	}

	fmt.Printf("DEBUG: Found record to delete with line: %d\n", *foundID)

	// 3. Remove the record by line using cPanel API v2
	delData := url.Values{}
	delData.Set("cpanel_jsonapi_user", c.User)
	delData.Set("cpanel_jsonapi_apiversion", "2")
	delData.Set("cpanel_jsonapi_module", "ZoneEdit")
	delData.Set("cpanel_jsonapi_func", "remove_zone_record")
	delData.Set("domain", zone) // Use the extracted zone
	delData.Set("line", fmt.Sprintf("%d", *foundID))

	delReq, err := http.NewRequest("POST", fullURL, bytes.NewBufferString(delData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	delReq.Header.Set("Authorization", fmt.Sprintf("cpanel %s:%s", c.User, c.APIKey))
	delReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	delResp, err := client.Do(delReq)
	if err != nil {
		return fmt.Errorf("delete request failed: %w", err)
	}
	defer delResp.Body.Close()

	delBody, _ := io.ReadAll(delResp.Body)
	if delResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %d: %s", delResp.StatusCode, string(delBody))
	}

	// Debug: log delete response
	fmt.Printf("DEBUG: remove_zone_record response: %s\n", string(delBody))

	// Parse and validate the delete response
	var delResult struct {
		CPanelResult struct {
			Data []struct {
				Result struct {
					NewSerial interface{} `json:"newserial"` // Can be string or int
					StatusMsg string      `json:"statusmsg"`
					Status    int         `json:"status"`
				} `json:"result"`
			} `json:"data"`
			Event struct {
				Result int `json:"result"`
			} `json:"event"`
		} `json:"cpanelresult"`
	}

	if err := json.Unmarshal(delBody, &delResult); err != nil {
		return fmt.Errorf("failed to parse remove_zone_record response: %w", err)
	}

	// Check if the operation was successful
	if delResult.CPanelResult.Event.Result != 1 {
		return fmt.Errorf("remove_zone_record failed: event result was %d", delResult.CPanelResult.Event.Result)
	}

	if len(delResult.CPanelResult.Data) > 0 && delResult.CPanelResult.Data[0].Result.Status != 1 {
		return fmt.Errorf("remove_zone_record failed: %s", delResult.CPanelResult.Data[0].Result.StatusMsg)
	}

	fmt.Printf("DEBUG: Record successfully deleted. New serial: %v\n",
		delResult.CPanelResult.Data[0].Result.NewSerial)

	return nil
}

func (c *CPanelConfig) EditTxtRecord(domain, key, oldValue, newValue string) error {
	// Extract the actual zone and record name
	zone, recordName := extractZoneAndName(domain)
	if recordName != "" {
		// If we have a subdomain, prepend the key to the record name
		recordName = key + "." + recordName
	} else {
		// If no subdomain, just use the key
		recordName = key
	}

	fmt.Printf("DEBUG: Using zone='%s', recordName='%s'\n", zone, recordName)

	// 1. Fetch all zone records using cPanel API v2 to find the record to edit
	fetchData := url.Values{}
	fetchData.Set("cpanel_jsonapi_user", c.User)
	fetchData.Set("cpanel_jsonapi_apiversion", "2")
	fetchData.Set("cpanel_jsonapi_module", "ZoneEdit")
	fetchData.Set("cpanel_jsonapi_func", "fetchzone")
	fetchData.Set("domain", zone) // Use the extracted zone
	fetchData.Set("customonly", "0")

	fullURL := fmt.Sprintf("%s/json-api/cpanel", c.URL)
	req, err := http.NewRequest("POST", fullURL, bytes.NewBufferString(fetchData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create fetch request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("cpanel %s:%s", c.User, c.APIKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, string(body))
	}

	// Debug: log the fetch response
	fmt.Printf("DEBUG: fetchzone response: %s\n", string(body))

	// 2. Parse cPanel API v2 response and find the record
	var fetchResp struct {
		CPanelResult struct {
			Data []struct {
				Record []struct {
					Line    int    `json:"Line"` // Capital L as per API docs
					Name    string `json:"name"`
					Type    string `json:"type"`
					TxtData string `json:"txtdata"`
				} `json:"record"`
			} `json:"data"`
		} `json:"cpanelresult"`
	}
	if err := json.Unmarshal(body, &fetchResp); err != nil {
		return fmt.Errorf("failed to parse fetchzone response: %w", err)
	}

	// Debug: log what we're searching for
	fmt.Printf("DEBUG: Looking for TXT record with name='%s' and txtdata='%s'\n", recordName+"."+zone+".", oldValue)

	var foundLine *int
	for _, data := range fetchResp.CPanelResult.Data {
		for _, rec := range data.Record {
			fmt.Printf("DEBUG: Found record - Line: %d, Name: '%s', Type: '%s', TxtData: '%s'\n",
				rec.Line, rec.Name, rec.Type, rec.TxtData)

			// Check if this is our TXT record
			if rec.Type == "TXT" && rec.Name == recordName+"."+zone+"." && rec.TxtData == oldValue {
				line := rec.Line
				foundLine = &line
				break
			}
		}
		if foundLine != nil {
			break
		}
	}
	if foundLine == nil {
		return fmt.Errorf("TXT record not found for editing")
	}

	fmt.Printf("DEBUG: Found record to edit at line: %d\n", *foundLine)

	// 3. Edit the record using cPanel API v2 edit_zone_record
	editData := url.Values{}
	editData.Set("cpanel_jsonapi_user", c.User)
	editData.Set("cpanel_jsonapi_apiversion", "2")
	editData.Set("cpanel_jsonapi_module", "ZoneEdit")
	editData.Set("cpanel_jsonapi_func", "edit_zone_record")
	editData.Set("Line", fmt.Sprintf("%d", *foundLine)) // Capital L as per API docs
	editData.Set("domain", zone)                        // Use the extracted zone
	editData.Set("name", recordName)                    // Use the extracted record name
	editData.Set("type", "TXT")
	editData.Set("txtdata", newValue)
	editData.Set("ttl", "300")
	editData.Set("class", "IN")

	editReq, err := http.NewRequest("POST", fullURL, bytes.NewBufferString(editData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create edit request: %w", err)
	}
	editReq.Header.Set("Authorization", fmt.Sprintf("cpanel %s:%s", c.User, c.APIKey))
	editReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	editResp, err := client.Do(editReq)
	if err != nil {
		return fmt.Errorf("edit request failed: %w", err)
	}
	defer editResp.Body.Close()

	editBody, _ := io.ReadAll(editResp.Body)
	if editResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %d: %s", editResp.StatusCode, string(editBody))
	}

	// Debug: log edit response
	fmt.Printf("DEBUG: edit_zone_record response: %s\n", string(editBody))

	return nil
}

// extractZoneAndName extracts the zone and record name from a full domain
// For example: "_acme-challenge.haos.iveronsoft.ro" -> zone: "iveronsoft.ro", name: "_acme-challenge.haos"
func extractZoneAndName(fullDomain string) (zone, name string) {
	parts := strings.Split(fullDomain, ".")
	if len(parts) < 3 {
		// If less than 3 parts, assume it's already a zone
		return fullDomain, ""
	}

	// Assume the zone is the last two parts (domain.tld)
	zone = strings.Join(parts[len(parts)-2:], ".")

	// The name is everything before the zone
	if len(parts) > 2 {
		name = strings.Join(parts[:len(parts)-2], ".")
	}

	return zone, name
}
