package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
	"crg.eti.br/go/monitor8/proc"
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

	fmt.Printf("cfg.MaxProcess %d\n", cfg.MaxProcess)
	loadLast1, loadLast5, loadLast15, err := proc.GetLoadAverage()

	fmt.Printf("Load average: %.2f %.2f %.2f\n", loadLast1, loadLast5, loadLast15)

	// le o arquivo /proc/uptime e pega o primerio campo, convertendo para int

	uptimeSeconds, err = proc.GetUptime()
	// uptime in hours:minutes:seconds
	uptime := time.Duration(uptimeSeconds * float64(time.Second))
	fmt.Printf("Uptime: %s\n", uptime)

	for {
		fmt.Printf("\033[2J\033[1;1H") // clear screen
		list, err := proc.GetTopProcesses()
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range list[:10] {
			fmt.Printf("%d. %s %.2f%% pid %d\n", k+1, v.Name, v.CPU, v.PID)
		}
		time.Sleep(1 * time.Second)
	}

}
