package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func readFromChannel(input chan string, output chan string) {
	var url string
	url = <-input
	output <- url
}
func checkAndSaveBody(output chan string) {
	var url string
	url = <-output
	resp, err := http.Get(url)

	if err != nil {

		s := fmt.Sprintf("%s is DOWN!\n", url)
		s += fmt.Sprintf("Error: %v  \n", err)

		// Sending the string over the channel.
		// This is a blocking call, so this goroutine will
		// wait for the main goroutine to receive it on the other part of the channel.
		output <- s
	} else {

		s := fmt.Sprintf("Status Code: %d  \n", resp.StatusCode)

		if resp.StatusCode == 200 {
			bodyBytes, err := ioutil.ReadAll(resp.Body)

			file := strings.Split(url, "//")[1]
			file += ".txt"

			s += fmt.Sprintf("Writing response Body to %s\n", file)

			err = ioutil.WriteFile(file, bodyBytes, 0664)
			if err != nil {

				s += "Error writing to file!\n"

				// sending s over the channel
				output <- s
			}
		}
		s += fmt.Sprintf("%s is UP\n", url)

		// sending s over the channel
		output <- s
	}
}

func main() {
	urls := []string{"https://www.golang.org", "https://www.google.com", "https://www.youtube.com", "https://www.facebook.com"}

	// Declaring a new channel
	input := make(chan string)
	output := make(chan string)

	go readFromChannel(input, output)

	for _, url := range urls {
		input <- url
	}

	for elem := range output {
		checkAndSaveBody(elem)
	}
}
