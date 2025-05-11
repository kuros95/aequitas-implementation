package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func chartItems(filePath string) [][]opts.LineData {
	itemsY := make([]opts.LineData, 0)
	itemsX := make([]opts.LineData, 0)
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
		if strings.Contains(line, "KBytes") {
			subs := strings.Split(line, ",")
			second, err := strconv.ParseFloat(subs[2], 32)
			if err != nil {
				log.Fatalf("failed to parse time: %v, error: %v", subs[2], err)
			}
			speed, err := strconv.ParseFloat(subs[5], 32)
			if err != nil {
				log.Fatalf("failed to parse speed: %v, error: %v", subs[5], err)
			}
			itemsY = append(itemsY, opts.LineData{Value: strconv.Itoa(int(speed))})
			itemsX = append(itemsX, opts.LineData{Value: strconv.Itoa(int(second))})
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
			Title:     "Przepływność dla RPC w czasie w Mb/s",
			TextAlign: "center",
		}),
	)
	itemsForChart := chartItems(os.Args[1])
	itemsForChart2 := chartItems(os.Args[2])
	itemsForChart3 := chartItems(os.Args[3])
	itemsForChart4 := chartItems(os.Args[4])
	itemsForChart5 := chartItems(os.Args[5])
	itemsForChart6 := chartItems(os.Args[6])
	line.SetXAxis(itemsForChart[0]).
		AddSeries("Podział 10/90, założenie 10MB/s", itemsForChart[1]).
		AddSeries("Podział 10/90, założenie 90MB/s", itemsForChart2[1]).
		AddSeries("Podział 30/70, założenie 30MB/s", itemsForChart3[1]).
		AddSeries("Podział 30/70, założenie 70MB/s", itemsForChart4[1]).
		AddSeries("Podział 50/50, założenie 50MB/s", itemsForChart5[1]).
		AddSeries("Podział 50/50, założenie 50MB/s", itemsForChart6[1])

	f, _ := os.Create(os.Args[1] + "-latency.html")
	line.Render(f)
}

func main() {
	chartItems(os.Args[1])
	chartItems(os.Args[2])
	chartItems(os.Args[3])
	chartItems(os.Args[4])
	chartItems(os.Args[5])
	chartItems(os.Args[6])
	chart()
}
