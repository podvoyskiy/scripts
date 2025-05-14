package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	h "scripts/components/helper"
	"strconv"
	"strings"
	"sync"
)

//go:embed .env
var env string

var searchDir string
var fileType string
var phone string
var maxGoroutines int

func main() {
	validate()

	mask := fmt.Sprintf("%s_*.csv", fileType)

	files, err := filepath.Glob(filepath.Join(searchDir, mask))
	if err != nil {
		h.Fatal("Error searching files:", err)
	}

	h.Log(h.ColorInfo, "Count files:", len(files))

	if fileType == "num" || maxGoroutines < 2 {
		for _, file := range files {
			searchInFile(file)
		}
	} else {
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, maxGoroutines)

		for _, file := range files {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(f string) {
				defer wg.Done()
				defer func() { <-semaphore }()
				searchInFile(f)
			}(file)
		}

		wg.Wait()
		close(semaphore)
	}
}

func searchInFile(file string) {
	fHandle, err := os.Open(file)
	if err != nil {
		h.Log(h.ColorError, "Error opening file:", err)
		return
	}
	defer fHandle.Close()

	scanner := bufio.NewScanner(fHandle)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), phone) {
			h.Log(h.ColorSuccess, scanner.Text(), file)
		}
	}
	if err := scanner.Err(); err != nil {
		h.Log(h.ColorError, "Error reading file:", err)
	}
}

func validate() {
	config := h.ParseEnv(env)
	if config["NUMBER_SEARCH_DIR"] == "" {
		h.Fatal("NUMBER_SEARCH_DIR not found")
	}
	if _, err := os.Stat(config["NUMBER_SEARCH_DIR"]); os.IsNotExist(err) {
		h.Fatal("Directory", config["NUMBER_SEARCH_DIR"], "does not exist")
	}
	searchDir = config["NUMBER_SEARCH_DIR"]
	h.Log(h.ColorInfo, "Path to files:", searchDir)

	if len(os.Args) < 3 || len(os.Args) > 4 {
		h.Fatal("Incorrect count args")
	}
	if os.Args[1] != "num" && os.Args[1] != "delta" {
		h.Fatal("Incorrect first arg. Only 'num' or 'delta'")
	}
	fileType = strings.ToUpper(os.Args[1])

	if !regexp.MustCompile(`^7\d{10}$`).MatchString(os.Args[2]) {
		h.Fatal("Incorrect second arg. Only phone 11 digits and first digit - 7")
	}
	phone = os.Args[2]

	if len(os.Args) == 4 && fileType == "DELTA" {
		if !regexp.MustCompile(`^\d{1,2}$`).MatchString(os.Args[3]) {
			h.Fatal("Incorrect third arg. Only two digits")
		}
		goroutinesCount, _ := strconv.Atoi(os.Args[3])
		maxGoroutines = goroutinesCount
		if maxGoroutines > 1 {
			h.Log(h.ColorInfo, "Max count goroutines:", maxGoroutines)
		}
	}
}
