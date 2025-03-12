package lspserver

import (
	"log"
	"strings"

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

type Token struct {
	Index int
	Word  string
}

func tokenizeFromIndex(input string, j int) []Token {
	// Ensure j is within bounds
	if j < 0 || j >= len(input) {
		return nil
	}

	// Get substring from index j
	substr := input[j:]

	// Split by spaces
	words := strings.Fields(substr)

	// Prepare tokens
	tokens := make([]Token, 0, len(words))
	index := j // Start index tracking from j

	//Chats code... does not include multiple spaces in the token
	for _, word := range words {
		tokens = append(tokens, Token{Index: index, Word: word})
		index += len(word) + 1 // Move index past word and space
	}

	return tokens
}

type TokenReady struct {
	line       int
	column     int
	length     int
	typeLegend int
}

func (self *Controller) SemanticTokens(request *messages.SemanticTokensRequest) error {
	textDocument := request.Params.TextDocument
	date, err := self.service.ParseDateFromName(request.Params.TextDocument.URI)
	tokens := []TokenReady{}
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
	if len(content) < 1 {
		msg := messages.SemanticTokenResponse{
			Response: response(request.Request),
			Result: &messages.SemanticTokens{
				Data: []int{},
			},
		}
		return self.writeResponse(msg)
	}

	/**
	{ line: 2, startChar: 10, length: 4, tokenType: "type", tokenModifiers: [] }, tokenType is an index from legend, tokenModifiers is an array of indices from legend but see in spec how it is calculated into uint32
		**/
	for i, line := range content {
		parsed := self.service.ParseLine(line, date)
		if parsed == nil {
			continue
		}
		switch e := parsed.(type) {
		case *model.TimesheetEntry:
			err := e.Validate()
			if err != nil {
				return nil
			}
			tokens = append(tokens, TokenReady{
				line:       i,
				column:     0,
				length:     len(e.Category),
				typeLegend: 0,
			})
			startOfTime := 0

			// endOfTime := 0
			for index, char := range line {
				log.Printf("Char: %s %d\n", string(char), index)
			}
			log.Printf("After category: '%s', %d", line[len(e.Category):], len(e.Category))
			j := len(e.Category)
			for ; j < len(line); j++ {
				if line[j] != ' ' && startOfTime == 0 {
					startOfTime = j
					continue
				}
				if line[j] == ' ' && startOfTime > 0 {
					log.Printf("Time: '%s'  -> %d %d %d", line[startOfTime:j], startOfTime, j, j-startOfTime)
					tokens = append(tokens, TokenReady{
						line:       i,
						column:     startOfTime,
						length:     j - startOfTime,
						typeLegend: 1,
					})
					break
				}
			}
			words := tokenizeFromIndex(line, j)
			for _, tocken := range words {
				tokens = append(tokens, TokenReady{
					line:       i,
					column:     tocken.Index + 1,
					length:     len(tocken.Word),
					typeLegend: 2,
				})
				log.Printf("Token %d: %s %d", i, tocken.Word, tocken.Index+1)
			}

		default:
			continue
		}

	}
	if len(tokens) == 0 {
		msg := messages.SemanticTokenResponse{
			Response: response(request.Request),
			Result: &messages.SemanticTokens{
				Data: []int{},
			},
		}
		return self.writeResponse(msg)
	}

	var tokensToSend []int = []int{tokens[0].line, tokens[0].column, tokens[0].length, tokens[0].typeLegend, 0}
	for i := 1; i < len(tokens); i++ {
		prev := tokens[i-1]
		current := tokens[i]
		tokensToSend = append(tokensToSend, current.line-prev.line)
		if current.line-prev.line == 0 {
			tokensToSend = append(tokensToSend, current.column-prev.column)
		} else {
			tokensToSend = append(tokensToSend, current.column)
		}
		tokensToSend = append(tokensToSend, current.length)
		tokensToSend = append(tokensToSend, current.typeLegend)
		tokensToSend = append(tokensToSend, 0)
	}
	log.Printf("Send %+v", tokensToSend)

	msg := messages.SemanticTokenResponse{
		Response: response(request.Request),
		Result: &messages.SemanticTokens{
			Data: tokensToSend,
		},
	}
	return self.writeResponse(msg)
}
