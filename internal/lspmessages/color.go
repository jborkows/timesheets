package lspmessages

type ColorRequest struct {
	Request
	Params ColorParams `json:"params"`
}
type ColorParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type ColorResponse struct {
	Response
	Result []ColorInformation `json:"result"`
}
type ColorInformation struct {
	Range Range `json:"range"`
	Color Color `json:"color"`
}
type Color struct {
	Red   float64 `json:"red"`
	Green float64 `json:"green"`
	Blue  float64 `json:"blue"`
	Alpha float64 `json:"alpha"`
}
