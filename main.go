package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jordanx8/lineup_optimizer/player"
	"github.com/jordanx8/lineup_optimizer/scrape"
)

var lineup []player.Player
var bench []player.Player
var sum float32

func main() {

	router := gin.Default()
	var files []string
	filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})

	router.LoadHTMLFiles(files...)
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
	router.GET("/error", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"error.html",
			gin.H{},
		)
	})
	router.POST("/error", returnToLogin)
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

func returnToLogin(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/")
}

func performLogin(c *gin.Context) {
	sum = 0
	username := c.PostForm("username")
	password := c.PostForm("password")
	lineup, bench = scrape.YahooScrape(username, password)
	if lineup == nil || bench == nil {
		c.Redirect(http.StatusMovedPermanently, "/error")
		return
	}
	for _, v := range lineup {
		sum += v.Points
	}
	c.Redirect(http.StatusMovedPermanently, "/table")
}
