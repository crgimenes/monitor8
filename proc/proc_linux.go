package proc

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ProcessInfo struct {
	PID          int
	Name         string
	CPU          float64
	lastReadTime float64
}

// GetLoadAverage returns the load average for the last 1, 5 and 15 minutes.
func GetLoadAverage() (float64, float64, float64, error) {
	loadavg, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		fmt.Println(err)
	}

	load := strings.Split(string(loadavg), " ")

	loadLast1, err := strconv.ParseFloat(load[0], 64)
	if err != nil {
		return 0, 0, 0, err
	}

	loadLast5, err := strconv.ParseFloat(load[1], 64)
	if err != nil {
		return 0, 0, 0, err
	}

	loadLast15, err := strconv.ParseFloat(load[2], 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return loadLast1, loadLast5, loadLast15, nil
}

// GetUptime returns the uptime in seconds.
func GetUptime() (float64, error) {
	uptime, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	uptimeFields := strings.Fields(string(uptime))
	return strconv.ParseFloat(uptimeFields[0], 64)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		log.Println(err)
		return false
	}
	return true
}

var processMap map[int]ProcessInfo

// getProcessInfo returns the CPU usage, PID and name of a process.
func getProcessInfo(fields []string, uptimeSeconds float64) (float64, int, string, error) {
	if processMap == nil {
		processMap = make(map[int]ProcessInfo)
	}
	utime, err := strconv.Atoi(fields[13])
	if err != nil {
		return 0, 0, "", err
	}

	stime, err := strconv.Atoi(fields[14])
	if err != nil {
		return 0, 0, "", err
	}

	cutime, err := strconv.Atoi(fields[15])
	if err != nil {
		return 0, 0, "", err
	}

	cstime, err := strconv.Atoi(fields[16])
	if err != nil {
		return 0, 0, "", err
	}

	starttime, err := strconv.ParseFloat(fields[21], 64)
	if err != nil {
		return 0, 0, "", err
	}

	totalTime := utime + stime + cutime + cstime
	//totalTime = utime + stime

	const _SYSTEM_CLK_TCK = 100

	Hertz := float64(_SYSTEM_CLK_TCK)

	seconds := uptimeSeconds - (starttime / Hertz)

	if seconds == 0 {
		seconds = 1
	}

	pid, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, "", err
	}

	//cpuUsage := 100.0 * ((float64(totalTime) / Hertz) / seconds)
	cpuUsage := 1000.0 * (((float64(totalTime) - processMap[pid].lastReadTime) / Hertz) / seconds)

	processMap[pid] = ProcessInfo{
		lastReadTime: float64(totalTime),
	}

	name := fields[1]

	return cpuUsage, pid, name, nil
}

// GetTopProcesses returns the top processes.
func GetTopProcesses() ([]ProcessInfo, error) {
	// read all directories in /proc
	files, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Println(err)
	}

	list := make([]ProcessInfo, len(files))

	uptime, err := GetUptime()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir() {
			if fileExists("/proc/" + file.Name() + "/stat") {

				stat, err := os.ReadFile("/proc/" + file.Name() + "/stat")
				if err != nil {
					fmt.Println(err)
				}

				fields := strings.Fields(string(stat))

				cpuUsage, pid, name, err := getProcessInfo(fields, uptime)
				if err != nil {
					return nil, err
				}

				list[i] = ProcessInfo{
					PID:  pid,
					Name: name,
					CPU:  cpuUsage,
				}
			}
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].CPU > list[j].CPU
	})

	return list[:10], nil
}
