package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

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

type ProcessInfo struct {
	PID  int
	Name string
	CPU  float64
}

func fileExists(path string) (ret bool) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		log.Println(err)
		return
	}

	ret = true
	return
}

func visit(path string, f os.FileInfo, perr error) error {
	if perr != nil {
		return perr
	}

	if f.IsDir() {
		return nil
	}

	/*
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
	*/
	fmt.Println(path)
	return nil
}

func readProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	err := filepath.Walk("/proc", visit)
	if err == nil || errors.Is(err, io.EOF) {
		return processes, nil
	}

	return processes, nil
}

func main() {
	if syscall.Getuid() != 0 {
		fmt.Printf("you must be root to run %q\n", os.Args[0])
		os.Exit(1)
	}

	cfg, err := load()
	if err != nil {
		fmt.Println(err)
	}

	println(cfg.MaxProcess)

	loadavg, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		fmt.Println(err)
	}

	load := strings.Split(string(loadavg), " ")

	loadLast1, err := strconv.ParseFloat(load[0], 64)
	if err != nil {
		fmt.Println(err)
	}

	loadLast5, err := strconv.ParseFloat(load[1], 64)
	if err != nil {
		fmt.Println(err)
	}

	loadLast15, err := strconv.ParseFloat(load[2], 64)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Load average: %.2f %.2f %.2f\n", loadLast1, loadLast5, loadLast15)

	readProcesses()
}
