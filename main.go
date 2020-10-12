package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	helpPtr := flag.Bool("help", false, "Lists possible arguments")
	urlPtr := flag.String("url", "", "URL to perform GET HTML request (ex: http://example.com/)")
	profilePtr := flag.Int("profile", 0, "Number of times to request from the URL")

	flag.Parse()
	if *helpPtr == true || *urlPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	addr, err := url.Parse(*urlPtr)
	if err != nil {
		log.Fatal(err)
	}
	path := addr.EscapedPath()
	if path == "" {
		path = "/"
	}
	if *profilePtr == 0 {
		*profilePtr = 1
	}

	var success int = *profilePtr
	var errorCodes = make([]string, *profilePtr)
	var minTime int64 = 0
	var maxTime int64 = 0
	var times = make([]int, *profilePtr)
	var minSize int = 0
	var maxSize int = 0

	for i := 0; i < *profilePtr; i++ {
		fmt.Printf("\nRequest #%d:\n", i)
		startTime := time.Now()
		conn, err := net.Dial("tcp", addr.Hostname()+":80")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(conn, "GET "+path+" HTTP/1.1\r\nHost: "+addr.Hostname()+"\r\nConnection: close\r\n\r\n")
		scanner := bufio.NewScanner(conn)

		scanner.Scan()
		endTime := time.Now()
		statusCode := scanner.Text()
		if !strings.Contains(statusCode, "200 OK") {
			errorCode := string(statusCode[0:len(statusCode)])
			fmt.Printf("Server responded with error code: %s\n", errorCode)
			errorCodes[i] = errorCode
			success--
		} else {
			fmt.Printf("Successfully connected to %s\n", addr)
		}

		currTime := endTime.Sub(startTime).Milliseconds()
		times[i] = int(currTime)
		if currTime > maxTime {
			maxTime = currTime
		}
		if currTime < minTime || i == 0 {
			minTime = currTime
		}
		fmt.Printf("Response time: %dms\n", currTime)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Content-Length:") {
				currSize, err := strconv.Atoi(string(line[16:len(line)]))
				if err != nil {
					log.Fatal(err)
				}
				if currSize > maxSize {
					maxSize = currSize
				}

				if currSize < minSize || i == 0 {
					minSize = currSize
				}
				fmt.Printf("Response size: %d bytes\n", currSize)
			} else if line == "" {
				break
			}
		}
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		conn.Close()
	}

	sort.Ints(times)
	median := times[len(times)/2]
	if len(times)%2 == 0 {
		median = median + times[len(times)/2-1]/2
	}
	sum := 0
	for _, i := range times {
		sum += int(i)
	}

	fmt.Println("\nSummary:")
	fmt.Printf("Number of requests: %d\n", *profilePtr)
	fmt.Printf("Fastest time: %d ms\n", minTime)
	fmt.Printf("Slowest time: %d ms\n", maxTime)
	fmt.Printf("Mean time: %d ms\n", sum/len(times))
	fmt.Printf("Median time: %d ms\n", median)
	fmt.Printf("Percentage of successful requests: %f%%\n", float64(success/(*profilePtr))*100)
	fmt.Print("Error Codes: ")
	for i := 0; i < len(errorCodes); i++ {
		if errorCodes[i] != "" {
			fmt.Printf("\n\tRequest #%d: %s", i, errorCodes[i])
		}
	}
	fmt.Printf("\nSmallest response size: %d bytes\n", minSize)
	fmt.Printf("Largest response size: %d bytes\n", maxSize)

	os.Exit(0)
}
