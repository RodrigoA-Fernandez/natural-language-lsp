package lsp

type DiagnosticResponse struct {
	Response
	Result DiagnosticResult `json:"result"`
}

type DiagnosticResult struct {
	Contents MarkupContent `json:"contents"`
}
