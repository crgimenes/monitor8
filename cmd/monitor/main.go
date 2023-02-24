package main

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
	MaxProcess int `json:"max_process" ini:"max_process" cfg:"max_process" cfgDefault:"10" cfgHelper:"Maximum number of processes to collect"`
}

func load() (Config, error) {
	var cfg = Config{}
	config.PrefixEnv = "MONITOR"
	config.File = "monitor.ini"
	err := config.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func main() {
	cfg, err := load()
	if err != nil {
		panic(err)
	}

	println(cfg.MaxProcess)
}
