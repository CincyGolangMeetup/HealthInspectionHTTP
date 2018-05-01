/*
 * badams: Brian Adams
 */

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var tokenFlag string
var prettyPrintFlag bool

func init() {
	flag.StringVar(&tokenFlag, "t", "", "Flag used for the app token")
	flag.BoolVar(&prettyPrintFlag, "p", false, "Use for pretty printing")

	flag.Parse()
}

func main() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://data.cincinnati-oh.gov/resource/2c8u-zmu9.json?", nil)
	req.Header.Add("X-App-Token", tokenFlag)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error!")
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if prettyPrintFlag {
		var prettyJSON bytes.Buffer
		json.Indent(&prettyJSON, body, "", "\t")

		fmt.Println(string(prettyJSON.Bytes()))
	} else {
		fmt.Println(string(body))
	}

}
