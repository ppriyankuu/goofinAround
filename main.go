package main

import (
	"fmt"
	"net/http"
	"sync"
)

var sites = []string{"test"}
var wg sync.WaitGroup
var mut sync.Mutex

func getStatusCode(endpoint string) {
	defer wg.Done()

	res, err := http.Get(endpoint)

	if err != nil {
		fmt.Println("something went wrong!!!")
	} else {
		mut.Lock()
		sites = append(sites, endpoint)
		mut.Unlock()
		fmt.Printf("%d status code for %s\n", res.StatusCode, endpoint)
	}
}

func main() {
	siteList := []string{
		"https://100xdevs.com",
		"https://go.dev",
		"https://google.com",
		"https://x.com",
		"https://facebook.com",
	}

	for _, web := range siteList {
		go getStatusCode(web)
		wg.Add(1)
	}

	wg.Wait()
	fmt.Println(sites)
}
