package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
	// Los tipos de los parámetros se especifican según el tipo de Request
}

type Response struct {
	RPC string `json:"jsonrpc"`
	Id  *int   `json:"id,omitempty"`

	// Result
	// Error
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
