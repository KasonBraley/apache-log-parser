package main

import (
	"bufio"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Represents some of the values of a line in a Common Apache log
type LogLine struct {
	gorm.Model
	RemoteHost  string
	DateTime    time.Time
	Method      string
	Route       string
	Status      int
	HTTPVersion int
}

func main() {
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

	r.GET("/", retrieveDataHandler)
	r.POST("/upload", uploadHandler)

	r.Run(":5000")
}

func retrieveDataHandler(c *gin.Context) {
	db := connectDB()
	logLines := retrieveAllRowsFromDB(db)

	c.JSON(200, gin.H{
		"data": logLines,
	})
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

	content, err := readLog(file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to open and read log",
		})
		return
	}

	logLines, err := parseLog(content)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to parse values from the log",
		})
		return
	}

	db := connectDB()
	storeData(db, logLines)

	c.JSON(200, "")
}

func readLog(file *multipart.FileHeader) ([]string, error) {
	f, err := file.Open()
	if err != nil {
		log.Printf("unable to open uploaded file %s", err)
		return []string{}, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func parseLog(lines []string) ([]LogLine, error) {
	var items []LogLine

	for _, line := range lines {

		fields := strings.Fields(line)
		remoteHost := fields[0]
		method := strings.Trim(fields[5], "\"")
		route := fields[6]

		status, err := strconv.Atoi(fields[8])
		if err != nil {
			log.Printf("Unable to convert status code to int %s", err)
			return []LogLine{}, err
		}

		layout := "02/Jan/2006:15:04:05 -0700"
		value := strings.Trim(fields[3], "[") + " " + strings.Trim(fields[4], "]")
		datetime, err := time.Parse(layout, value)
		if err != nil {
			log.Printf("unable to parse date %s", err)
			return []LogLine{}, err
		}

		httpVersion, err := strconv.Atoi(strings.Trim(fields[7], "HTTP/.0\""))
		if err != nil {
			log.Printf("Unable to convert http Version string to int %s", err)
			return []LogLine{}, err
		}

		lineData := LogLine{
			RemoteHost:  remoteHost,
			DateTime:    datetime,
			Method:      method,
			Route:       route,
			Status:      status,
			HTTPVersion: httpVersion,
		}

		items = append(items, lineData)
	}

	return items, nil
}

func connectDB() *gorm.DB {
	dsn := "host=localhost user=kason password=pass dbname=apache_logs port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func storeData(db *gorm.DB, logLines []LogLine) {

	// Migrate the schema
	db.AutoMigrate(&LogLine{})

	for _, line := range logLines {
		db.Create(&line)
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

func retrieveAllRowsFromDB(db *gorm.DB) []LogLine {
	var logLines []LogLine

	// Get all records
	db.Find(&logLines)

	return logLines
}
