package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Represents some of the values of a line in a Common Apache log
type LogLine struct {
	remoteHost  string
	datetime    time.Time
	method      string
	route       string
	status      int
	httpVersion int
}

func main() {
	parse("logs/apache.log")
}

func readLog(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func parse(file string) []LogLine {
	var items []LogLine

	lines, err := readLog(file)
	if err != nil {
		log.Fatalf("read lines %s", err)
	}

	for _, line := range lines {

		fields := strings.Fields(line)
		remoteHost := fields[0]
		method := strings.Trim(fields[5], "\"")
		route := fields[6]

		status, err := strconv.Atoi(fields[8])
		if err != nil {
			log.Fatalf("Unable to convert status code to int %s", err)
		}

		layout := "02/Jan/2006:15:04:05 -0700"
		value := strings.Trim(fields[3], "[") + " " + strings.Trim(fields[4], "]")
		datetime, err := time.Parse(layout, value)
		if err != nil {
			fmt.Printf("unable to parse date %s", err)
		}

		httpVersion, err := strconv.Atoi(strings.Trim(fields[7], "HTTP/.0\""))
		if err != nil {
			log.Fatalf("Unable to convert http Version string to int %s", err)
		}

		lineData := LogLine{
			remoteHost:  remoteHost,
			datetime:    datetime,
			method:      method,
			route:       route,
			status:      status,
			httpVersion: httpVersion,
		}

		items = append(items, lineData)
	}

	fmt.Println(items)
	return items
}
