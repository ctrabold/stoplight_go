package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	_ "os"
	"regexp"
	_ "regexp"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

// Filter is a filter function applied to a single record.
type Filter func(Project) bool

// FilterBulk is a bulk filter function applied to an entire slice of records.
type FilterBulk func([]Project) []Project

// Projects XML & JSON structure

// LastBuildStatus <implement comment>
type LastBuildStatus string

// CurrentStatus <implement comment>
type CurrentStatus string

// Projects <implement comment>
type Projects struct {
	XMLName  xml.Name  `xml:"Projects"`
	Projects []Project `xml:"Project"`
}

// Project describes the XML structure of a project
type Project struct {
	XMLName         xml.Name        `xml:"Project"`
	Name            string          `json:"name"              xml:"name,attr"`
	BuildURL        string          `json:"build_url"         xml:"webUrl,attr"`
	LastBuildID     string          `json:"last_build_id"     xml:"lastBuildLabel,attr"`
	LastBuildTime   string          `json:"last_build_time"   xml:"lastBuildTime,attr"`
	LastBuildStatus LastBuildStatus `json:"last_build_status" xml:"lastBuildStatus,attr"`
	CurrentStatus   CurrentStatus   `json:"current_status"    xml:"activity,attr"`
}

// UnmarshalXMLAttr converts the XML value to JSON value
func (s *LastBuildStatus) UnmarshalXMLAttr(attr xml.Attr) error {
	status := attr.Value
	var attrValue string

	switch strings.ToLower(status) {
	default:
		attrValue = "unknown"
	case "failure":
		attrValue = "failed"
	case "success":
		attrValue = "passed"
	}

	*s = LastBuildStatus(attrValue)
	return nil
}

// UnmarshalXMLAttr converts the XML value to JSON value
func (s *CurrentStatus) UnmarshalXMLAttr(attr xml.Attr) error {
	status := attr.Value
	var attrValue string

	switch strings.ToLower(status) {
	default:
		attrValue = "unknown"
	case "building":
		attrValue = "building"
	case "sleeping":
		attrValue = "done"
	}

	*s = CurrentStatus(attrValue)
	return nil
}

func config() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	viper.SetDefault("Port", 1323)
	viper.SetDefault("Host", "0.0.0.0")

	// Show all projects by default
	viper.SetDefault("server.projects", ".*")
	viper.SetDefault("server.ignored_projects", "")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
	// ENV Variables Setup
	// Prefix with "STOPLIGHT_"
	viper.SetEnvPrefix("stoplight")
	viper.AutomaticEnv()
}

func main() {
	config()

	// HTTP
	host := viper.GetString("Host")
	port := viper.GetInt("Port")
	address := fmt.Sprintf("%s:%v", host, port)

	fmt.Printf("Got server address %v\n", address)

	s := &http.Server{
		Addr: address,
	}

	e := echo.New()
	e.HideBanner = true

	// the file server for rice. "public" is the folder where the files come from.
	assetHandler := http.FileServer(rice.MustFindBox("public").HTTPBox())
	// serves the index.html from rice
	e.GET("/", echo.WrapHandler(assetHandler))
	// servers other static files
	e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", assetHandler)))

	e.GET("/projects.json", FetchProjects)
	e.Logger.Fatal(e.StartServer(s))
}

// FetchProjects e.GET("/projects.json", FetchProjects)
// CONNECT -> FETCH CC XML -> PARSE -> FILTER -> JSON
func FetchProjects(c echo.Context) error {
	server := viper.GetString("server.url")

	if !viper.IsSet("server") {
		log.Fatal("missing a server to connect to. Aborting...")
	}

	// Parse URL
	uri, err := url.Parse(server)
	if err != nil {
		log.Fatal("Failed on parsing URL", err)
	}

	// Fetch XML from CI server
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", uri.String(), nil)
	if viper.IsSet("server.username") && viper.IsSet("server.password") {
		username := viper.GetString("server.username")
		password := viper.GetString("server.password")
		req.SetBasicAuth(username, password)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed on client", err)
	}

	ccxml, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal("Failed on parsing Body", err)
	}

	// Parse XML
	data := &Projects{}
	err = xml.Unmarshal(ccxml, data)
	if err != nil {
		// TODO ct 2019-04-23 Handle empty responses gracefully
		// eg. so project.json shows an empty JSON object
		log.Println("Failed on XML Unmarshal:", err, "Please check your XML")
	}

	// fmt.Println("DEBUG data")
	// fmt.Printf("DEBUG %T", data.Project)

	filteredProjects := FilterProjects(data.Projects)

	// Return JSON
	return c.JSON(http.StatusOK, filteredProjects)
}

// FilterProjects applies a set of filters removing any unwanted projects
func FilterProjects(records []Project) []Project {
	return ApplyBulkFilters(
		ApplyFilters(records,
			FilterIncludeProjects,
			FilterExcludedProjects,
		),
		FilterDuplicates,
	)
}

// ApplyFilters applies a set of filters to a record list.
// Each record will be checked against each filter.
// The filters are applied in the order they are passed in.
func ApplyFilters(records []Project, filters ...Filter) []Project {
	// Make sure there are actually filters to be applied.
	if len(filters) == 0 {
		return records
	}

	filteredRecords := make([]Project, 0, len(records))

	// Range over the records and apply all the filters to each record.
	// If the record passes all the filters, add it to the final slice.
	for _, r := range records {
		keep := true

		for _, f := range filters {
			if !f(r) {
				keep = false
				break
			}
		}

		if keep {
			filteredRecords = append(filteredRecords, r)
		}
	}

	return filteredRecords
}

// ApplyBulkFilters applies a set of filters to the entire slice of records.
// Used when each record filter requires knowledge of the other records, e.g. de-duping.
func ApplyBulkFilters(records []Project, filters ...FilterBulk) []Project {
	for _, f := range filters {
		records = f(records)
	}

	return records
}

// FilterDuplicates is a bulk filter to remove any duplicates from the set.
func FilterDuplicates(records []Project) []Project {
	recordMap := map[string]bool{}
	filteredRecords := []Project{}

	for _, record := range records {
		if ok := recordMap[record.Name]; ok {
			continue
		}
		recordMap[record.Name] = true
		filteredRecords = append(filteredRecords, record)
	}

	return filteredRecords
}

// FilterIncludeProjects makes sure *only* these projects are listed.
func FilterIncludeProjects(record Project) bool {
	patterns := viper.GetStringSlice("server.projects")
	// Compile regular expressions just once
	regexps := make([]*regexp.Regexp, len(patterns))
	for i := range patterns {
		regexps[i] = regexp.MustCompile(patterns[i])
	}

	for _, regexp := range regexps {
		if regexp.MatchString(record.Name) {
			fmt.Printf("allowing element %s because it matches allow pattern\n", record.Name)
			return true
		}
	}

	return false
}

// FilterExcludedProjects filters out all excluded projects.
func FilterExcludedProjects(record Project) bool {
	patterns := viper.GetStringSlice("server.ignored_projects")
	regexps := make([]*regexp.Regexp, len(patterns))
	for i := range patterns {
		regexps[i] = regexp.MustCompile(patterns[i])
	}

	for _, regexp := range regexps {
		if regexp.MatchString(record.Name) {
			fmt.Printf("deleting element %s due to ignore pattern\n", record.Name)
			return false
		}
	}

	return true
}
