package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/charts"
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
			itemsY = append(itemsY, opts.LineData{Value: subs[1]})
		} else {
			subs := strings.Split(line, ", ")
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
			Title: "Opóźnienie dla RPC w czasie",
		}),
	)
	itemsForChart := chartItems()
	line.SetXAxis(itemsForChart[0]).
		AddSeries("Opóźnienie", itemsForChart[1])

	f, _ := os.Create(os.Args[1] + "-latency.html")
	line.Render(f)
}

func main() {
	chartItems()
	chart()
}
