package lspmessages

import "github.com/jborkows/timesheets/internal/model"

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	// ... there's tons more that goes here
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}
type ExecuteCommandClientCapabilities struct {
	Commands []string `json:"commands"`
}

type SaveOptions struct {
	IncludeText bool `json:"includeText"`
}

type TextDocumentSyncOptions struct {
	OpenClose bool        `json:"openClose"`
	Change    int         `json:"change"`
	Save      SaveOptions `json:"save"`
}

type ServerCapabilities struct {
	TextDocumentSync   TextDocumentSyncOptions          `json:"textDocumentSync"`
	HoverProvider      bool                             `json:"hoverProvider"`
	DefinitionProvider bool                             `json:"definitionProvider"`
	CodeActionProvider bool                             `json:"codeActionProvider"`
	CompletionProvider map[string]any                   `json:"completionProvider"`
	ExecuteCommand     ExecuteCommandClientCapabilities `json:"executeCommand"`
	ColorProvider      bool                             `json:"colorProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(response Response) InitializeResponse {
	return InitializeResponse{
		Response: response,
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: TextDocumentSyncOptions{
					OpenClose: true,
					Change:    1,
					Save:      SaveOptions{IncludeText: true},
				},
				// HoverProvider:      true,
				// DefinitionProvider: true,
				// CodeActionProvider: true,
				// ColorProvider:      false,
				CompletionProvider: map[string]any{},
				// ExecuteCommand: ExecuteCommandClientCapabilities{
				// 	Commands: []string{"some_command"},
				// },
			},
			ServerInfo: ServerInfo{
				Name:    "timesheets",
				Version: model.Version,
			},
		},
	}
}
