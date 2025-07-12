package commands

import (
	"fmt"
	"strings"
)

// CPanelCommandHandler implements CommandHandler interface
type CPanelCommandHandler struct{}

// NewCPanelCommandHandler creates a new command handler
func NewCPanelCommandHandler() CommandHandler {
	return &CPanelCommandHandler{}
}

// HandleCreate handles creating a TXT record
func (h *CPanelCommandHandler) HandleCreate(cmd *CreateTxtRecordCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Extract the actual zone and record name
	zone, recordName := extractZoneAndName(cmd.Request.Domain)
	if recordName != "" {
		// If we have a subdomain, prepend the key to the record name
		recordName = cmd.Request.Key + "." + recordName
	} else {
		// If no subdomain, just use the key
		recordName = cmd.Request.Key
	}

	fmt.Printf("DEBUG: Creating TXT record - zone='%s', recordName='%s', value='%s'\n", zone, recordName, cmd.Request.Value)

	// Implementation will be moved from cpanel.go
	return h.createTxtRecordAPI(zone, recordName, cmd.Request.Value)
}

// HandleDelete handles deleting a TXT record
func (h *CPanelCommandHandler) HandleDelete(cmd *DeleteTxtRecordCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Extract the actual zone and record name
	zone, recordName := extractZoneAndName(cmd.Request.Domain)
	if recordName != "" {
		// If we have a subdomain, prepend the key to the record name
		recordName = cmd.Request.Key + "." + recordName
	} else {
		// If no subdomain, just use the key
		recordName = cmd.Request.Key
	}

	fmt.Printf("DEBUG: Deleting TXT record - zone='%s', recordName='%s', value='%s'\n", zone, recordName, cmd.Request.Value)

	// Implementation will be moved from cpanel.go
	return h.deleteTxtRecordAPI(zone, recordName, cmd.Request.Value)
}

// HandleEdit handles editing a TXT record
func (h *CPanelCommandHandler) HandleEdit(cmd *EditTxtRecordCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Extract the actual zone and record name
	zone, recordName := extractZoneAndName(cmd.Request.Domain)
	if recordName != "" {
		// If we have a subdomain, prepend the key to the record name
		recordName = cmd.Request.Key + "." + recordName
	} else {
		// If no subdomain, just use the key
		recordName = cmd.Request.Key
	}

	fmt.Printf("DEBUG: Editing TXT record - zone='%s', recordName='%s', oldValue='%s', newValue='%s'\n", 
		zone, recordName, cmd.Request.OldValue, cmd.Request.NewValue)

	// Implementation will be moved from cpanel.go
	return h.editTxtRecordAPI(zone, recordName, cmd.Request.OldValue, cmd.Request.NewValue)
}

// Private helper methods - these will contain the actual cPanel API calls
func (h *CPanelCommandHandler) createTxtRecordAPI(zone, recordName, value string) error {
	// TODO: Move implementation from cpanel.go CreateTxtRecord method
	return fmt.Errorf("not implemented yet")
}

func (h *CPanelCommandHandler) deleteTxtRecordAPI(zone, recordName, value string) error {
	// TODO: Move implementation from cpanel.go DeleteTxtRecord method
	return fmt.Errorf("not implemented yet")
}

func (h *CPanelCommandHandler) editTxtRecordAPI(zone, recordName, oldValue, newValue string) error {
	// TODO: Move implementation from cpanel.go EditTxtRecord method
	return fmt.Errorf("not implemented yet")
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
