package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	//Skip SSL
	config = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	netClient = &http.Client{
		Transport: transport,
	}
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func findTitle(pageContent string) string {
	var title string
	start := "<link rel=\"profile\" href=\"https://gmpg.org/xfn/11\">"
	titleStartIndex := strings.Index(pageContent, start)
	if titleStartIndex == -1 {
		title = "Not Found"
	} else {
		titleStartIndex += len(start) + 8

		titleEndIndex := strings.Index(pageContent, "<link rel='dns-prefetch' href='//fonts.googleapis.com' />")
		if titleEndIndex == -1 {
			title = "No Title Close Found"
		} else {
			titleEndIndex -= 31
			title = string([]byte(pageContent[titleStartIndex:titleEndIndex]))
		}
	}

	return title
}

func main() {
	baseURL := "https://mflixblog.xyz/archives/"

	/*
		Found and save at rawMovieLink.txt from 10000 - 11111
		Not Found from 11111 - 12379
	*/

	for number := 10000; number < 11111; number++ {
		response, err := netClient.Get(baseURL + fmt.Sprint(number)) // A Get request to URL
		checkError(err)                                              // Check if any Error Exists
		defer response.Body.Close()
		if response.StatusCode == 404 {
			fmt.Print("\n", response.Status, " for ", fmt.Sprint(number))
		} else {
			dataInBytes, err := io.ReadAll(response.Body)
			checkError(err) // Check if any Error Exists
			pageContent := string(dataInBytes)

			fmt.Print("\n", findTitle(pageContent)+" | Link => "+baseURL+fmt.Sprint(number))
		}
	}
}
