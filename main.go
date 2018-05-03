package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Inspection struct {
	ActionDate           string `json:"action_date"`
	ActionSequence       string `json:"action_sequence"`
	ActionStatus         string `json:"action_status"`
	Address              string `json:"address"`
	BusinessName         string `json:"business_name"`
	City                 string `json:"city"`
	Code                 string `json:"code"`
	InspSubType          string `json:"insp_subtype"`
	InspType             string `json:"insp_type"`
	LastTableUpdate      string `json:"last_table_update"` // TODO this should be time.Time but unmarshal error...
	Latitude             string `json:"latitude"`
	LicenseNo            string `json:"license_no"`
	LicenseStatus        string `json:"license_status"`
	Longitude            string `json:"longitude"`
	PhoneNumber          string `json:"phone_number"`
	PostalCode           string `json:"postal_code"`
	RecordNumInsp        string `json:"recordnum_insp"`
	RecordNumLicense     string `json:"recordnum_license"`
	State                string `json:"state"`
	ViolationComments    string `json:"violation_comments"`
	ViolationDescription string `json:"violation_description"`
	ViolationKey         string `json:"violation_key"`
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

const (
	BusinessName = "Kroger"
	BaseUrl      = "https://data.cincinnati-oh.gov/resource/2c8u-zmu9.json?"
	ResultLimit  = 10
)

func main() {
	// use a channel and goroutine because why not...
	var inspChan = make(chan []Inspection, 1)

	// TODO: not sure if a WaitGroup has value here...
	//var wg sync.WaitGroup
	//wg.Add(1)
	go GetInspections(inspChan)
	//wg.Wait()

	var results []Inspection
	if r, ok := <-inspChan; ok {
		results = r
	}

	fmt.Println("Total count is", len(results))
	for i := 0; i < len(results); i++ {
		fmt.Printf("Business Name: %v, Inspection Type: %v, Date: %v, Action Status: %v \n", results[i].BusinessName, results[i].InspType, results[i].LastTableUpdate, results[i].ActionStatus)
	}
}

func GetInspections(ch chan<- []Inspection) {
	// make sure we defer so we signal we're done
	//defer wg.Done()
	defer close(ch)

	requestUrl := BaseUrl + "$where=" +
		url.QueryEscape("business_name like '%"+strings.ToUpper(BusinessName)+"%'") +
		"&$limit=" + fmt.Sprintf("%d", ResultLimit)

	resp, err := httpClient.Get(requestUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	if resp == nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Failure! Status: %v", resp.StatusCode)
		return
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	results := make([]Inspection, ResultLimit)
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		fmt.Println(err.Error())
		return
	}

	ch <- results
}
