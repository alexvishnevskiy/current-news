package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/alexvishnevskiy/current-news/api"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

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
	db.Connect("localhost:6379")

	// setup config
	var conf api.ConfigAPI
	conf.GetConf()
	continents := reflect.ValueOf(conf.GetContinents()).MapKeys()
	continentsNames := make([]string, len(continents))
	for i := 0; i < len(continents); i++ {
		continentsNames[i] = continents[i].String()
	}

	// run background job to update table
	s := gocron.NewScheduler(time.UTC)
	job, _ := s.Every(8).Hour().Do(api.UpdateTable, &db, &conf)
	s.StartAsync()

	//if it is a first time, we should update table first
	if exists, _ := db.SetExists(continents[0].String()); !exists {
		for {
			if !job.IsRunning() {
				break
			}
		}
	}

	r.GET("/data", func(c *gin.Context) {
		data := make(map[string][]string)
		// get data from db
		for _, continent := range continentsNames {
			res, err := api.GetData(&db, continent)
			if err != nil {
				log.Fatal("There is no data")
			}
			data[continent] = res
		}

		jsonData, _ := json.Marshal(data)
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8888")
}
