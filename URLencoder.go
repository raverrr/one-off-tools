package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	var output []string
	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, strings.TrimSpace(input))
	}
	for j := 0; j < len(output); j++ {
		fmt.Println(url.QueryEscape(output[j]))
	}
}
