package lsp

type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

type MarkupKind string

const (
	Plaintext MarkupKind = "plaintext"
	Markdown  MarkupKind = "markdown"
)
