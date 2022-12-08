package main

import (
	"path/filepath"

	"github.com/BurntSushi/toml"

	"utm5-sber-osmp/internal/service"
)

func LoadConfig(confFile string) (conf *service.Config) {
	if _, err := toml.DecodeFile(filepath.Clean(confFile), &conf); err != nil {
		panic(err)
	}
	return conf
}
