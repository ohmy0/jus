package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Permits []Permit `toml:"permit"`
}

type Permit struct {
	User    string   `toml:"user"`
	As      string   `toml:"as"`
	Paths   []string `toml:"paths"`
	KeepEnv bool     `toml:"keepenv"`
}

// LoadConfig from _configPath
func LoadConfig() {
	_, err := toml.DecodeFile(_configPath, &_config)
	if err != nil {
		fmt.Println(_cantLoadConfig)
		os.Exit(1)
	}
}
