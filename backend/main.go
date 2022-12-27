package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/alexvishnevskiy/current-news/api"
	"github.com/alexvishnevskiy/current-news/api/config"
	"github.com/alexvishnevskiy/current-news/api/extract"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func main() {
	// setup router
	r := gin.Default()
	r.Use(CORSMiddleware())

	// setup config for categories
	var conf config.ConfigCat
	conf.GetConf()
	continents := reflect.ValueOf(conf.GetContinents()).MapKeys()
	continentsNames := make([]string, len(continents))
	for i := 0; i < len(continents); i++ {
		continentsNames[i] = continents[i].String()
	}
	// setup config for headlines
	var confHead config.ConfigHead
	confHead.GetConf()
	// setup backend config
	var confBack config.ConfigBackend
	confBack.GetConf()
	// setup config for archive
	var confArchive config.ConfigArchive
	confArchive.GetConf()

	// setup databases
	var db = api.RedisDB{Ctx: context.TODO()}
	var dbTime = api.TimeseriesDB{Ctx: context.TODO()}
	err := db.Connect(confBack.RedisAddr)
	if err != nil {
		panic(err)
	}
	err = dbTime.Connect(confBack.RedisAddr)
	if err != nil {
		panic(err)
	}

	// run job to initialize archive
	if exists, err := db.SetExists(confArchive.SetKey); !exists && err == nil {
		err = api.InitArchive(&dbTime, confArchive)
		if err != nil {
			panic(err)
		}
	}

	// run background job to update categories
	s := gocron.NewScheduler(time.UTC)
	jobCat, _ := s.Every(conf.GetFrequency()).Hour().Do(api.UpdateSortedSet, &db, &conf)
	// run background job to update headlines
	jobHead, _ := s.Every(confHead.Frequency).Hour().Do(api.UpdateHeadlines, &db, confHead)
	// run background job to update headlines
	updateFreq := 24 * 30
	if confArchive.Frequency == "year" {
		updateFreq *= 12
	}
	jobArchive, _ := s.Every(updateFreq).Hour().Do(api.UpdateArchive, &dbTime, confArchive)
	s.StartAsync()

	// if it is a first time, we should update tables first
	if exists, _ := db.SetExists(continents[0].String()); !exists {
		for {
			if !jobCat.IsRunning() && !jobHead.IsRunning() && !jobArchive.IsRunning() {
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
		articles, err := extract.FetchHead(confHead.Url, confHead.ApiKey, []string{}, confHead.PageSize, question)
		// when there is no apis calls left display smth
		if err != nil {
			log.Fatal("Failed to fetched articles for specific question\n", err)
		}

		jsonData, _ := json.Marshal(articles)
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	r.GET("/archive", func(c *gin.Context) {
		timeseries, err := api.GetArchive(&dbTime, confArchive)
		if err != nil {
			log.Fatal("Failed to fetched archive articles\n", err)
		}

		sendDict := make(map[string][]api.TimeStamp)
		sendDict["archive"] = timeseries
		jsonData, _ := json.Marshal(sendDict)
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	// Listen and Serve
	r.Run(confBack.BackendPort)
}

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
