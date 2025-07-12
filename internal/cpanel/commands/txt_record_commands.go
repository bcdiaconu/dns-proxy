package commands

import (
	"fmt"
)

// Execute implements TxtRecordCommand interface
func (cmd *CreateTxtRecordCommand) Execute() error {
	handler := NewCPanelCommandHandler()
	return handler.HandleCreate(cmd)
}

// Execute implements TxtRecordCommand interface
func (cmd *DeleteTxtRecordCommand) Execute() error {
	handler := NewCPanelCommandHandler()
	return handler.HandleDelete(cmd)
}

// Execute implements TxtRecordCommand interface
func (cmd *EditTxtRecordCommand) Execute() error {
	handler := NewCPanelCommandHandler()
	return handler.HandleEdit(cmd)
}

// Validate validates the create command
func (cmd *CreateTxtRecordCommand) Validate() error {
	if cmd.Request.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if cmd.Request.Key == "" {
		return fmt.Errorf("key is required")
	}
	if cmd.Request.Value == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}

// Validate validates the delete command
func (cmd *DeleteTxtRecordCommand) Validate() error {
	if cmd.Request.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if cmd.Request.Key == "" {
		return fmt.Errorf("key is required")
	}
	if cmd.Request.Value == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}

// Validate validates the edit command
func (cmd *EditTxtRecordCommand) Validate() error {
	if cmd.Request.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if cmd.Request.Key == "" {
		return fmt.Errorf("key is required")
	}
	if cmd.Request.OldValue == "" {
		return fmt.Errorf("old value is required")
	}
	if cmd.Request.NewValue == "" {
		return fmt.Errorf("new value is required")
	}
	return nil
}
