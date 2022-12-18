package main

import (
	"context"
	"fmt"
	"github.com/alexvishnevskiy/current-news/backend/api"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"log"
	"net/http"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	//r := setupRouter()
	db := api.RedisDB{Ctx: context.TODO()}
	err := db.Connect("localhost:6379")
	if err != nil {
		log.Fatal("Can't connect to the server")
	}

	db.Add("Africa", 10, "business")
	db.Add("Africa", 5, "sports")
	db.Add("Africa", 7, "business")
	db.Add("Africa", 6, "sports")

	res, err := db.Exists("Africa")
	if err == nil {
		fmt.Println(res)
	}
	res, err = db.Exists("America")
	if err == nil {
		fmt.Println(res)
	}

	top, err := db.GetTop("Africa")
	if err == nil {
		fmt.Println(top)
	}

	//a := api.GetCountryData(resp)
	//fmt.Println(a)

	//fmt.Println(resp)
	// Listen and Server in 0.0.0.0:8080
	//r.Run(":8888")
}
