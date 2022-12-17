package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
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

type config struct {
	Url        string `mapstructure:"url"`
	Continents struct {
		NorthAmerica []string `mapstructure:"NorthAmerica"`
		Europe       []string `mapstructure:"Europe"`
		Asia         []string `mapstructure:"Asia"`
		Africa       []string `mapstructure:"Africa"`
		SouthAmerica []string `mapstructure:"SouthAmerica"`
	} `mapstructure:"continent"`
	Categories []string `mapstructure:"category"`
}

func Fetch() (map[string]map[string]int, error) {
	apiKey := os.Getenv("API_KEY")
	continentsDict := make(map[string]map[string]int) //continent->category->#articles

	var (
		c          config
		url        string
		continent  string
		countries  []string
		categories []string
	)
	c.getConf()
	categories = c.Categories

	v := reflect.ValueOf(c.Continents)
	typeOfS := v.Type()
	// iterate over each continent and category
	for i := 0; i < v.NumField(); i++ {
		countries = v.Field(i).Interface().([]string)
		continent = typeOfS.Field(i).Name
		continentsDict[continent] = make(map[string]int)

		for _, category := range categories {
			url = getUrl(c.Url, countries, category, apiKey)
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
			continentsDict[continent][category] = responseObject.TotalResults
		}
	}
	return continentsDict, nil
}

func getUrl(url string, countries []string, category string, apiKey string) string {
	url = fmt.Sprintf("%s?apikey=%s", url, apiKey)
	url = fmt.Sprintf("%s&category=%s", url, category)

	if len(countries) != 0 {
		url = fmt.Sprintf("%s&country=%s", url, strings.Join(countries, ","))
	}
	return url
}

func (c *config) getConf() *config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	path, _ := os.Getwd()
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath(fmt.Sprintf("%s/configs", path))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Failed to load config")
	}

	err := viper.Unmarshal(c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}
