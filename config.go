package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"syscall"
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

// CheckConfigPerm check perms
func CheckConfigPerm() {
	info, err := os.Stat(_configPath)
	if err != nil {
		fmt.Printf(_cantLoadConfig)
		os.Exit(1)
	}

	if info.Mode().Perm() > 0644 || info.Sys().(*syscall.Stat_t).Uid != 0 {
		fmt.Println(_configPermError)
		os.Exit(1)
	}
}

// LoadConfig from _configPath
func LoadConfig() {
	_, err := toml.DecodeFile(_configPath, &_config)
	if err != nil {
		fmt.Println(_cantLoadConfig)
		os.Exit(1)
	}
}
