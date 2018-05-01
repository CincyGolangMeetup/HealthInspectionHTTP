/*
 * bgardner87: Brad Gardner
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//Inspection Represents a health inspection
type Inspection struct {
	ActionDate           string `json:"action_date"`
	Address              string `json:"address"`
	ViolationKey         string `json:"violation_key"`
	ViolationDescription string `json:"violation_description"`
	ViolationComments    string `json:"violation_comments"`
}

func main() {
	companyName := os.Args[1]

	fmt.Println(companyName)

	baseURL := "https://data.cincinnati-oh.gov/resource/2c8u-zmu9.json?"

	url := baseURL + "$where=" + url.QueryEscape("business_name like '%"+strings.ToUpper(companyName)+"%'") + "&$limit=10"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Fail!")
	}

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	var inspections []Inspection

	json.Unmarshal(body, &inspections)

	fmt.Printf("Inspections : %+v", inspections)
}
