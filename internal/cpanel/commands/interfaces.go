package commands

// CreateTxtRecordRequest represents a request to create a TXT record
type CreateTxtRecordRequest struct {
	Domain string
	Key    string
	Value  string
}

// DeleteTxtRecordRequest represents a request to delete a TXT record
type DeleteTxtRecordRequest struct {
	Domain string
	Key    string
	Value  string
}

// EditTxtRecordRequest represents a request to edit a TXT record
type EditTxtRecordRequest struct {
	Domain   string
	Key      string
	OldValue string
	NewValue string
}

// TxtRecordCommand represents a command that modifies TXT records
type TxtRecordCommand interface {
	Execute() error
}

// CreateTxtRecordCommand handles creating TXT records
type CreateTxtRecordCommand struct {
	Request CreateTxtRecordRequest
}

// DeleteTxtRecordCommand handles deleting TXT records
type DeleteTxtRecordCommand struct {
	Request DeleteTxtRecordRequest
}

// EditTxtRecordCommand handles editing TXT records
type EditTxtRecordCommand struct {
	Request EditTxtRecordRequest
}

// CommandHandler handles command execution
type CommandHandler interface {
	HandleCreate(cmd *CreateTxtRecordCommand) error
	HandleDelete(cmd *DeleteTxtRecordCommand) error
	HandleEdit(cmd *EditTxtRecordCommand) error
}
