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
	db            = api.RedisDB{Ctx: context.TODO()}
	timeseries    = api.TimeseriesDB{Ctx: context.TODO()}
	configTest    config.TestConfig
	configApi     config.ConfigCat
	configHead    config.ConfigHead
	configArchive config.ConfigArchive
	configBack    config.ConfigBackend
)

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

	err = configHead.GetConf()
	url = configHead.Url
	apiKey := configHead.ApiKey
	sources := configHead.Sources
	assert.Nil(t, err, nil, "Headlines config is not loaded")
	assert.Greater(t, len(url), 0, "Config for headlines is not loaded correctly")
	assert.Greater(t, len(sources), 0, "Config for headlines is not loaded correctly")
	assert.Greater(t, len(apiKey), 0, "Config for headlines is not loaded correctly")

	err = configArchive.GetConf()
	url = configHead.Url
	apiKey = configHead.ApiKey
	assert.Nil(t, err, nil, "Archive config is not loaded")
	assert.Greater(t, len(url), 0, "Config for archive is not loaded correctly")
	assert.Greater(t, len(apiKey), 0, "Config for archive is not loaded correctly")
}

func TestDB(t *testing.T) {
	// connect databases
	err := db.Connect("localhost:6379")
	assert.Nil(t, err, "Failed to connect to database")
	err = timeseries.Connect("localhost:6379")
	assert.Nil(t, err, "Failed to connect to timeseries database")

	// timeseries tests
	_ = timeseries.Add("timeseries1", 0, 10)
	_ = timeseries.Add("timeseries1", 1, 5)
	_ = timeseries.Add("timeseries1", 2, 8)

	point, err := timeseries.GetLast("timeseries1")
	assert.Nil(t, err, "Timeseries is not working correctly")
	assert.Equal(t, point.Timestamp, int64(2), "Timeseries is not working correctly")
	assert.Equal(t, point.Value, float64(8), "Timeseries is not working correctly")

	points, err := timeseries.GetSeries("timeseries1", 0, 2)
	assert.Nil(t, err, "Timeseries data getter is not working correctly")
	assert.Equal(t, len(points), 3, "Timeseries data getter is not working correctly")
	assert.Equal(t, points[1].Value, float64(5), "Timeseries data getter is not working correctly")

	// clear timeseries
	err = timeseries.Delete("timeseries1")
	assert.Nil(t, err, "It is not possible to delete key from timeseries")

	// database tests
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

	// testing archive
	_ = configArchive.GetConf()
	url = configArchive.Url
	apiKey = configArchive.ApiKey
	numberArticles, err := extract.FetchArchive(url, apiKey, 2019, 3)
	assert.Nil(t, err, "Failed to fetch archive data")
	assert.Greater(t, numberArticles, 0, "Failed to fetch data for headlines")
}

func TestDBApi(t *testing.T) {
	_ = configTest.GetConf()
	_ = configArchive.GetConf()
	_ = configHead.GetConf()
	_ = configBack.GetConf()

	// connect to databases
	db.Connect(configBack.RedisAddr)
	timeseries.Connect(configBack.RedisAddr)

	// test table update
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
	err = api.UpdateHeadlines(&db, configHead)
	assert.Nil(t, err, "Failed update headlines")
	articles, _ := api.GetHeadlines(&db, configHead.SetKey)
	assert.Greater(t, len(articles), 0, "Failed to get headlines")
	// clean
	db.RemoveSet(configHead.SetKey)

	// test archive data
	err = api.InitArchive(&timeseries, configArchive)
	assert.Nil(t, err, "Failed to initialize archive")
	err = api.UpdateArchive(&timeseries, configArchive)
	assert.Nil(t, err, "Failed to update archive")
	archive, err := api.GetArchive(&timeseries, configArchive)
	assert.Greater(t, len(archive), 0, "Failed to get archive")
	err = timeseries.Delete(configArchive.SetKey)
	assert.Nil(t, err, "Failed to delete archive")
}
