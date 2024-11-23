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

type InvalidCategory struct {
	Err error
}

func (e *InvalidCategory) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

type InvalidTime struct {
	Err error
}

func (e *InvalidTime) Error() string {
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
	return parser.doParseLine
}

type tokenAnalyzer struct {
	*Parser
	temp  []token
	state string
	entry *TimesheetEntry
}

func (analyzer *tokenAnalyzer) analyze(token token) (err error) {
	if analyzer.state == "category" {
		return analyzer.analizeCategory(token)
	} else if analyzer.state == "hours" {
		return analyzer.analizeHours(token)
	} else if analyzer.state == "task" {
		//TODO: implement task
		aaa := "AAAAAA"
		analyzer.entry.Task = &aaa
		return nil
	}
	return nil

}
func (analyzer *tokenAnalyzer) analizeCategory(t token) error {
	if _, ok := t.(*space); ok {
		if len(analyzer.temp) == 0 {
			return &InvalidCategory{Err: errors.New("empty category")}
		}
		var categoryBuilder strings.Builder
		for _, tt := range analyzer.temp {
			categoryBuilder.WriteString(tt.value())
		}
		potentialCategory := categoryBuilder.String()
		if analyzer.IsCategory(potentialCategory) {
			analyzer.entry.Category = potentialCategory
			analyzer.state = "hours"
			analyzer.resetTemp()
		} else {
			return &InvalidCategory{Err: errors.New("invalid category")}
		}

	} else {
		analyzer.temp = append(analyzer.temp, t)
	}
	return nil
}

func (analyzer *tokenAnalyzer) resetTemp() {
	analyzer.temp = make([]token, 0)
}

func invalidTime() error {
	return &InvalidTime{Err: errors.New("invalid time use 1.5 or 1h30m")}
}

func (analyzer *tokenAnalyzer) analizeHours(t token) error {
	temp := analyzer.temp
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
			if ok {
				var timeBuilder strings.Builder
				for _, leter := range word.Value {

					if leter >= '0' && leter <= '9' {
						timeBuilder.WriteRune(leter)
					} else if leter == 'h' {
						if timeBuilder.Len() == 0 {
							return invalidTime()
						}
						tempValue := parseNumber([]rune(timeBuilder.String()))

						if tempValue >= 24 {
							return invalidTime()
						}
						entry.Hours = uint8(tempValue)
						timeBuilder.Reset()
					} else if leter == 'm' {
						if timeBuilder.Len() == 0 {
							return invalidTime()
						}
						tempValue := parseNumber([]rune(timeBuilder.String()))
						if tempValue >= 60 {
							return invalidTime()
						}
						entry.Minutes = uint8(tempValue)
					} else {
						return invalidTime()
					}
					timeBuilder.WriteRune(leter)
				}
			}

			var hoursBuilder strings.Builder
			for _, tt := range temp {
				hoursBuilder.WriteString(tt.value())
			}
			potentialHours := hoursBuilder.String()
			if _, err := strconv.ParseFloat(potentialHours, 64); err == nil {
				analyzer.state = "task"
				analyzer.resetTemp()
			} else {
				return invalidTime()
			}
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
			analyzer.state = "task"
			analyzer.resetTemp()
			return nil
		}
		return invalidTime()
	} else {
		analyzer.temp = append(temp, t)
	}
	return nil
}

func (parser *Parser) doParseLine(line string) (WorkItem, error) {
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

	analyzer := &tokenAnalyzer{
		Parser: parser,
		temp:   make([]token, 0),
		state:  "category",
		entry:  &TimesheetEntry{},
	}
	for _, t := range tokens {
		if err := analyzer.analyze(t); err != nil {
			return nil, err
		}

	}

	return analyzer.entry, nil
}

func parseNumber(temp []rune) uint64 {
	value, _ := strconv.ParseUint(string(temp), 10, 8)
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
