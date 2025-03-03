package lspserver

import (
	"encoding/json"
	"fmt"
	"log"

	messages "github.com/jborkows/timesheets/internal/lspmessages"
)

func (self *Controller) HandleMessage(method string, contents []byte) {
	response, err := self.route(method, contents)
	if err != nil {
		log.Printf("Got an error: %s", err)
		return
	}
	if response != nil {
		log.Printf("Sending response for %s", method)
		err = self.writeResponse(response)
		if err != nil {
			log.Printf("Error writing response: %s", err)
		}
	}
}

func (self *Controller) route(method string, contents []byte) (any, error) {

	log.Printf("Received msg with method: %s", method)
	if method == "initialize" {
		log.Printf("Received msg with contents: %s", contents)
	}
	switch method {
	case "initialize":
		var request messages.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			return nil, err
		}

		log.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		msg := messages.NewInitializeResponse(response(request.Request))
		return msg, nil

	case "textDocument/didOpen":
		var request messages.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			return nil, err
		}
		return nil, nil
	case "textDocument/didChange":
		var request messages.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			return nil, err
		}
		self.onChange(&request)
		return nil, nil
	case "textDocument/didSave":
		var request messages.DidSaveTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			return nil, err
		}
		self.onSave(&request)
		return nil, nil

	case "textDocument/completion":
		var request messages.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			return nil, err
		}
		error := self.completion(&request)
		if error != nil {
			return nil, fmt.Errorf("Error getting completions: %w", error)
		}
		return nil, error

	case "textDocument/hover":
		var request messages.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			return nil, err
		}
		err := self.Hover(&request)
		if err != nil {
			return nil, fmt.Errorf("Error getting hover: %w", err)
		}
		return nil, nil
	default:
		return nil, nil
	}
}

func response(request messages.Request) messages.Response {
	return messages.Response{
		RPC: "2.0",
		ID:  &request.ID,
	}
}
