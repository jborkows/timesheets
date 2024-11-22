package model

import (
	"errors"
	"fmt"
)

type CategeryProvider interface {
	GetCategories() []CategoryType
}

type EmptyLine struct {
	Err error
}

func (e *EmptyLine) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func ParseLine(line string, cp CategeryProvider) (*WorkItem, error) {

	return nil, &EmptyLine{Err: errors.New("empty line")}
}
