package proc

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
