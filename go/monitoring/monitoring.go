package main

import (
	_ "embed"
	"fmt"
	"os"
	h "scripts/components/helper"
	"scripts/components/telegram"
	"scripts/monitoring/metrics"
	"strconv"
)

//go:embed .env
var env string

func main() {
	config := h.ParseEnv(env)
	maxMemoryPercent, err := strconv.Atoi(config["ALERT_MEMORY_PERCENT_USAGE"])
	if err != nil {
		h.Fatal(err)
	}
	maxDiskPercent, err := strconv.Atoi(config["ALERT_DISK_PERCENT_USAGE"])
	if err != nil {
		h.Fatal(err)
	}

	bot := telegram.NewBot(os.Getenv("TELEGRAM_TOKEN"), os.Getenv("TELEGRAM_CHAT_ID"))

	usageDiskPercent := metrics.CheckDiskUsage()
	if usageDiskPercent > maxDiskPercent {
		h.Log(h.ColorWarning, "disk:", usageDiskPercent)
		bot.Send(fmt.Sprintf("disk usage: %d percent", usageDiskPercent))
	} else {
		h.Log(h.ColorInfo, "disk:", usageDiskPercent, "%")
	}

	usageMemoryPercent := metrics.CheckMemoryUsage()
	if usageMemoryPercent > maxMemoryPercent {
		h.Log(h.ColorWarning, "memory:", usageMemoryPercent, "%")
		bot.Send(fmt.Sprintf("memory usage: %d percent", usageMemoryPercent))
	} else {
		h.Log(h.ColorInfo, "memory:", usageMemoryPercent, "%")
	}

	metrics.CheckCpuUsage()
}
