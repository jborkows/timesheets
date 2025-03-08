package model

import (
	"errors"
	"io"
	"strings"

	"slices"

	"github.com/BurntSushi/toml"
)

type categories struct {
	Regular  []string
	Overtime []string
}
type holidays struct {
	Repeatable []string
	AddHoc     []string
}

type taskDefinition struct {
	Prefix      string
	OnlyNumbers bool
}
type Config struct {
	Categories categories
	Holidays   holidays
	Tasks      taskDefinition
}

func ReadConfig(r io.Reader) (*Config, error) {
	tomlData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var config Config
	if _, err := toml.Decode(string(tomlData), &config); err != nil {
		return nil, err
	}
	err = config.validate()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func NewConfig(regularCategories []string) *Config {
	return &Config{
		Categories: categories{
			Regular:  regularCategories,
			Overtime: []string{},
		},
		Holidays: holidays{
			Repeatable: []string{},
			AddHoc:     []string{},
		},
		Tasks: taskDefinition{
			Prefix:      "Task-",
			OnlyNumbers: true,
		},
	}
}

func (config *Config) validate() error {
	for _, category := range config.Categories.Regular {
		if category == "" {
			return errors.New("empty category")
		}
		if strings.Contains(category, " ") {
			return errors.New("category cannot contain spaces")
		}
	}
	return nil
}

func (config *Config) IsHoliday(info *DateInfo) bool {
	if slices.Contains(config.Holidays.AddHoc, info.Value) {
		return true
	}
	for _, repeatable := range config.Holidays.Repeatable {
		if !(len(info.Value) == len("2024-12-24")) {
			return false
		}
		if repeatable == info.Value[5:] {
			return true
		}
	}
	return false
}

func insideOfCategory(category string, categories []string) bool {
	return slices.Contains(categories, category)
}

func (config *Config) IsCategory(category string) bool {
	matched := insideOfCategory(category, config.Categories.Regular)
	if matched {
		return true
	}
	return insideOfCategory(category, config.Categories.Overtime)

}
func (config *Config) IsOvertime(category string) bool {
	return insideOfCategory(category, config.Categories.Overtime)

}

func (config *Config) PossibleCategories() []string {
	return append(config.Categories.Regular, config.Categories.Overtime...)
}

func (config *Config) IsTask(text string) bool {
	if !(strings.HasPrefix(text, config.Tasks.Prefix)) {
		return false
	}
	if !config.Tasks.OnlyNumbers {
		return true
	}
	for _, r := range text[len(config.Tasks.Prefix):] {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

const Version = "0.0.2"
const TimeSheetExtension = ".tsf"
