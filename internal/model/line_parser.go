package model

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
)

var (
	ErrEmptyLine       = errors.New("empty line")
	ErrInvalidCategory = errors.New("invalid category")
	ErrInvalidTime     = errors.New("invalid time format. Use X.Y or XhYm (e.g., 1.5 or 1h30m)")
)

type Parser struct {
	HolidayClassifier HolidayClassifier
	IsCategory        func(text string) bool
	IsTask            func(text string) bool
}

func (parser *Parser) ParseLine(dateInfo DateInfo) func(line string) (WorkItem, error) {
	if parser.HolidayClassifier(&dateInfo) {
		return func(line string) (WorkItem, error) {
			return NewHoliday(dateInfo.Value, line)
		}
	}
	return parser.doParseLine
}

type analyzerState int

const (
	StateCategory analyzerState = iota
	StateHours
	StateTask
	StateComment
)

func (s analyzerState) String() string {
	switch s {
	case StateCategory:
		return "category"
	case StateHours:
		return "hours"
	case StateTask:
		return "task"
	case StateComment:
		return "comment"
	default:
		return "unknown" // Important for handling unexpected values
	}
}

type tokenAnalyzer struct {
	*Parser
	tokens []token
	state  analyzerState
	entry  *TimesheetEntry
}

func (analyzer *tokenAnalyzer) analyze(t token) (err error) {
	switch analyzer.state {
	case StateCategory:
		return analyzer.analizeCategory(t)
	case StateHours:
		return analyzer.analizeHours(t)
	case StateTask:
		return analyzer.analizeTask(t)
	case StateComment:
		analyzer.tokens = append(analyzer.tokens, t)
	default:
		return fmt.Errorf("unknown analyzer state: %v", analyzer.state)
	}
	return nil
}

func (analyzer *tokenAnalyzer) finish() *TimesheetEntry {
	var commentBuilder strings.Builder
	for _, t := range analyzer.tokens {
		commentBuilder.WriteString(t.value())
	}
	analyzer.entry.Comment = strings.TrimSpace(commentBuilder.String())
	return analyzer.entry
}

func (analyzer *tokenAnalyzer) analizeTask(t token) error {
	if word, ok := t.(*word); ok {
		if analyzer.IsTask(word.Value) {
			analyzer.entry.Task = &word.Value
		} else {
			analyzer.tokens = append(analyzer.tokens, t)
		}
	} else {
		analyzer.tokens = append(analyzer.tokens, t)
	}
	analyzer.state = StateComment
	return nil
}

func (analyzer *tokenAnalyzer) analizeCategory(t token) error {
	if _, ok := t.(*space); ok {
		if len(analyzer.tokens) == 0 {
			return ErrInvalidCategory
		}
		var categoryBuilder strings.Builder
		for _, tt := range analyzer.tokens {
			categoryBuilder.WriteString(tt.value())
		}
		potentialCategory := categoryBuilder.String()
		if analyzer.IsCategory(potentialCategory) {
			analyzer.entry.Category = potentialCategory
			analyzer.state = StateHours
			analyzer.resetTemp()
		} else {
			return ErrInvalidCategory
		}

	} else {
		analyzer.tokens = append(analyzer.tokens, t)
	}
	return nil
}

func (analyzer *tokenAnalyzer) resetTemp() {
	analyzer.tokens = nil
}

func invalidTime() error {
	debug.PrintStack()
	return ErrInvalidTime
}

func (analyzer *tokenAnalyzer) analizeHours(t token) error {
	temp := analyzer.tokens
	entry := analyzer.entry
	if _, ok := t.(*space); ok {
		if len(temp) == 0 {
			return invalidTime()
		}
		if len(temp) == 1 {
			tt, ok := temp[0].(*number)
			if ok {
				if tt.Value >= 24 {
					return invalidTime()
				}
				entry.Hours = uint8(tt.Value)
				analyzer.resetTemp()
				return nil
			}
			word, ok := temp[0].(*word)
			if !ok {
				return invalidTime()
			}
			var timeBuilder strings.Builder
			for _, leter := range word.Value {

				if leter >= '0' && leter <= '9' {
					timeBuilder.WriteRune(leter)
				} else if leter == 'h' {
					if timeBuilder.Len() == 0 {
						return invalidTime()
					}
					tempValue := parseNumber([]rune(timeBuilder.String()))

					if tempValue >= uint64(24) {
						return invalidTime()
					}
					entry.Hours = uint8(tempValue)
					timeBuilder.Reset()
				} else if leter == 'm' {
					if timeBuilder.Len() == 0 {
						return invalidTime()
					}
					tempValue := parseNumber([]rune(timeBuilder.String()))
					if tempValue >= uint64(60) {
						return invalidTime()
					}
					entry.Minutes = uint8(tempValue)
				} else {
					return invalidTime()
				}
			}
			analyzer.state = StateTask
			return nil
		}
		if len(temp) == 3 {
			hours, ok := temp[0].(*number)
			if !ok {
				return invalidTime()
			}
			_, ok = temp[1].(*dot)
			if !ok {
				return invalidTime()
			}
			minutes, ok := temp[2].(*number)
			if !ok {
				return invalidTime()
			}
			if hours.Value >= 24 || minutes.Value >= 100 {
				return invalidTime()
			}
			entry.Hours = uint8(hours.Value)
			if minutes.Value < 10 {
				entry.Minutes = uint8(minutes.Value * 6)
			} else {
				entry.Minutes = uint8(minutes.Value * 3 / 5)
			}
			analyzer.state = StateTask
			analyzer.resetTemp()
			return nil
		}
		return invalidTime()
	} else {
		analyzer.tokens = append(temp, t)
	}
	return nil
}

func (parser *Parser) doParseLine(line string) (WorkItem, error) {
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return nil, ErrEmptyLine
	}
	var debugMessage strings.Builder
	tokens := tokenize(trimmedLine)
	for _, t := range tokens {
		debugMessage.WriteString(t.represent())
		debugMessage.WriteString(",")
	}
	log.Println(debugMessage.String())

	analyzer := &tokenAnalyzer{
		Parser: parser,
		tokens: make([]token, 0),
		state:  StateCategory,
		entry:  &TimesheetEntry{},
	}
	for _, t := range tokens {
		if err := analyzer.analyze(t); err != nil {
			return nil, err
		}

	}
	return analyzer.finish(), nil

}

func parseNumber(temp []rune) uint64 {
	value, err := strconv.ParseUint(string(temp), 10, 64)
	if err != nil {
		err := fmt.Errorf("error parsing number %v for %s", err, string(temp))
		panic(err)
	}
	return value
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
				tokens = append(tokens, &number{Value: parseNumber(temp)})
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
				tokens = append(tokens, &number{Value: parseNumber(temp)})
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
		tokens = append(tokens, &number{Value: parseNumber(temp)})
	} else if kind == "word" {
		tokens = append(tokens, &word{Value: string(temp)})
	}

	return tokens
}

type token interface {
	represent() string
	value() string
}
type space struct {
}

type word struct {
	Value string
}
type number struct {
	Value uint64
}
type dot struct {
}

func (s *space) represent() string {
	return "space"
}
func (s *space) value() string {
	return " "
}

func (w *word) represent() string {
	return fmt.Sprintf("Word: %s", w.Value)
}
func (w *word) value() string {
	return w.Value
}
func (n *number) represent() string {
	return fmt.Sprintf("Number:%d", n.Value)
}
func (n *number) value() string {
	return fmt.Sprintf("%d", n.Value)
}
func (d *dot) represent() string {
	return "dot"
}
func (d *dot) value() string {
	return "."
}
