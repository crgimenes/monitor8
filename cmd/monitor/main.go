package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
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

var (
	uptimeSeconds float64
)

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

func visit(path string, f fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if !f.IsDir() {
		return nil
	}

	if !fileExists(path + "/stat") {
		return nil
	}

	pid, err := strconv.Atoi(f.Name())
	if err != nil {
		return nil
	}
	if pid == 0 {
		return nil
	}

	stat, err := os.ReadFile(path + "/stat")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fields := strings.Fields(string(stat))
	/*
		for i, field := range fields {
			fmt.Printf("%d %q\n", i, field)
		}
	*/
	/*
	   #14 utime - CPU time spent in user code, measured in clock ticks
	   #15 stime - CPU time spent in kernel code, measured in clock ticks
	   #16 cutime - Waited-for children's CPU time spent in user code (in clock ticks)
	   #17 cstime - Waited-for children's CPU time spent in kernel code (in clock ticks)
	   #22 starttime - Time when the process started, measured in clock ticks
	*/

	utime, _ := strconv.Atoi(fields[13])
	stime, _ := strconv.Atoi(fields[14])
	cutime, _ := strconv.Atoi(fields[15])
	cstime, _ := strconv.Atoi(fields[16])
	starttime, _ := strconv.ParseFloat(fields[21], 64)

	total_time := utime + stime + cutime + cstime

	const _SYSTEM_CLK_TCK = 100

	Hertz := float64(_SYSTEM_CLK_TCK)

	seconds := uptimeSeconds - (starttime / Hertz)

	if seconds == 0 {
		seconds = 1
	}

	cpu_usage := 100 * ((float64(total_time) / Hertz) / seconds)

	fmt.Printf("PID: %d path %q cpu_usage %v\n", pid, path, cpu_usage)
	return nil
}

func readProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	err := filepath.WalkDir("/proc", visit)
	if err == nil || errors.Is(err, io.EOF) {
		return processes, nil
	}

	return processes, nil
}

func main() {

	// set uid to root
	//syscall.Seteuid(0)

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

	// le o arquivo /proc/uptime e pega o primerio campo, convertendo para int
	uptime, err := os.ReadFile("/proc/uptime")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("uptime", string(uptime))
	uptimeFields := strings.Fields(string(uptime))
	fmt.Println("uptimeFields", uptimeFields)
	uptimeSeconds, err = strconv.ParseFloat(uptimeFields[0], 64)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Uptime: %v seconds\n", uptimeSeconds)

	// readProcesses()

	// read all directories in /proc
	files, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if file.IsDir() {
			if fileExists("/proc/" + file.Name() + "/stat") {
				fmt.Println(file.Name())
			}
		}
	}

}
