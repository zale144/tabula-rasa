package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {
	urls := []string{"http://python.org", "http://golang.org"}
	responses := make(chan string)

	for _, url := range urls {
		go func(url string) {
			resp, _ := http.Get(url)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			responses <- string(body)
		}(url)
	}

	for response := range responses {
		fmt.Println(response)
	}
}
