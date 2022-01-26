package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Represents some of the values of a line in a Common Apache log
type LogLine struct {
	remoteHost  string
	datetime    time.Time
	Method      string
	route       string
	status      int
	httpVersion int
}

func main() {
	// parse("logs/apache.log")
	setupHTTP()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupHTTP() {
	r := gin.Default()

	r.Use(CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", handler)
	r.POST("/", parseHandler)
	r.POST("/upload", uploadHandler)

	r.Run(":5000")
}

func handler(c *gin.Context) {
	c.JSON(200, "Get DATA")
}

func parseHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatalf("unable to read request body %s", err)
	}

	fmt.Print(string(body))
}

func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}
	fmt.Println(file.Filename)
	f, err := file.Open()
	if err != nil {
		log.Fatalf("unable to open uploaded file %s", err)
	}
	defer f.Close()

	out, err := os.Create("uploadedfile")
	if err != nil {
		log.Fatalf("unable to create file %s", err)
	}
	defer out.Close()
	_, err = io.Copy(out, f)
	if err != nil {
		log.Fatalf("unable to copy to file %s", err)
	}
}

func readLogFile(file string) ([]string, error) {
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

	lines, err := readLogFile(file)
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
			Method:      method,
			route:       route,
			status:      status,
			httpVersion: httpVersion,
		}

		items = append(items, lineData)
	}

	// fmt.Println(items)
	return items
}
