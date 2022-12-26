package extract

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type responseHead struct {
	Status       string `json:"status"`
	TotalResults int64  `json:"totalResults"`
	Articles     []struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Source struct {
			ID   string `json:"Id"`
			Name string `json:"Name"`
		} `json:"source"`
		PublishedAt string `json:"publishedAt"`
		URL         string `json:"url"`
	} `json:"articles"`
}

type Article struct {
	Title  string
	URL    string
	Source string
}

// fetch headlines
func FetchHead(url string, apiKey string, sources []string, pageSize int, args ...string) ([]Article, error) {
	url = GetUrlHead(url, apiKey, sources, pageSize, args...)
	res, err := http.Get(url)
	// get data from news api
	if err != nil {
		return nil, err
	}
	// read data
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// parse response
	var responseObject responseHead
	json.Unmarshal(resp, &responseObject)
	// get articles
	articles := make([]Article, len(responseObject.Articles))
	for i, article := range responseObject.Articles {
		articles[i] = Article{Title: article.Title, URL: article.URL, Source: article.Source.Name}
	}
	return articles, nil
}

func GetUrlHead(url string, apiKey string, sources []string, pageSize int, args ...string) string {
	url = fmt.Sprintf("%s?apiKey=%s&pageSize=%d", url, apiKey, pageSize)
	if len(args) > 0 {
		url = fmt.Sprintf("%s&q=%s", url, args[0])
	}
	if len(sources) != 0 {
		url = fmt.Sprintf("%s&sources=%s", url, strings.Join(sources, ","))
	}
	return url
}
