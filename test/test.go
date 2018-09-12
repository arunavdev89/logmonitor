package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	methods  = []string{"GET", "POST"}
	status   = []string{"200", "500", "404", "503", "401"}
	websites = []string{
		"/api/v1/users/arun",
		"/api/v1/users/dorothy",
		"/api/v1/users/gopher",
		"/api/v1/products/soap",
		"/api/v1/products/cookies",
		"/api/v1/products/donuts",
		"/api/v1/usage/products",
		"/api/v1/usage/rate",
		"/api/v1/usage/users",
		"/api/v1/accounts/arun",
		"/api/v1/accounts/dorothy",
		"/api/v1/accounts/gopher",
		"/api/v1/login",
		"/api/v1/metrics",
		"/api/v1/config",
		"/api/v1/usage"}
	callerIds = []string{"arun", "dorothy", "harry", "marie", "raymond"}
)

type FileWriter struct {
	file *os.File
}

func NewFileWriter(f *os.File) *FileWriter {
	return &FileWriter{
		file: f,
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (f *FileWriter) GetRandomLogLine() string {
	return fmt.Sprintf("%s - %s [%s] \"%s %s %s\" %s %s",
		"127.0.0.1",
		GetRandomItemFromList(callerIds),
		time.Now().Format(time.RFC3339),
		GetRandomItemFromList(methods),
		GetRandomItemFromList(websites),
		"Http/1.1",
		GetRandomItemFromList(status),
		"2390")
}

func (f *FileWriter) GetLogLine(website string, method string, status string) string {
	return fmt.Sprintf("%s - %s [%s] \"%s %s %s\" %s %s",
		"127.0.0.1",
		GetRandomItemFromList(callerIds),
		time.Now().Format(time.RFC3339),
		method,
		website,
		"Http/1.1",
		status,
		"2390")
}

func (f *FileWriter) Write(line string) {
	_, err := f.file.WriteString(line)
	if err != nil {
		fmt.Println(err)
	}
}

func (f *FileWriter) BurstWriteForNSec(n int, website string, method string, status string) {
	timer := time.NewTimer(time.Duration(n) * time.Second)
	done := make(chan bool)
	//w := bufio.NewWriter(f.file)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Done writing for %s ", n, " seconds")

			default:
				f.Write(f.GetLogLine(website, method, status))
				f.Write("\n")
				f.Write(f.GetRandomLogLine())
				f.Write("\n")
				f.file.Sync()
				//Uncomment if high CPU usage
				//time.Sleep(n/10)
			}

		}
	}()

	<-timer.C
	done <- true
}

func GetRandomItemFromList(s []string) string {
	return s[rand.Intn(len(s))]
}

func main() {
	name := flag.String("file", "/Users/arunadeb/access.log", "Log file to monitor")
	flag.Parse()

	f, err := os.OpenFile(*name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
		panic("Unable to open the file")
	}
	defer f.Close()

	writer := NewFileWriter(f)
	for {
		time.Sleep(10)
		writer.BurstWriteForNSec(10, "/api/v1/users/arun", "GET", "200")
		time.Sleep(10)
		writer.BurstWriteForNSec(10, "/api/v1/users/arun", "GET", "200")
		time.Sleep(30)
		writer.BurstWriteForNSec(20, "/api/v1/users/arun", "GET", "200")
		time.Sleep(20)
		time.Sleep(120)
	}
}
