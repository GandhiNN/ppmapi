package ppmapi

import (
	"encoding/csv"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestAPI(t *testing.T) {

	// assertEqual is a helper function to test file equality with reference
	assertEqual := func(t *testing.T, bytein, byteout interface{}) {
		t.Helper()
		if !reflect.DeepEqual(bytein, byteout) {
			t.Fatal("not getting csv file")
		}
	}

	// Test API call
	t.Run("Starts local server and invoke API call", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			record := [][]string{{"line1", "test1"}, {"line2", "test2"}}
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", "attachment;filename=test.csv")
			wr := csv.NewWriter(w)
			defer wr.Flush()

			// Write the csv file to response body
			for _, val := range record {
				err := wr.Write(val)
				if err != nil {
					t.Errorf("Error sending csv: %s %v", err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}))
		defer server.Close()

		// Get the resource
		URL := server.URL
		api := API{
			Client:          server.Client(),
			URI:             "",
			NodeIP:          "",
			NodeName:        "",
			IntervalTypeKey: "",
			DurationSelect:  "",
			StartDate:       "",
			EndDate:         "",
			User:            "",
			Password:        "",
		}
		csvOut, err := api.GetCSV(URL)
		if err != nil {
			t.Errorf("API call cannot be invoked: %s", err.Error())
		}

		// Load the goldenfile into byteslice
		golden, err := ioutil.ReadFile("../testdata/golden.csv")
		if err != nil {
			t.Errorf("golden file not found: %s", err.Error())
		}

		// Do the test
		assertEqual(t, csvOut, golden)
	})
}

func TestWriteCSV(t *testing.T) {

	// Test API write to CSV
	t.Run("Write to CSV", func(t *testing.T) {
		record := [][]string{{"line3", "test3"}, {"line4", "test4"}}
		f, err := os.Create("../testdata/silver.csv")
		if err != nil {
			t.Errorf("cannot open silver file: %s", err.Error())
		}

		// Write records to file
		w := csv.NewWriter(f)
		err = w.WriteAll(record)
		if err != nil {
			t.Errorf("cannot write to file: %s", err.Error())
		}
	})
}
