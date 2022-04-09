package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	client *http.Client
	wg     sync.WaitGroup
)

var feed int
var timeout int

func main() {

	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "Concurrency")
	flag.IntVar(&feed, "f", 50, " delay the ingestion of URLs by X milliseconds (May help when scanning large lists on low ram)")
	flag.IntVar(&timeout, "t", 5, " Timout (seconds)")

	flag.Parse()

	var input io.Reader
	input = os.Stdin

	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Printf("failed to open file: %s\n", err)
			os.Exit(1)
		}
		input = file
	}

	sc := bufio.NewScanner(input)

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        concurrency,
			MaxIdleConnsPerHost: concurrency,
			MaxConnsPerHost:     concurrency,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 5 * time.Second,
	}

	semaphore := make(chan bool, concurrency)

	for sc.Scan() {
		raw := sc.Text()
		time.Sleep(time.Duration(feed) * time.Millisecond)
		wg.Add(1)
		semaphore <- true
		go func(raw string) {
			defer wg.Done()
			u, err := url.ParseRequestURI(raw)
			if err != nil {
				return
			}
			fetchURL(u)

		}(raw)
		<-semaphore
	}

	wg.Wait()

	if sc.Err() != nil {
		fmt.Printf("error: %s\n", sc.Err())
	}
}

func fetchURL(u *url.URL) (*http.Response, error) {

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()
	fmt.Println(finalURL)

	io.Copy(ioutil.Discard, resp.Body)

	return resp, err

}
