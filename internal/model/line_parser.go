package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type HolidayClassifier = func(aDate *DateInfo) bool

type EmptyLine struct {
	Err error
}

func (e *EmptyLine) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

type DateInfo struct {
	Value string
}

type Parser struct {
	HolidayClassifier HolidayClassifier
	IsCategory        func(text string) bool
	IsTask            func(text string) bool
}

func (parser *Parser) ParseLine(dateInfo DateInfo) func(line string) (WorkItem, error) {
	if parser.HolidayClassifier(&dateInfo) {
		return func(line string) (WorkItem, error) {
			return NewHoliday(dateInfo.Value)
		}
	}
	return doParseLine
}

func doParseLine(line string) (WorkItem, error) {
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return nil, &EmptyLine{Err: errors.New("empty line")}
	}
	var debugMessage strings.Builder
	tokens := tokenize(trimmedLine)
	for _, t := range tokens {
		debugMessage.WriteString(t.represent())
		debugMessage.WriteString(",")
	}
	log.Println(debugMessage.String())
	return nil, errors.New("not implemented")
}

func parseNumber(temp []rune) uint8 {
	value, _ := strconv.ParseUint(string(temp), 10, 8)
	return uint8(value)
}

func tokenize(line string) []token {
	var tokens []token = make([]token, 0)
	var lineAsRunes = []rune(line)
	var temp []rune
	kind := ""

	for i := 0; i < len(lineAsRunes); i++ {
		if lineAsRunes[i] >= '0' && lineAsRunes[i] <= '9' {
			if kind == "" {
				kind = "number"
				temp = []rune{lineAsRunes[i]}
			} else if kind == "space" {
				kind = "number"
				temp = []rune{lineAsRunes[i]}
			} else {
				temp = append(temp, lineAsRunes[i])
			}
			continue
		}
		if lineAsRunes[i] == ' ' {
			if kind == "space" {
				continue
			}
			if kind == "number" {
				tokens = append(tokens, &number{value: parseNumber(temp)})
			} else if kind == "word" {
				tokens = append(tokens, &word{Value: string(temp)})
			}
			temp = []rune{}
			kind = "space"
			tokens = append(tokens, &space{})
			continue
		}
		if lineAsRunes[i] == '.' {
			if kind == "number" {
				tokens = append(tokens, &number{value: parseNumber(temp)})
			} else if kind == "word" {
				tokens = append(tokens, &word{Value: string(temp)})
			}
			tokens = append(tokens, &dot{})
			temp = []rune{}
			kind = ""
			continue
		}
		if kind == "" {
			temp = []rune{lineAsRunes[i]}
		} else if kind == "space" {
			temp = []rune{lineAsRunes[i]}
		} else {
			temp = append(temp, lineAsRunes[i])
		}
		kind = "word"
	}
	if kind == "number" {
		tokens = append(tokens, &number{value: parseNumber(temp)})
	} else if kind == "word" {
		tokens = append(tokens, &word{Value: string(temp)})
	}

	return tokens
}

type token interface {
	represent() string
}
type space struct {
}

type word struct {
	Value string
}
type number struct {
	value uint8
}
type dot struct {
}

func (s *space) represent() string {
	return "space"
}

func (w *word) represent() string {
	return fmt.Sprintf("Word: %s", w.Value)
}
func (n *number) represent() string {
	return fmt.Sprintf("Number:%d", n.value)
}
func (d *dot) represent() string {
	return "dot"
}
