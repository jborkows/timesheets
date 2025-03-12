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
