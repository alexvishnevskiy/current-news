package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/alexvishnevskiy/current-news/api"
	cat "github.com/alexvishnevskiy/current-news/api/categories"
	head "github.com/alexvishnevskiy/current-news/api/headlines"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

type DataOutput struct {
	Total      int
	Categories []string
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// setup router
	r := gin.Default()
	r.Use(CORSMiddleware())

	// setup database
	var db = api.RedisDB{Ctx: context.TODO()}
	err := db.Connect("localhost:6379")
	// err := db.Connect("redis:6379")
	if err != nil {
		panic(err)
	}

	// setup config for categories
	var conf cat.ConfigAPI
	conf.GetConf()
	continents := reflect.ValueOf(conf.GetContinents()).MapKeys()
	continentsNames := make([]string, len(continents))
	for i := 0; i < len(continents); i++ {
		continentsNames[i] = continents[i].String()
	}
	// setup config for headlines
	var confHead head.ConfigAPI
	confHead.GetConf()

	// run background job to update categories
	s := gocron.NewScheduler(time.UTC)
	jobCat, _ := s.Every(8).Hour().Do(api.UpdateSortedSet, &db, &conf)
	// run background job to update headlines
	jobHead, _ := s.Every(2).Hour().Do(api.UpdateHeadlines, &db, confHead)
	s.StartAsync()

	// if it is a first time, we should update tables first
	if exists, _ := db.SetExists(continents[0].String()); !exists {
		for {
			if !jobCat.IsRunning() && !jobHead.IsRunning() {
				break
			}
		}
	}

	r.GET("/categories", func(c *gin.Context) {
		data := make(map[string]DataOutput)
		// get data from db
		for _, continent := range continentsNames {
			res, err := api.GetCategoriesTop(&db, continent)
			if err != nil {
				log.Fatal("There is no data for categories\n", err)
			}

			total := 0
			categories := make([]string, 0)
			for _, cand := range res {
				categories = append(categories, cand.Member)
				total += cand.Score
			}
			data[continent] = DataOutput{Total: total, Categories: categories}
		}

		jsonData, _ := json.Marshal(data)
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	r.GET("/headlines", func(c *gin.Context) {
		headlines, err := api.GetHeadlines(&db, confHead.SetKey)
		if err != nil {
			log.Fatal("There is no data for headlines\n", err)
		}

		jsonData, _ := json.Marshal(headlines)
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	r.GET("/headlines/:q", func(c *gin.Context) {
		question := c.Param("q")
		articles, err := head.Fetch(confHead.Url, confHead.ApiKey, []string{}, confHead.PageSize, question)
		// when there is no apis calls left display smth
		if err != nil {
			log.Fatal("Failed to fetched articles for specific question\n", err)
		}

		jsonData, _ := json.Marshal(articles)
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
	// r.Run(":80")
}
