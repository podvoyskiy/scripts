package metrics

import (
	h "scripts/components/helper"
	"strings"
)

func CheckCpuUsage() {
	output, err := h.ExecCliCmd("mpstat 1 1 | grep -A 1 'Average' | tail -n 1 | awk '{print 100 - $NF\"%\"}'")
	if err != nil {
		h.Log(h.ColorError, "Error:", err, "Stderr:", output)
	} else {
		output := strings.ReplaceAll(output, "\n", "")
		h.Log(h.ColorInfo, "cpu:", output)
	}
}
