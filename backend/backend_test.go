package main_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/alexvishnevskiy/current-news/api"
	"github.com/stretchr/testify/assert"
)

var (
	db         = api.RedisDB{Ctx: context.TODO()}
	configTest api.TestConfig
	configApi  api.ConfigAPI
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

	// test for the topk
	_ = db.Add("Africa", 0, "business")
	_ = db.Add("Africa", 5, "politics")
	_ = db.Add("Africa", 9, "sports")
	_ = db.Add("Africa", 4, "education")
	topk, _ := db.GetTop("Africa")
	assert.Equal(
		t, topk, []string{"sports", "politics", "education", "business"}, "Topk is not right")

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
}

func TestFetch(t *testing.T) {
	_ = configTest.GetConf()
	url := configTest.GetUrl()
	countries := configTest.GetContinents()["Europe"]
	categories := configTest.GetCategories()
	apiKey := configTest.GetApiKey()

	res, err := api.Fetch(url, countries, categories, apiKey)
	assert.Greater(t, res[categories[0]], 0, "Failed to fetch data")
	assert.Nil(t, err, "Failed to fetch data")
}

func TestDBApi(t *testing.T) {
	db.Connect("localhost:6379")

	// test table update
	_ = configTest.GetConf()
	err := api.UpdateTable(&db, &configTest)
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
}
