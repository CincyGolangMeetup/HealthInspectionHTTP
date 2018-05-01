/*
 * whh8b: Will Hawkins
 */
package main

import (
	"net/http"
	"net/url"
	"context"
	"fmt"
	"io"
	"os"
	"time"
	"log"
	"strings"
)

func main() {
	business_name := "Kroger"
	base_url := "https://data.cincinnati-oh.gov/resource/2c8u-zmu9.json?"
	limit := 2

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	complete_url := base_url + "$where=" +
		url.QueryEscape("business_name like '%" +
			strings.ToUpper(business_name) +
			"%'") +
		"&$limit=" + fmt.Sprintf("%d", limit)

	req, error := http.NewRequest(http.MethodGet, complete_url, nil)

	if (error != nil) {
		fmt.Printf("Fatal error: %v\n", error)
		return
	}

	req = req.WithContext(ctx)

	res, error := http.DefaultClient.Do(req)
	if (error != nil) {
		log.Fatal("Oops: %v\n", error)
	}
	defer res.Body.Close()

	if (res != nil && res.StatusCode == http.StatusOK) {
		io.Copy(os.Stdout, res.Body)
	} else {
		fmt.Printf("Failure!\n");
	}
	return
}
