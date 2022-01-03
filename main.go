package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jordanx8/lineup_optimizer/scrape"
	"golang.org/x/term"
)

func getPassword() string {
	fmt.Println("Enter password: ")
	// https://godoc.org/golang.org/x/crypto/ssh/terminal#ReadPassword
	// terminal.ReadPassword accepts file descriptor as argument, returns byte slice and error.
	password, e := term.ReadPassword(int(os.Stdin.Fd()))
	if e != nil {
		log.Fatal(e)
	}
	// Type cast byte slice to string.
	return string(password)
}

func main() {
	var username string

	fmt.Println("Enter email address/username: ")
	fmt.Scanln(&username)

	password := getPassword()

	lineup, bench := scrape.YahooScrape(username, password)
	var sum float32
	for _, v := range lineup {
		sum += v.Points
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/styles.css", "./static/css/styles.css")
	router.StaticFile("/loginstyles.css", "./static/css/loginstyles.css")
	router.GET("/", func(c *gin.Context) {

		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"playertable.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"playerlineup": lineup,
				"bench":        bench,
				"total":        sum,
			},
		)
	})
	router.Run()
}
