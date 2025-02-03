package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	filePath := os.Args[1]
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	hiFile, err := os.Create(readFile.Name() + "-hi-latency.csv")
	if err != nil {
		fmt.Printf("failed to create a file for logging, error: %v", err)
		os.Exit(1)
	}

	loFile, err := os.Create(readFile.Name() + "-lo-latency.csv")
	if err != nil {
		fmt.Printf("failed to create a file for logging, error: %v", err)
		os.Exit(1)
	}

	hiFile.WriteString("time, latency\n")
	loFile.WriteString("time, latency\n")
	for _, line := range fileLines {
		if strings.Contains(line, "completed") {
			subs := strings.Fields(line)
			if subs[10] == "hi" {
				writeIn := subs[1] + ", " + subs[12] + "\n"
				hiFile.WriteString(writeIn)
			} else if subs[10] == "lo" {
				writeIn := subs[1] + ", " + subs[12] + "\n"
				loFile.WriteString(writeIn)
			}
		}
	}

}
