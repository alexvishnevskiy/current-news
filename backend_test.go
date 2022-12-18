package backend_test

import (
	"context"
	"github.com/alexvishnevskiy/current-news/backend/api"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	db = api.RedisDB{Ctx: context.TODO()}
	c  api.Config
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
	err := c.GetConf()
	assert.Nil(t, err, nil, "Config is not loaded")
	assert.Greater(t, len(c.Continents.Africa), 0, "Config is not loaded correctly")
	assert.Greater(t, len(c.Url), 0, "Config is not loaded correctly")
}

func TestFetch(t *testing.T) {
	_ = c.GetConf()
	res, err := api.Fetch(c.Url, c.Continents.Europe, c.Categories)
	assert.Greater(t, res[c.Categories[0]], 0, "Failed to fetch data")
	assert.Nil(t, err, "Failed to fetch data")
}

func TestDBApi(t *testing.T) {
	db.Connect("localhost:6379")
	_ = c.GetConf()
	err := api.UpdateTable(&db)
	assert.Nil(t, err, "Failed to update table with api")

	v := reflect.ValueOf(c.Continents)
	typeOfS := v.Type()

	// test table update
	exists, err := db.SetExists(typeOfS.Field(0).Name)
	assert.Nil(t, err, "Api didn't update table")
	assert.Equal(t, exists, true, "Api didn't update table")
	exists, err = db.SetExists(typeOfS.Field(1).Name)
	assert.Nil(t, err, "Api didn't update table")
	assert.Equal(t, exists, true, "Api didn't update table")

	// test fetching data from table
	res, err := api.GetData(&db, typeOfS.Field(0).Name)
	assert.Nil(t, err, "Failed to fetch data")
	assert.Greater(t, len(res), 0, "Failed to fetch data")

	// clean data
	for i := 0; i < v.NumField(); i++ {
		continent := typeOfS.Field(i).Name
		db.RemoveSet(continent)
	}
}
