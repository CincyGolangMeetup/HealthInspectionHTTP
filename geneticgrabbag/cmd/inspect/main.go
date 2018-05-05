package main

import (
	"flag"
	"log"
	"os"

	"github.com/geneticgrabbag/HealthInspectionHTTP/geneticgrabbag"
	fsp "github.com/geneticgrabbag/HealthInspectionHTTP/geneticgrabbag/cincyfsp"
)

func main() {

	// Grab command line arguments.
	dryRunOpt := flag.Bool("dryrun", false, "perform a dry run with example data (no network activity)")
	limitOpt := flag.Int("limit", 0, "maximum number of inspections to return")
	tokenOpt := flag.String("token", os.Getenv("API_TOKEN"), "token to authenticate against API")
	flag.Parse()

	// Create the Cincy Foods inspector repository implementation. For a dry
	// run, create a repository that returns example data without a network
	// call.  Otherwise, create a repository that calls the real API.
	var repo inspect.InspectionRepository
	var err error
	if *dryRunOpt {
		repo, err = fsp.NewRepository(
			fsp.WithName("Example Data"),
			fsp.WithExampleData(),
			fsp.WithLimit(*limitOpt))
	} else {
		repo, err = fsp.NewRepository(
			fsp.WithName("Real API Data"),
			fsp.WithToken(*tokenOpt),
			fsp.WithLimit(*limitOpt))
	}
	if err != nil {
		log.Fatalf("Unable to create repository: %s\n", err)
	}

	// Get the inspections from the repository
	ii, err := repo.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	// And dump each inspection.
	log.Printf("%s - %d inspection(s):\n", repo.Name(), len(ii))
	for n, i := range ii {
		log.Println("  - Inspection", n+1)
		log.Println("    - Summary of inspection:", i.Summarize())
		log.Println("    - Map of inspection site:", i.MapURL())
	}
}
