package ppmapi

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/viper"
)

// API is data structure for our client
type API struct {
	Client          *http.Client
	URI             string
	NodeIP          string
	NodeName        string
	IntervalTypeKey string
	DurationSelect  string
	StartDate       string
	EndDate         string
	User            string
	Password        string
}

// GetCSV will invoke HTTP GET and return a byte slice
func (api *API) GetCSV(URL string) ([]byte, error) {

	// Set client object
	client := api.Client

	// Set the API request object and auth
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(api.User, api.Password)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// WriteCSV will process the response body into a csv file
func WriteCSV(path string, respBody []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Convert respBody into type that implements io.Reader
	r := bytes.NewReader(respBody)

	// Copy the byte
	io.Copy(f, r)
	log.Println("done writing CSV file")

	return nil
}

// URLBuilder will create an URL for the API call
func URLBuilder(key string, api API) (URL string, err error) {
	path, err := os.Executable()
	if err != nil {
		log.Fatalf("error reading path, %s", err)
	}
	dir := filepath.Dir(path)
	viper.SetConfigFile(dir + "/uri.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
		return "", err
	}
	fmt.Printf("using config: %s ; for node: %s ; resource is: %s\n", viper.ConfigFileUsed(), api.NodeName, key)

	// Get resource value
	tmpl := viper.Get(key + ".Resource")

	// Templating the URL
	t := template.New("URL template")
	t, err = t.Parse(tmpl.(string)) // assert t.interface{} into t.string
	if err != nil {
		log.Fatal("Parse: ", err)
		return "", err
	}
	// tpl is a buffer pointer that implement io.Writer
	// it is used to store the template execution result
	// then parse them to string
	var tpl bytes.Buffer
	if err = t.Execute(&tpl, api); err != nil {
		return "", err
	}
	URL = tpl.String()
	return URL, nil
}
