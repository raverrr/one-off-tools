package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

const threshold = 3.0

var wg sync.WaitGroup
var (
	entropyThreshold = flag.Float64("e", 3.0, "Maximum entropy level")
)

func entropy(s string) float64 {
	m := make(map[rune]int)
	for _, r := range s {
		m[r]++
	}
	var e float64
	for _, c := range m {
		p := float64(c) / float64(len(s))
		e += p * math.Log2(p)
	}
	return -e
}

func isValidWord(word string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
	return re.MatchString(word) && !strings.Contains(word, ".")
}

func main() {
	flag.Parse()
	// Use a buffered reader to read from stdin
	reader := bufio.NewReader(os.Stdin)

	// Create a concurrent map to store the directory wordlist
	var wordlist sync.Map

	// Use a goroutine to read each line from stdin
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// Read a line from stdin
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}

			// Parse the line as a URL
			u, err := url.Parse(string(line))
			if err != nil {
				continue
			}

			// Extract the path from the URL
			path := u.Path

			// Split the path into individual words
			words := strings.Split(path, "/")

			// Add each word to the wordlist
			for _, word := range words {
				if isValidWord(word) && entropy(word) <= threshold {
					wordlist.Store(word, true)
				}
			}
		}
	}()
	wg.Wait()
	// Print the wordlist
	wordlist.Range(func(word, _ interface{}) bool {
		fmt.Println(word)
		return true
	})
}
