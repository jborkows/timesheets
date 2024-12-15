package model

import "fmt"

type InvalidTime struct {
	Err error
}

func (e *InvalidTime) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

type DateInfo struct {
	Value string
}

type HolidayClassifier = func(aDate *DateInfo) bool
