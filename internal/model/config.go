package model

import (
	"github.com/BurntSushi/toml"
	"io"
)

type categories struct {
	Regular  []string
	Overtime []string
}
type holidays struct {
	Repeatable []string
	AddHoc     []string
}

type Config struct {
	Categories categories
	Holidays   holidays
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
