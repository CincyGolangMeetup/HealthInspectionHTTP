package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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

type ApiError struct {
	Code    string `json:"code"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    struct {
		Query string `json:"query"`
	} `json:"data"`
}

func main() {
	req, err := http.NewRequest("GET", "https://data.cincinnati-oh.gov/resource/2c8u-zmu9.json", nil)
	if err != nil {
		fmt.Printf("There was an error creating the request: %s", err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	req.WithContext(ctx)

	query := req.URL.Query()
	query.Add("$limit", "20")
	query.Add("license_status", "'PAID'")
	query.Add("postal_code", "45202")
	query.Add("$where", "city NOT 'Cincinnati'")
	req.URL.RawQuery = query.Encode()

	fmt.Printf("The URL is %s.\n\n", req.URL)

	req.Header.Add("X-App-Token", os.Getenv("API_TOKEN"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("There was an error making the request: %s", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errorResponse ApiError
		if err := decoder.Decode(&errorResponse); err != nil {
			fmt.Printf("There was an error decoding the errored request: %s\n", err.Error())
		} else {
			fmt.Printf("The API returned an error - %s\n", errorResponse.Message)
			fmt.Printf("The SoQL query was '%s'\n", errorResponse.Data.Query)
		}

		os.Exit(1)
	}

	var data DataSet
	if err := decoder.Decode(&data); err != nil {
		fmt.Printf("There was an error decoding the request: %s\n", err.Error())
		os.Exit(1)
	}

	for i, d := range data {
		fmt.Printf("%d: %s - %s\n", i+1, d.BusinessName, d.ViolationDescription)
	}
}
