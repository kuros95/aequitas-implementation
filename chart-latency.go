package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func chartItems() [][]opts.LineData {
	itemsY := make([]opts.LineData, 0)
	itemsX := make([]opts.LineData, 0)
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

	for _, line := range fileLines {
		if strings.Contains(line, ":") {
			subs := strings.Split(line, ", ")
			normalized, err := time.ParseDuration(subs[1])
			if err != nil {
				log.Fatalf("failed to parse time: %v, error: %v", subs[1], err)
			}
			intNorm := int(normalized.Milliseconds())
			itemsY = append(itemsY, opts.LineData{Value: strconv.Itoa(intNorm)})
			itemsX = append(itemsX, opts.LineData{Value: subs[0]})
		}
	}

	return [][]opts.LineData{itemsX, itemsY}
}

func chart() {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Opóźnienie dla RPC w czasie w ms",
		}),
	)
	var title string
	var prio string
	if strings.Contains(os.Args[1], "hi") {
		prio = " wysokim"
	} else {
		prio = " niskim"
	}
	if strings.Contains(os.Args[1], "-n-") {
		title = " bez Aequitasa"
	}
	itemsForChart := chartItems()
	line.SetXAxis(itemsForChart[0]).
		AddSeries("Opóźnienie w priorytecie"+prio+title, itemsForChart[1])

	f, _ := os.Create(os.Args[1] + "-latency.html")
	line.Render(f)
}

func main() {
	chartItems()
	chart()
}
