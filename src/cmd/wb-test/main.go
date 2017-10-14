package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const k = 5

func init() {
	log.SetFlags(0)
}

// Other funcs, if any.
func worker(ch <-chan string, totalCh chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range ch {
		cnt, err := getCount(url)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Count for %s: %d", url, cnt)

		totalCh <- cnt
	}
}

func getCount(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("response incorrect, status: %s\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return bytes.Count(body, []byte("Go")), nil
}

func summer(totalCh <-chan int, res chan<- int) {
	var sum int
	for total := range totalCh {
		sum += total
	}
	res <- sum
}

func main() {

	// Initialization, if any.
	var (
		workersNum int
		urlCh      = make(chan string)
		totalCh    = make(chan int)
		resCh      = make(chan int)
		wg         = &sync.WaitGroup{}
	)

	go summer(totalCh, resCh)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		// Do something with strings here.
		url := strings.TrimSpace(scanner.Text())
		if len(url) == 0 {
			continue
		}

		if workersNum < k {
			wg.Add(1)
			go worker(urlCh, totalCh, wg)
			workersNum++
		}
		urlCh <- url
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	close(urlCh)
	wg.Wait()
	close(totalCh)

	// Other code, if any.

	total := 0

	// Other code, if any.

	total = <-resCh

	log.Printf("Total: %v", total)
}
