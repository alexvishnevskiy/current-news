package extract

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// fetch archive
// returns the number of articles for year/month
func FetchArchive(url string, apiKey string, year int, month int) (int, error) {
	url = GetUrlArchive(url, apiKey, month, year)
	// get data from news api
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	// read data
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	// parse response
	responseObject := make(map[string]interface{})
	err = json.Unmarshal(resp, &responseObject)
	if err != nil {
		return 0, err
	}
	// preprocess, extract result
	response := responseObject["response"].(map[string]interface{})
	docs := response["docs"].([]interface{})
	return len(docs), nil
}

func GetUrlArchive(url string, apiKey string, month int, year int) string {
	url = fmt.Sprintf("%s/%d/%d.json?api-key=%s", url, year, month, apiKey)
	return url
}
