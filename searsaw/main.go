package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DataSet []struct {
	ActionDate           string `json:"action_date"`
	ActionSequence       string `json:"action_sequence"`
	ActionStatus         string `json:"action_status"`
	Address              string `json:"address"`
	BusinessName         string `json:"business_name"`
	City                 string `json:"city"`
	Code                 string `json:"code"`
	InspSubtype          string `json:"insp_subtype"`
	InspType             string `json:"insp_type"`
	LastTableUpdate      string `json:"last_table_update"`
	Latitude             string `json:"latitude"`
	LicenseNo            string `json:"license_no"`
	LicenseStatus        string `json:"license_status"`
	Longitude            string `json:"longitude"`
	PhoneNumber          string `json:"phone_number"`
	PostalCode           string `json:"postal_code"`
	RecordnumInsp        string `json:"recordnum_insp"`
	RecordnumLicense     string `json:"recordnum_license"`
	State                string `json:"state"`
	ViolationComments    string `json:"violation_comments"`
	ViolationDescription string `json:"violation_description"`
	ViolationKey         string `json:"violation_key"`
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://data.cincinnati-oh.gov/resource/2c8u-zmu9.json", nil)
	if err != nil {
		fmt.Printf("There was an error creating the request: %s", err.Error())
		os.Exit(1)
	}

	query := req.URL.Query()
	query.Add("$limit", "20")
	query.Add("license_status", "'PAID'")
	query.Add("postal_code", "45202")
	req.URL.RawQuery = query.Encode()

	fmt.Printf("The URL is %s.\n\n", req.URL)

	req.Header.Add("X-App-Token", os.Getenv("API_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("There was an error making the request: %s", err.Error())
		os.Exit(1)
	}

	var data DataSet
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Printf("There was an error decoding the request: %s\n", err.Error())
		os.Exit(1)
	}

	for i, d := range data {
		fmt.Printf("%d: %s - %s\n", i+1, d.BusinessName, d.ViolationDescription)
	}
}
