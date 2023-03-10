package extract

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alexvishnevskiy/current-news/api/config"
)

type responseCat struct {
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

// fetch categories
func FetchCat(url string, countries []string, categories []string, apiKey string) (map[string]int, error) {
	categoryDict := make(map[string]int) //category->#articles

	// iterate over each continent and category
	for _, category := range categories {
		newUrl := GetUrlCat(url, countries, category, apiKey)
		// get data from news api
		res, err := http.Get(newUrl)
		if err != nil {
			return nil, err
		}
		// read data
		resp, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		// parse response
		var responseObject responseCat
		json.Unmarshal(resp, &responseObject)
		categoryDict[category] = responseObject.TotalResults
	}
	return categoryDict, nil
}

func FetchAllData(c config.Config) map[string]map[string]int {
	// dict to store all data
	continentsDict := make(map[string]map[string]int)
	categories := c.GetCategories()
	continents := c.GetContinents()
	url := c.GetUrl()
	apiKey := c.GetApiKey()

	// iterate over each continent
	for continent, countries := range continents {
		continentsDict[continent] = make(map[string]int)
		// get data for each continent
		dict, err := FetchCat(url, countries, categories, apiKey)
		if err == nil {
			continentsDict[continent] = dict
		} else {
			fmt.Printf("Failed to fetch %s\n", continent)
		}
	}
	return continentsDict
}

func GetUrlCat(url string, countries []string, category string, apiKey string) string {
	url = fmt.Sprintf("%s?apikey=%s&category=%s", url, apiKey, category)

	if len(countries) != 0 {
		url = fmt.Sprintf("%s&country=%s", url, strings.Join(countries, ","))
	}
	return url
}
