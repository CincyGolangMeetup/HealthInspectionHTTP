// Package cincyfsp provides an interace to data provided by the Cincinnati
// Food Safety Program.
package cincyfsp

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	inspect "github.com/geneticgrabbag/HealthInspectionHTTP/geneticgrabbag"
)

const (
	apiURL = "https://data.cincinnati-oh.gov/resource/2c8u-zmu9.json"

	defaultTimeout = 30 * time.Second
	defaultLimit   = 10
)

var (
	// Ensure our repository implements health inspection repository.
	_ inspect.InspectionRepository = &InspectionRepository{}

	// Custom errors.
	errAPIFailure    = errors.New("API did not return expected 200 HTTP OK")
	errNegativeLimit = errors.New("limit must not be negative")
)

// InspectionRepository provides access to health inspection records via the
// Cincinnati Food Service API.
type InspectionRepository struct {
	name   string
	client *http.Client
	limit  int
	token  string
}

var defaultRepository = InspectionRepository{
	name:   "Default Repository",
	client: &http.Client{Timeout: defaultTimeout},
	limit:  defaultLimit,
}

// Name of this repository instance.
func (s *InspectionRepository) Name() string {
	return s.name
}

// GetAll returns all known health inspections.
func (s *InspectionRepository) GetAll() (inspect.Inspections, error) {

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	if s.token != "" {
		log.Println("Using supplied API token.")
		req.Header.Add("X-API-Token", s.token)
	}

	q := req.URL.Query()
	if s.limit > 0 {
		q.Add("$limit", strconv.Itoa(s.limit))
	}
	req.URL.RawQuery = q.Encode()

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errAPIFailure
	}

	var aa inspect.Inspections

	d := json.NewDecoder(res.Body)
	for err := d.Decode(&aa); err != nil && err != io.EOF; {
		return nil, err
	}

	return aa, nil
}

// InspectionRepositoryOption is a function that takes an inspection
// repository and augments it.  Uses the functional option pattern as
// popularized by Dave Cheney.
type InspectionRepositoryOption func(*InspectionRepository) error

// WithName sets the name of the repository instance.
func WithName(name string) InspectionRepositoryOption {
	return func(s *InspectionRepository) error {
		s.name = name
		return nil
	}
}

// WithToken sets the token used to authenticate against the API.
func WithToken(token string) InspectionRepositoryOption {
	return func(s *InspectionRepository) error {
		s.token = token
		return nil
	}
}

// WithTimeout sets the maximum time a request can take before giving up.
func WithTimeout(timeout time.Duration) InspectionRepositoryOption {
	return func(s *InspectionRepository) error {
		s.client.Timeout = timeout
		return nil
	}
}

// WithLimit sets the limit on the number of inspections that will be
// returned.
func WithLimit(limit int) InspectionRepositoryOption {
	return func(s *InspectionRepository) error {
		if limit < 0 {
			return errNegativeLimit
		}
		s.limit = limit
		return nil
	}
}

// WithClient sets the HTTP client used for transport.
func WithClient(client *http.Client) InspectionRepositoryOption {
	return func(s *InspectionRepository) error {
		s.client = client
		return nil
	}
}

// WithExampleData replaces the HTTP client with a mock transport that returns
// example data.
func WithExampleData() InspectionRepositoryOption {
	return WithClient(&http.Client{Transport: newMockTransport()})
}

// NewRepository returns a new inspection repository with the default HTTP
// client.
func NewRepository(opts ...InspectionRepositoryOption) (*InspectionRepository, error) {
	repo := defaultRepository
	for _, opt := range opts {
		err := opt(&repo)
		if err != nil {
			return nil, err
		}
	}

	return &repo, nil
}

type mockTransport struct{}

func newMockTransport() http.RoundTripper {
	return &mockTransport{}
}

func (t *mockTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {

	resBody := `
  [{
    "action_date": "2013-10-10T00:00:00.000",
    "action_sequence": "A012642688",
    "action_status": "Not Abated",
    "address": "1801 RACE ST",
    "business_name": "\"KROEGER & SONS MEATS\"",
    "city": "CINCINNATI",
    "code": "3717-1-03.4(C)",
    "insp_subtype": "STANDARD INSPECTION",
    "insp_type": "ROUTINE",
    "last_table_update": "2015-03-05T21:56:43.000",
    "latitude": "39.1154039286902",
    "license_no": "RFE-008760-C3S",
    "license_status": "PAID",
    "longitude": "-84.5184463268576",
    "phone_number": "513 651-5543",
    "postal_code": "45202",
    "recordnum_insp": "CFSI131698",
    "recordnum_license": "H200801025",
    "state": "OH",
    "violation_comments": "\"A BOX OF FROZEN PORK PRODUCT WAS LEFT ON THE FLOOR IN THE CUSTOMER SERVICE AREA BY A DELIVERY DRIVER.  UPON BRINGING THIS TO THE EMPLOYEES ATTENTION, THE PRODUCT WAS MOVED TO THE WALK IN COOLER.\"",
    "violation_description": "\"3717-1-03.4(C)  - Violation - Thawing TCS - Temperature & Time ControlTCS food was improperly thawed.\"",
    "violation_key": "HLE130775V"
	}]`

	res = &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	res.Body = ioutil.NopCloser(strings.NewReader(resBody))

	return res, nil
}
