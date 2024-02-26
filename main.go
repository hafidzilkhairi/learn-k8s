package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 200,
			MaxIdleConns:        200,
		},
	}
	totalRequest, totalError, totalSuccess int
	totalDuration, maxDuration             time.Duration
	isStop                                 bool
	wg                                     = &sync.WaitGroup{}
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	now := time.Now()
	for i := 0; i < 100; i++ {
		go hit()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	totalRuntime := time.Since(now)

	isStop = true
	wg.Wait()

	RPS := float64(totalRequest)

	if totalRuntime > 1 {
		RPS = RPS / totalRuntime.Seconds()
	}

	fmt.Println("\nTotal Request		", totalRequest)
	fmt.Println("RPS			", RPS)
	fmt.Println("Total Success		", totalSuccess)
	fmt.Println("Success Rate		", float64(totalSuccess)/float64(totalRequest)*100)
	fmt.Println("Total duration		", totalRuntime)
	fmt.Println("Max duration		", maxDuration)
	fmt.Println("Average duration	", totalDuration/time.Duration(totalRequest))
}

func hit() {
	req, err := http.NewRequest(http.MethodGet, "http://192.168.49.2:31234", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	wg.Add(1)

	for {
		if isStop {
			wg.Done()
			break
		}
		now := time.Now()
		resp, err := client.Do(req)
		currentTotalDuration := time.Since(now)
		totalDuration += currentTotalDuration
		if currentTotalDuration > maxDuration {
			maxDuration = currentTotalDuration
		}
		switch {
		case err != nil:
			log.Println(err.Error())
			totalError++
		case resp.StatusCode != http.StatusOK:
			log.Println("Status not OK", resp.StatusCode)
			totalError++
		default:
			_, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			// log.Println("Response:", string(body))
			totalSuccess++
		}
		totalRequest++
	}
}
