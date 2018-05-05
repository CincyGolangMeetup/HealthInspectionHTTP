package inspect_test

import (
	"testing"

	inspect "github.com/geneticgrabbag/HealthInspectionHTTP/geneticgrabbag"
)

func TestInspection_Summarize(t *testing.T) {
	t.Run("All values provided", func(t *testing.T) {
		i := inspect.Inspection{
			BusinessName: "Acme, LLC",
			ActionDate:   "2010-09-08T00:00:00.000",
			ActionStatus: "harsh words",
		}
		want := `On 2010-09-08T00:00:00.000, Acme, LLC had an inspection resulting in "harsh words"`
		got := i.Summarize()
		if got != want {
			t.Errorf("Inspection.Summarize() = %v, want %v", got, want)
		}
	})
}

func TestInspection_MapURL(t *testing.T) {
	t.Run("All values provided", func(t *testing.T) {
		i := inspect.Inspection{
			Latitude:  "25.5",
			Longitude: "-82.1",
		}
		want := `http://www.openstreetmap.org/?mlat=25.5&mlon=-82.1#map=17/25.5/-82.1`
		got := i.MapURL()
		if got != want {
			t.Errorf("Inspection.MapURL() = %v, want %v", got, want)
		}
	})

}
