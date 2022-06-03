package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type AlpContent struct {
	Count   int
	Http1XX int
	Http2XX int
	Http3XX int
	Http4XX int
	Http5XX int
	Min     float64
	Max     float64
	Sum     float64
	Avg     float64
	P90     float64
	P95     float64
	P99     float64
}

var (
	alpFormat string
)

func init() {
	flag.StringVar(&alpFormat, "m", "", "alpの-mオプションに渡す値")
}

func getAlpResult(logFile string) (string, error) {
	alpArgs := []string{"ltsv", "-o", "count,1xx,2xx,3xx,4xx,5xx,method,uri,min,max,sum,avg,p90,p95,p99", "--file"}
	alpArgs = append(alpArgs, logFile)
	if alpFormat != "" {
		alpArgs = append(alpArgs, "-m", alpFormat)
	}
	result, err := exec.Command("alp", alpArgs...).Output()
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func parseAlpResult(alpResult string) map[string]AlpContent {
	splittedResult := strings.Split(alpResult, "\n")
	alpMap := make(map[string]AlpContent, len(splittedResult)-4)

	// 上部のヘッダーと最後の罫線を除いてループ
	for _, a := range splittedResult[3 : len(splittedResult)-2] {
		sa := strings.Split(strings.Trim(a, "|"), "|")

		intSa := make([]int, 0, 6)
		floatSa := make([]float64, 0, 7)
		for i, sav := range sa {
			if i < 6 {
				if v, err := strconv.Atoi(strings.TrimSpace(sav)); err == nil {
					intSa = append(intSa, v)
				} else {
					fmt.Println(err)
				}
			} else if i > 7 {
				if v, err := strconv.ParseFloat(strings.TrimSpace(sav), 64); err == nil {
					floatSa = append(floatSa, v)
				} else {
					fmt.Println(err)
				}
			}
		}

		// keyはmethod+uriとする
		key := strings.TrimSpace(sa[6]) + " " + strings.TrimSpace(sa[7])
		alpMap[key] = AlpContent{
			Count:   intSa[0],
			Http1XX: intSa[1],
			Http2XX: intSa[2],
			Http3XX: intSa[3],
			Http4XX: intSa[4],
			Http5XX: intSa[5],
			Min:     floatSa[0],
			Max:     floatSa[1],
			Sum:     floatSa[2],
			Avg:     floatSa[3],
			P90:     floatSa[4],
			P95:     floatSa[5],
			P99:     floatSa[6],
		}
	}

	return alpMap
}

func getDiffAlpMap(oldMap, newMap map[string]AlpContent) map[string]AlpContent {
	diffMap := make(map[string]AlpContent, len(newMap))
	for k, newac := range newMap {
		// NOTE: newMapにあってoldMapにないkey(method+uri)は無視される
		if oldac, ok := oldMap[k]; ok {
			diffMap[k] = AlpContent{
				Count:   newac.Count - oldac.Count,
				Http1XX: newac.Http1XX - oldac.Http1XX,
				Http2XX: newac.Http2XX - oldac.Http2XX,
				Http3XX: newac.Http3XX - oldac.Http3XX,
				Http4XX: newac.Http4XX - oldac.Http4XX,
				Http5XX: newac.Http5XX - oldac.Http5XX,
				Min:     newac.Min - oldac.Min,
				Max:     newac.Max - oldac.Max,
				Sum:     newac.Sum - oldac.Sum,
				Avg:     newac.Avg - oldac.Avg,
				P90:     newac.P90 - oldac.P90,
				P95:     newac.P95 - oldac.P95,
				P99:     newac.P99 - oldac.P99,
			}
		}
	}
	return diffMap
}

func getAlpHeader(alpResult string) []string {
	headerRow := strings.Split(alpResult, "\n")[1]
	splittedHeaderRow := strings.Split(strings.Trim(headerRow, "|"), "|")
	headerItems := make([]string, 0, len(splittedHeaderRow))
	for _, hr := range splittedHeaderRow {
		headerItems = append(headerItems, strings.TrimSpace(hr))
	}
	return headerItems
}

// TODO: リファクタしたくなってきた
// 一気にパースしていい感じの構造体に色々持っておくと楽そう
func getAlpPrintOrder(alpResult string) []string {
	splittedResult := strings.Split(alpResult, "\n")
	printOrder := make([]string, 0, len(splittedResult)-4)
	for _, a := range splittedResult[3 : len(splittedResult)-2] {
		sa := strings.Split(strings.Trim(a, "|"), "|")
		key := strings.TrimSpace(sa[6]) + " " + strings.TrimSpace(sa[7])
		printOrder = append(printOrder, key)
	}
	return printOrder
}

func printAlpDiff(header, printOrder []string, diffMap map[string]AlpContent) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for _, k := range printOrder {
		// printOrderに存在しないkey(method+uri)は無視する
		if dm, ok := diffMap[k]; ok {
			methodAndUri := strings.Split(k, " ")
			row := []string{
				strconv.Itoa(dm.Count),
				strconv.Itoa(dm.Http1XX),
				strconv.Itoa(dm.Http2XX),
				strconv.Itoa(dm.Http3XX),
				strconv.Itoa(dm.Http4XX),
				strconv.Itoa(dm.Http5XX),
				methodAndUri[0],
				methodAndUri[1],
				strconv.FormatFloat(dm.Min, 'f', 3, 64),
				strconv.FormatFloat(dm.Max, 'f', 3, 64),
				strconv.FormatFloat(dm.Sum, 'f', 3, 64),
				strconv.FormatFloat(dm.Avg, 'f', 3, 64),
				strconv.FormatFloat(dm.P90, 'f', 3, 64),
				strconv.FormatFloat(dm.P95, 'f', 3, 64),
				strconv.FormatFloat(dm.P99, 'f', 3, 64),
			}
			table.Append(row)
		}
	}
	table.Render()
}

func getLogFiles(args []string) (string, string, error) {
	if len(args) < 2 {
		return "", "", fmt.Errorf("Error: please pass two args as `./alpdiff <old_log_file> <new_log_file>`")
	}
	isFileExists := func(filepath string) bool {
		_, err := os.Stat(filepath)
		return err == nil
	}
	oldLogFile := args[0]
	newLogFile := args[1]
	if !isFileExists(oldLogFile) {
		return "", "", fmt.Errorf("Error: %s does not exist", oldLogFile)
	} else if !isFileExists(newLogFile) {
		return "", "", fmt.Errorf("Error: %s does not exist", newLogFile)
	}
	return oldLogFile, newLogFile, nil
}

func main() {
	flag.Parse()

	oldLogFile, newLogFile, err := getLogFiles(flag.Args())
	if err != nil {
		fmt.Println(err)
		return
	}

	oldAlpResult, err := getAlpResult(oldLogFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	oldAlpMap := parseAlpResult(oldAlpResult)

	newAlpResult, err := getAlpResult(newLogFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	newAlpMap := parseAlpResult(newAlpResult)

	tableHeader := getAlpHeader(newAlpResult)
	tablePrintOrder := getAlpPrintOrder(newAlpResult)
	diffAlpMap := getDiffAlpMap(oldAlpMap, newAlpMap)

	printAlpDiff(tableHeader, tablePrintOrder, diffAlpMap)
}
