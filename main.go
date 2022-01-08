package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordanx8/lineup_optimizer/player"
	"github.com/jordanx8/lineup_optimizer/scrape"
)

var lineup []player.Player
var bench []player.Player
var sum float32

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/styles.css", "./static/css/styles.css")
	router.StaticFile("/loginstyles.css", "./static/css/loginstyles.css")
	router.StaticFile("/loading.css", "./static/css/loading.css")
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"login.html",
			gin.H{},
		)
	})
	router.POST("/", performLogin)
	router.GET("/loading", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"loading.html",
			gin.H{},
		)
	})
	router.GET("/table", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"playertable.html",
			gin.H{
				"playerlineup": lineup,
				"bench":        bench,
				"total":        sum,
			},
		)
	})
	router.Run()
}

func performLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	ctx, url, cancel := scrape.LogIn(username, password)
	if url == "" {
		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		c.Redirect(http.StatusMovedPermanently, "/loading")
		lineup, bench = scrape.YahooScrape(ctx, url, cancel)
		for _, v := range lineup {
			sum += v.Points
		}
		c.Redirect(http.StatusMovedPermanently, "/table")
	}
}
