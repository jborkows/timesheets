package lspmessages

type SemanticTokensRequest struct {
	Request
	Params SemanticTokenParams `json:"params"`
}

type SemanticTokenParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type SemanticTokenResponse struct {
	Response
	Result *SemanticTokens `json:"result"`
}

type SemanticTokens struct {
	Data     []int  `json:"data"`
	ResultId string `json:"resultId"`
}

// SemanticTokensRefreshNotification is sent from server to client
// to request that the client re-requests semantic tokens for all buffers.
// See: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#semanticTokens_refreshRequest
type SemanticTokensRefreshNotification struct {
	Notification
}
