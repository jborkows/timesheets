package lspserver

import (
	"log"

	"github.com/jborkows/timesheets/internal/model"

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

func TranslateSemanticTokens(tokens []model.TokenReady) []int {
	if len(tokens) == 0 {
		return []int{}
	}

	var tokensToSend []int = []int{tokens[0].Line, tokens[0].Column, tokens[0].Length, int(tokens[0].Type), 0}
	for i := 1; i < len(tokens); i++ {
		prev := tokens[i-1]
		current := tokens[i]
		tokensToSend = append(tokensToSend, current.Line-prev.Line)
		if current.Line-prev.Line == 0 {
			tokensToSend = append(tokensToSend, current.Column-prev.Column)
		} else {
			tokensToSend = append(tokensToSend, current.Column)
		}
		tokensToSend = append(tokensToSend, current.Length)
		tokensToSend = append(tokensToSend, int(current.Type))
		tokensToSend = append(tokensToSend, 0)
	}
	return tokensToSend
}
