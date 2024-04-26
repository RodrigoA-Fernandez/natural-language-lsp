package lsp

type DidOpenTextDocumentNotification struct {
  Notification
  Params DidOpenTextDocumentNotificationParams `json:"params"`
}

type DidOpenTextDocumentNotificationParams struct {
  TextDocument TextDocumentItem `json:"textDocument"`
}
