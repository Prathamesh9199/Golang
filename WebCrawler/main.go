package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/steelx/extractlinks"
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

	queue     = make(chan string)
	isVisited = make(map[string]bool)
)

func main() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		fmt.Print("Missing URL")
		os.Exit(1)
	}

	// Concurrency channel
	go func() {
		queue <- arguments[0]
	}()

	for href := range queue {
		if !isVisited[href] {
			crawlURL(href)
		}
	}
}

func fixURL(link, baseURL string) string {
	uri, err := url.Parse(link)
	if err != nil {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	// take basehost + uri path is uri does not have a host
	fixedURI := base.ResolveReference(uri)

	return fixedURI.String()
}

func isSameDomain(link, baseURL string) bool {
	uri, err := url.Parse(link)
	if err != nil {
		return false
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	if uri.Host != base.Host {
		return false
	}

	return true
}

func crawlURL(baseURL string) {
	isVisited[baseURL] = true

	fmt.Printf("Crawling -> %v\n", baseURL)

	response, err := netClient.Get(baseURL) // A Get request to URL
	checkError(err)                         // Check if any Error Exists
	defer response.Body.Close()

	links, err := extractlinks.All(response.Body) // Byte response
	checkError(err)
	for _, link := range links {
		if !isVisited[link.Href] && isSameDomain(link.Href, baseURL) {
			absoluteURL := fixURL(link.Href, baseURL)

			// Concurrency channel
			go func() {
				queue <- absoluteURL
			}()
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Print(queue)
		os.Exit(1)
	}
}
