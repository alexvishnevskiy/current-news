package main_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/alexvishnevskiy/current-news/api"
	"github.com/alexvishnevskiy/current-news/api/config"
	"github.com/alexvishnevskiy/current-news/api/extract"
	"github.com/stretchr/testify/assert"
)

var (
	db         = api.RedisDB{Ctx: context.TODO()}
	configTest config.TestConfig
	configApi  config.ConfigCat
	configHead config.ConfigHead
)

func TestDB(t *testing.T) {
	db.Connect("localhost:6379")
	_ = db.Add("Africa", 10, "business")
	_ = db.Add("Africa", 8, "politics")
	_ = db.Add("America", 9, "sports")

	// tests for member existence
	res, _ := db.MemberExists("Africa", "business")
	assert.Equal(t, res, true, "Africa:business should exist in database")
	res, _ = db.MemberExists("Africa", "sports")
	assert.Equal(t, res, false, "Africa:sports shouldn't exist in database")
	db.RemoveMember("Africa", "politics")
	res, _ = db.MemberExists("Africa", "politics")
	assert.Equal(t, res, false, "Africa:politics shouldn't exist in database")

	// tests for set existence
	res, _ = db.SetExists("Africa")
	assert.Equal(t, res, true, "Africa should exist in database")
	res, _ = db.SetExists("America")
	assert.Equal(t, res, true, "America value should exist in database")
	res, _ = db.SetExists("Europe")
	assert.Equal(t, res, false, "Europe value shouldn't exist in database")
	_ = db.RemoveSet("America")
	res, _ = db.SetExists("America")
	assert.Equal(t, res, false, "America value shouldn't exist in database")

	// clear everything
	_ = db.RemoveSet("Africa")
	_ = db.RemoveSet("America")

	// test for normal set
	db.AddSet("articles", extract.Article{Title: "article1", URL: "url1"})
	db.AddSet("articles", extract.Article{Title: "article2", URL: "url2"})
	db.AddSet("articles", extract.Article{Title: "article3", URL: "url3"})
	size, _ := db.Size("articles")
	assert.Equal(t, size, int64(3), "The size of the articles set should be 3")
	_ = db.RemoveSet("articles")

	// test for the topk
	_ = db.Add("Africa", 0, "business")
	_ = db.Add("Africa", 5, "politics")
	_ = db.Add("Africa", 9, "sports")
	_ = db.Add("Africa", 4, "education")
	topk, _ := db.GetTop("Africa")
	assert.Equal(
		t, topk, []api.Member{
			{Member: "sports", Score: 9},
			{Member: "politics", Score: 5},
			{Member: "education", Score: 4},
			{Member: "business", Score: 0},
		}, "Topk is not right")

	// clear everything
	_ = db.RemoveSet("Africa")
}

func TestConfig(t *testing.T) {
	var (
		err        error
		url        string
		continents []string
	)
	err = configTest.GetConf()
	url = configTest.GetUrl()
	continents = configTest.GetContinents()["Europe"]

	assert.Nil(t, err, nil, "Test config is not loaded")
	assert.Greater(t, len(continents), 0, "Test config is not loaded correctly")
	assert.Greater(t, len(url), 0, "Test config is not loaded correctly")

	err = configApi.GetConf()
	url = configApi.GetUrl()
	continents = configApi.GetContinents()["Europe"]
	assert.Nil(t, err, nil, "Api config is not loaded")
	assert.Greater(t, len(continents), 0, "Config is not loaded correctly")
	assert.Greater(t, len(url), 0, "Config is not loaded correctly")

	_ = configHead.GetConf()
	url = configHead.Url
	apiKey := configHead.ApiKey
	sources := configHead.Sources
	assert.Greater(t, len(url), 0, "Config for headlines is not loaded correctly")
	assert.Greater(t, len(sources), 0, "Config for headlines is not loaded correctly")
	assert.Greater(t, len(apiKey), 0, "Config for headlines is not loaded correctly")
}

func TestFetch(t *testing.T) {
	// testing categories
	_ = configTest.GetConf()
	url := configTest.GetUrl()
	countries := configTest.GetContinents()["Europe"]
	categories := configTest.GetCategories()
	apiKey := configTest.GetApiKey()

	res, err := extract.FetchCat(url, countries, categories, apiKey)
	assert.Greater(t, res[categories[0]], 0, "Failed to fetch data for categories")
	assert.Nil(t, err, "Failed to fetch data for categories")

	// testings headlines
	_ = configHead.GetConf()
	url = configHead.Url
	apiKey = configHead.ApiKey
	pageSize := configHead.PageSize
	sources := configHead.Sources
	articles, err := extract.FetchHead(url, apiKey, sources, pageSize)
	assert.Greater(t, len(articles), 0, "Failed to fetch data for headlines")
	assert.Nil(t, err, "Failed to fetch data for headlines")
}

func TestDBApi(t *testing.T) {
	db.Connect("localhost:6379")

	// test table update
	_ = configTest.GetConf()
	err := api.UpdateSortedSet(&db, &configTest)
	assert.Nil(t, err, "Failed to update table with api")

	continents := reflect.ValueOf(configTest.GetContinents()).MapKeys()
	exists, err := db.SetExists(continents[0].String())
	assert.Nil(t, err, "Api didn't update table")
	assert.Equal(t, exists, true, "Api didn't update table")

	exists, err = db.SetExists(continents[1].String())
	assert.Nil(t, err, "Api didn't update table")
	assert.Equal(t, exists, true, "Api didn't update table")

	// test fetching data from table
	res, err := api.GetCategoriesTop(&db, continents[0].String())
	assert.Nil(t, err, "Failed to fetch data")
	assert.Greater(t, len(res), 0, "Failed to fetch data")

	// clean data
	for _, continent := range continents {
		db.RemoveSet(continent.String())
	}

	// test fetching headlines
	_ = configHead.GetConf()
	err = api.UpdateHeadlines(&db, configHead)
	assert.Nil(t, err, "Failed update headlines")
	articles, _ := api.GetHeadlines(&db, configHead.SetKey)
	assert.Greater(t, len(articles), 0, "Failed to get headlines")
	// clean
	db.RemoveSet(configHead.SetKey)
}
