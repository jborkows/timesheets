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

func ParseLine(dateInfo DateInfo, hollidayClassifier HolidayClassifier) func(line string) (WorkItem, error) {
	if hollidayClassifier(&dateInfo) {
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
