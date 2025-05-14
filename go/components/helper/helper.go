package helper

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ColorText int

const ColorError = ColorText(31)
const ColorSuccess = ColorText(32)
const ColorWarning = ColorText(33)
const ColorInfo = ColorText(36)

func Log(color ColorText, args ...any) {
	var msg string
	for _, arg := range args {
		msg += fmt.Sprintf("%v ", arg)
	}
	log.Printf("\033[%dm%s\033[0m", color, msg)
}

func Fatal(args ...any) {
	Log(ColorError, args)
	os.Exit(1)
}

func ParseEnv(content string) map[string]string {
	lines := strings.Split(content, "\n")
	config := make(map[string]string)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}
	//Log(ColorSuccess, "Config:", config)
	return config
}

func ExecCliCmd(cmd string) (string, error) {
	command := exec.Command("sh", "-c", cmd)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()

	if err != nil {
		return stderr.String(), err
	}

	return stdout.String(), nil
}
