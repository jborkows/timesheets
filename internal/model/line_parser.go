package model

import (
	"errors"
	"fmt"
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
	if line == "" {
		return nil, &EmptyLine{Err: errors.New("empty line")}
	}
	return nil, errors.New("not implemented")
}
