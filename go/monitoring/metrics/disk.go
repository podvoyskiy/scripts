package metrics

import (
	h "scripts/components/helper"
	"strconv"
	"strings"
)

func CheckDiskUsage() int {
	var result int
	output, err := h.ExecCliCmd("df -h | grep -e '/$'")
	if err != nil {
		h.Fatal("Error:", err, "Stderr:", output)
	} else {
		output := strings.ReplaceAll(output, "\n", "")
		diskInfo := strings.Fields(output)
		if len(diskInfo) != 6 {
			h.Fatal("Unexpected output disk info")
		}

		usageDiskPercent, _ := strconv.Atoi(strings.ReplaceAll(diskInfo[4], "%", ""))
		result = usageDiskPercent
	}
	return result
}
