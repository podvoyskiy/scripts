package metrics

import (
	"fmt"
	"regexp"
	h "scripts/components/helper"
	"strconv"
)

type MemoryType string

const (
	MemAvailable MemoryType = "MemAvailable"
	MemTotal     MemoryType = "MemTotal"
)

func CheckMemoryUsage() int {
	memAvailable := getMemoryInGb(MemAvailable)
	memTotal := getMemoryInGb(MemTotal)
	return int((memTotal - memAvailable) / memTotal * 100)
}

func getMemoryInGb(memType MemoryType) float64 {
	var result float64
	output, err := h.ExecCliCmd(fmt.Sprintf("cat /proc/meminfo | grep %v", memType))
	if err != nil {
		h.Fatal("Error:", err, "Stderr:", output)
	} else {
		output := regexp.MustCompile(`\D+`).ReplaceAllString(output, "")
		memoryInKb, _ := strconv.ParseFloat(output, 64)

		result = memoryInKb / 1024 / 1024
	}
	return result
}
