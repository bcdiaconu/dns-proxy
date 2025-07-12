package queries

// ListTxtRecordsRequest represents a request to list TXT records
type ListTxtRecordsRequest struct {
	Domain    string
	KeyFilter string // Optional filter by key
}

// TxtRecord represents a TXT DNS record (duplicated from cpanel package for now)
type TxtRecord struct {
	Line  int    `json:"line"`
	Key   string `json:"key"`   // The record name without the zone
	Value string `json:"value"` // The txtdata
	Name  string `json:"name"`  // Full name including zone
}

// TxtRecordQuery represents a query for TXT records
type TxtRecordQuery interface {
	Execute() (interface{}, error)
}

// ListTxtRecordsQuery handles listing TXT records
type ListTxtRecordsQuery struct {
	Request ListTxtRecordsRequest
}

// QueryHandler handles query execution
type QueryHandler interface {
	HandleList(query *ListTxtRecordsQuery) ([]TxtRecord, error)
}
