package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordanx8/lineup_optimizer/player"
	"github.com/jordanx8/lineup_optimizer/scrape"
)

// func getPassword() string {
// 	fmt.Println("Enter password: ")
// 	// https://godoc.org/golang.org/x/crypto/ssh/terminal#ReadPassword
// 	// terminal.ReadPassword accepts file descriptor as argument, returns byte slice and error.
// 	password, e := term.ReadPassword(int(os.Stdin.Fd()))
// 	if e != nil {
// 		log.Fatal(e)
// 	}
// 	// Type cast byte slice to string.
// 	return string(password)
// }

var lineup []player.Player
var bench []player.Player
var sum float32

func main() {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> 97cc3f0 (login page working/implemented new optimization)
	// var username string

	// fmt.Println("Enter email address/username: ")
	// fmt.Scanln(&username)

	// // password := getPassword()

	// lineup, bench := scrape.YahooScrape("***REMOVED***", "***REMOVED***")
	// var sum float32
<<<<<<< HEAD
>>>>>>> c38d8ca (exe)
=======
>>>>>>> 30f91a2 (loading page created)
=======
>>>>>>> 97cc3f0 (login page working/implemented new optimization)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/styles.css", "./static/css/styles.css")
	router.StaticFile("/loginstyles.css", "./static/css/loginstyles.css")
	router.StaticFile("/loading.css", "./static/css/loading.css")
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
<<<<<<< HEAD
			"login.html",
=======
			// Use the index.html template
			"login.html",
			// Pass the data that the page uses (in this case, 'title')
>>>>>>> 97cc3f0 (login page working/implemented new optimization)
			gin.H{},
		)
	})
	router.POST("/", performLogin)
<<<<<<< HEAD
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
=======
	router.GET("/table", func(c *gin.Context) {

		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
>>>>>>> 97cc3f0 (login page working/implemented new optimization)
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

<<<<<<< HEAD
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
=======
func loginScrape(username, password string) bool {
	return false
}

func performLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	lineup, bench = scrape.YahooScrape(username, password)
	for _, v := range lineup {
		sum += v.Points
	}
	c.Redirect(http.StatusMovedPermanently, "/table")
>>>>>>> 97cc3f0 (login page working/implemented new optimization)
}
