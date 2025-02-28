package lspserver

import (
	"encoding/json"
	"log"

	messages "github.com/jborkows/timesheets/internal/lspmessages"
)

func Route(method string, contents []byte) (interface{}, error) {

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
