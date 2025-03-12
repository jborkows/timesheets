package lspserver

import (
	"github.com/jborkows/timesheets/internal/model"
)

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
