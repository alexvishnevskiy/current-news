package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type response struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Results      []struct {
		Title       string      `json:"title"`
		Link        string      `json:"link"`
		Keywords    []string    `json:"keywords"`
		Creator     []string    `json:"creator"`
		VideoUrl    interface{} `json:"video_url"`
		Description string      `json:"description"`
		Content     interface{} `json:"content"`
		PubDate     string      `json:"pubDate"`
		ImageUrl    interface{} `json:"image_url"`
		SourceId    string      `json:"source_id"`
		Country     []string    `json:"country"`
		Category    []string    `json:"category"`
		Language    string      `json:"language"`
	} `json:"results"`
}

func Fetch(url string, countries []string, categories []string) (map[string]int, error) {
	apiKey := os.Getenv("API_KEY")
	categoryDict := make(map[string]int) //category->#articles

	// iterate over each continent and category
	for _, category := range categories {
		url = getUrl(url, countries, category, apiKey)
		// get data from news api
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		// read data
		resp, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		// parse response
		var responseObject response
		json.Unmarshal(resp, &responseObject)
		categoryDict[category] = responseObject.TotalResults
	}
	return categoryDict, nil
}

func FetchAllData() map[string]map[string]int {
	var (
		c          Config
		categories []string
		countries  []string
		continent  string
	)
	c.GetConf()

	// dict to store all data
	continentsDict := make(map[string]map[string]int)
	// make struct iterable
	categories = c.Categories
	v := reflect.ValueOf(c.Continents)
	typeOfS := v.Type()

	// iterate over each continent
	for i := 0; i < v.NumField(); i++ {
		countries = v.Field(i).Interface().([]string)
		continent = typeOfS.Field(i).Name
		continentsDict[continent] = make(map[string]int)
		// get data for each continent
		dict, err := Fetch(c.Url, countries, categories)
		if err == nil {
			continentsDict[continent] = dict
		} else {
			fmt.Printf("Failed to fetch %s\n", continent)
		}
	}
	return continentsDict
}

func getUrl(url string, countries []string, category string, apiKey string) string {
	url = fmt.Sprintf("%s?apikey=%s", url, apiKey)
	url = fmt.Sprintf("%s&category=%s", url, category)

	if len(countries) != 0 {
		url = fmt.Sprintf("%s&country=%s", url, strings.Join(countries, ","))
	}
	return url
}
