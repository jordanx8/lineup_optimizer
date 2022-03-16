package scrape

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/jordanx8/lineup_optimizer/player"
)

var playerPointsStrings []string
var playerPoints []float32
var playerNames []string
var playerData []string

func YahooScrape(username string, password string) ([]player.Player, []player.Player, string) {
	playerPointsStrings = nil
	playerPoints = nil
	playerNames = nil
	playerData = nil

	// creates context with ExecAllocator
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		// block all images
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
	)
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	// new browser, first tab
	browserContext, cancel := chromedp.NewContext(
		allocatorCtx,
	)
	defer cancel()

	// ensure first tab starts
	if err := chromedp.Run(browserContext); err != nil {
		log.Fatal(err)
		return nil, nil, err.Error()
	}

	// same browser, another tab for login with a 25 second timeout to ensure login went through
	loginTab, cancel := context.WithTimeout(browserContext, 25*time.Second)
	defer cancel()

	fmt.Println("Received username/password; logging in")
	// logs into yahoo with username/password and gets link for weekly lineup
	editWeeklyLineupURL, err := attemptLogin(loginTab, username, password)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, nil, err.Error()
		}
		log.Fatal(err)
		return nil, nil, err.Error()
	}
	fmt.Println("Login successful. Attempting to scrape players' information.")

	urls := getDateURLs(editWeeklyLineupURL)

	playerNames, playerData, err := gatherPlayerInfo(browserContext, editWeeklyLineupURL)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, nil, err.Error()
		}
		log.Fatal(err)
		return nil, nil, err.Error()
	}

	scanDay(browserContext, urls[0], 1)
	for _, b := range playerPointsStrings {
		if s, err := strconv.ParseFloat(b, 32); err == nil {
			playerPoints = append(playerPoints, float32(s))
		}
	}

	// loops and goes through each day of projected fantasy scores for the week and adds them together
	day := 2
	for day < 8 {
		scanDay(browserContext, urls[day-1], day)
		day++
	}

	//wait for variables to be set, then create necessary data structures
	var players []player.Player
	var positions []string
	a := 0
	c := 0
	for _, b := range playerNames {
		//this checks for empty spots in the set lineup
		if playerData[a] == "" {
			a++
			c++
		}
		if playerData[a+1] == "INJ" {
			playerPoints[c] = 0
		}
		var info string = playerData[a]
		positions = strings.Split(playerData[a][6:], ",")
		positions = player.AddExtraPositions(positions)
		players = append(players, *player.NewPlayer(b, positions, playerData[a+1], info, playerPoints[c]))
		a = a + 2
		c++
	}

	lineup, bench := player.OptimizeLineup(players)
	return lineup, bench, "nil"
}

func scanDay(browser context.Context, url string, day int) {
	fmt.Printf("Scanning Day %d\n", day)
	defer fmt.Printf("Day %d Scanned\n", day)
	newTab, cancel := context.WithTimeout(browser, 25*time.Second)
	defer cancel()
	err := chromedp.Run(newTab,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`td > div > span.Fw-b`),
		chromedp.Evaluate(`[...document.querySelectorAll('td > div > span.Fw-b')].map((e) => e.innerText)`, &playerPointsStrings),
	)
	if err != nil {
		log.Fatal(err)
	}
	if day > 1 {
		for a, b := range playerPointsStrings {
			if s, err := strconv.ParseFloat(b, 32); err == nil {
				playerPoints[a] = playerPoints[a] + float32(s)
			}
		}
	}
}

func getDateURLs(originalURL string) []string {
	var dateURLs []string

	arrayURL := []rune(originalURL)
	dateSubstring := string(arrayURL[len(originalURL)-10 : len(originalURL)])
	beginningSubstring := string(arrayURL[:len(originalURL)-16])

	day, err := time.Parse("2006-01-02", dateSubstring)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 7; i++ {
		dateURLs = append(dateURLs, beginningSubstring+"/team?&date="+day.Format("2006-01-02")+"&stat1=P&stat2=P")
		day = day.AddDate(0, 0, 1)
	}

	return dateURLs
}

func gatherPlayerInfo(browser context.Context, url string) ([]string, []string, error) {
	newTab, cancel := context.WithTimeout(browser, 25*time.Second)
	defer cancel()

	// navigates to weekly lineup and gathers players' names and info
	err := chromedp.Run(newTab,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#statTable0 a.Nowrap`),
		chromedp.WaitVisible(`#statTable0 span.Fz-xxs`),
		chromedp.Evaluate(`[...document.querySelectorAll('#statTable0 a.Nowrap')].map((e) => e.innerText)`, &playerNames),
		chromedp.Evaluate(`[...document.querySelectorAll('#statTable0 span.Fz-xxs')].map((e) => e.innerText)`, &playerData),
	)
	return playerNames, playerData, err
}

func attemptLogin(tab context.Context, username string, password string) (string, error) {
	var nodes []*cdp.Node
	var editWeeklyLineupURL string
	err := chromedp.Run(tab,
		chromedp.Navigate(`https://login.yahoo.com/`),
		chromedp.WaitVisible(`username-challenge`),
		chromedp.SetValue(`phone-no`, username),
		chromedp.Click(`login-signin`, chromedp.NodeVisible),
		chromedp.WaitVisible(`password-container`),
		chromedp.SetValue(`login-passwd`, password),
		chromedp.Click(`login-signin`, chromedp.NodeVisible),
		chromedp.WaitVisible(`ybar-logo`),
		chromedp.Navigate(`https://basketball.fantasysports.yahoo.com/`),
		chromedp.Sleep(1*time.Second),
		chromedp.Nodes(`a[href*="date"]`, &nodes),
	)
	if err == nil {
		editWeeklyLineupURL = "https://basketball.fantasysports.yahoo.com" + nodes[0].AttributeValue("href")
	}
	return editWeeklyLineupURL, err
}
