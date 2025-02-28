package lspmessages

type HoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}

// TODO something something
type HoverParams struct {
	TextDocumentPositionParams
}

type HoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}

type HoverResult struct {
	Contents string `json:"contents"`
}
