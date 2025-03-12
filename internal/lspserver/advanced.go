package lspserver

import (
	"log"

	messages "github.com/jborkows/timesheets/internal/lspmessages"
)

func (self *Controller) Formatting(request *messages.FormattingRequest) error {

	textEdits := []messages.TextEdit{}
	content := self.content.get(request.Params.TextDocument.URI)
	log.Printf("Formatting content: '%s' %d", content, len(content))
	for i, line := range content {
		log.Printf("Line %d: '%s'", i, line)
	}
	if len(content) == 0 {
		msg := messages.FormattingResponse{
			Response: response(request.Request),
			Result:   textEdits,
		}
		log.Printf("Formatting response: %+v", msg)
		return self.writeResponse(msg)
	}

	lastLine := content[len(content)-1]
	if len(lastLine) > 0 && lastLine[len(lastLine)-1] != '\n' {
		textEdits = append(textEdits, messages.TextEdit{
			Range: messages.Range{
				Start: messages.Position{
					Line:      len(content) - 1,
					Character: len(lastLine),
				},
				End: messages.Position{
					Line:      len(content) - 1,
					Character: len(lastLine),
				},
			},
			NewText: "\n",
		})
	}

	for i, line := range content {
		if len(line) == 0 {
			continue
		}

		endSpaceIndex := -1
		for j := len(line) - 1; j > 0; j-- {
			if line[j] == ' ' && endSpaceIndex == -1 {
				endSpaceIndex = j
			}
			if line[j] != ' ' {
				startSpaceIndex := j + 1
				if endSpaceIndex-startSpaceIndex > 0 {
					textEdits = append(textEdits, messages.TextEdit{
						Range: messages.Range{
							Start: messages.Position{
								Line:      i,
								Character: startSpaceIndex,
							},
							End: messages.Position{
								Line:      i,
								Character: endSpaceIndex,
							},
						},
						NewText: "",
					})
				}
				endSpaceIndex = -1
			}
		}

	}

	msg := messages.FormattingResponse{
		Response: response(request.Request),
		Result:   textEdits,
	}
	log.Printf("Formatting response: %+v", msg)
	return self.writeResponse(msg)

}

func (self *Controller) SemanticTokens(request *messages.SemanticTokensRequest) error {
	textDocument := request.Params.TextDocument
	date, err := self.service.ParseDateFromName(request.Params.TextDocument.URI)
	if err != nil {

		msg := messages.SemanticTokenResponse{
			Response: response(request.Request),
			Result: &messages.SemanticTokens{
				Data: []int{},
			},
		}
		return self.writeResponse(msg)
	}
	content := self.content.get(textDocument.URI)
	tokens := self.service.SemanaticTokenFrom(content, date)
	tokensToSend := TranslateSemanticTokens(tokens)
	log.Printf("Send %+v", tokensToSend)
	msg := messages.SemanticTokenResponse{
		Response: response(request.Request),
		Result: &messages.SemanticTokens{
			Data: tokensToSend,
		},
	}
	return self.writeResponse(msg)
}
