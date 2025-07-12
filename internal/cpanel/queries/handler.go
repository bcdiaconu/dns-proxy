package queries

import (
	"fmt"
	"strings"
)

// CPanelQueryHandler implements QueryHandler interface
type CPanelQueryHandler struct{}

// NewCPanelQueryHandler creates a new query handler
func NewCPanelQueryHandler() QueryHandler {
	return &CPanelQueryHandler{}
}

// HandleList handles listing TXT records
func (h *CPanelQueryHandler) HandleList(query *ListTxtRecordsQuery) ([]TxtRecord, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	// Extract the actual zone
	zone, recordPrefix := extractZoneAndName(query.Request.Domain)
	
	fmt.Printf("DEBUG: Listing TXT records for zone='%s', recordPrefix='%s', keyFilter='%s'\n", 
		zone, recordPrefix, query.Request.KeyFilter)

	// Implementation will be moved from cpanel.go
	return h.listTxtRecordsAPI(zone, recordPrefix, query.Request.KeyFilter)
}

// Execute implements TxtRecordQuery interface
func (q *ListTxtRecordsQuery) Execute() (interface{}, error) {
	handler := NewCPanelQueryHandler()
	return handler.HandleList(q)
}

// Validate validates the list query
func (q *ListTxtRecordsQuery) Validate() error {
	if q.Request.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	return nil
}

// Private helper methods
func (h *CPanelQueryHandler) listTxtRecordsAPI(zone, recordPrefix, keyFilter string) ([]TxtRecord, error) {
	// TODO: Move implementation from cpanel.go ListTxtRecords method
	return nil, fmt.Errorf("not implemented yet")
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
