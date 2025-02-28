package lspmessages

type DidSaveTextDocumentNotification struct {
	Notification
	Params DidSaveTextDocumentParams `json:"params"`
}

type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
	Text         string           `json:"text"`
}
