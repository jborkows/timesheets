package model

import (
	"io"
	"strings"

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
	return &config, nil
}

func (config *Config) IsHoliday(info *DateInfo) bool {
	for _, addHoc := range config.Holidays.AddHoc {
		if addHoc == info.Value {
			return true
		}
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
	for _, c := range categories {
		if c == category {
			return true
		}
	}
	return false
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
