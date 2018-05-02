/*
 * bt: Ben Thornburg
 */
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		print("Searches the Violation Comments for a given keyword.\nUsage: ./bt keyword\n")
		return
	}

	searchString := os.Args[1]
	violationContents := " " + searchString + " "
	baseURL := "https://data.cincinnati-oh.gov/resource/2c8u-zmu9.json?"
	limit := 10

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	completeURL := baseURL + "$select=business_name,violation_comments&$where=" +
		url.QueryEscape("violation_comments like '%"+
			violationContents+
			"%'") +
		"&$limit=" + fmt.Sprintf("%d", limit)

	req, error := http.NewRequest(http.MethodGet, completeURL, nil)

	if error != nil {
		fmt.Printf("Fatal error: %v\n", error)
		return
	}

	req = req.WithContext(ctx)

	res, error := http.DefaultClient.Do(req)
	if error != nil {
		log.Fatal("Oops: %v\n", error)
	}
	defer res.Body.Close()

	if res != nil && res.StatusCode == http.StatusOK {
		io.Copy(os.Stdout, res.Body)
	} else {
		fmt.Printf("Failure!\n")
	}
	return
}
