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

func YahooScrape(username string, password string) ([]player.Player, []player.Player) {
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		// block all images
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
	)
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	// create chrome instance
	bctx, cancel := chromedp.NewContext(
		allocatorCtx,
	)
	defer cancel()

	// start the browser first
	if err := chromedp.Run(bctx); err != nil {
		log.Fatal(err)
	}

	// create a timeout
	timeoutctx, cancel := context.WithTimeout(bctx, 25*time.Second)
	defer cancel()

	fmt.Println("Received username/password; logging in")
	var nodes []*cdp.Node
	// logs into yahoo with username/password and gets link for weekly lineup
	err := chromedp.Run(timeoutctx,
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
		chromedp.Nodes(`I Navtarget yfa-rapid-beacon`, &nodes),
	)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, nil
		}
		log.Fatal(err)
	}
	fmt.Println("Logged in. Attempting to scrape players' information.")

	var url = "https://basketball.fantasysports.yahoo.com/" + nodes[19].AttributeValue("href")

	var playerNames []string
	var playerData []string
	nodes = nil
	// navigates to weekly lineup and gathers players' names and info
	err = chromedp.Run(bctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(`[...document.querySelectorAll('#statTable0 a.Nowrap')].map((e) => e.innerText)`, &playerNames),
		chromedp.Evaluate(`[...document.querySelectorAll('#statTable0 span.Fz-xxs')].map((e) => e.innerText)`, &playerData),
		chromedp.Nodes(`ul.Nav-h.Nav-bot-pointers-north.No-bdr > li.Navitem.Mstart-xxl.Ta-c > a.Navtarget.yfa-rapid-beacon`, &nodes),
	)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, nil
		}
		log.Fatal(err)
	}
	var playerPointsStrings []string
	var playerPoints []float32
	url = nodes[2].AttributeValue("href")
	fmt.Println("Scanning Day 1")
	// navigates to tab of daily projected fantasy scores and gathers first day of players' projected scores
	err = chromedp.Run(bctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(`[...document.querySelectorAll('td > div > span.Fw-b')].map((e) => e.innerText)`, &playerPointsStrings),
		chromedp.Click(`Js-next Grid-u No-bdr-radius-start No-bdrstart Pstart-med Td-n Fz-xs`),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, nil
		}
		log.Fatal(err)
	}
	for _, b := range playerPointsStrings {
		if s, err := strconv.ParseFloat(b, 32); err == nil {
			playerPoints = append(playerPoints, float32(s))
		}
	}

	// loops and goes through each day of projected fantasy scores for the week and adds them together
	day := 2
	for day < 8 {
		fmt.Printf("Scanning Day %d\n", day)
		err = chromedp.Run(bctx,
			chromedp.Evaluate(`[...document.querySelectorAll('td > div > span.Fw-b')].map((e) => e.innerText)`, &playerPointsStrings),
			chromedp.Click(`Js-next Grid-u No-bdr-radius-start No-bdrstart Pstart-med Td-n Fz-xs `),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			if err == context.DeadlineExceeded {
				return nil, nil
			}
			log.Fatal(err)
		}
		for a, b := range playerPointsStrings {
			if s, err := strconv.ParseFloat(b, 32); err == nil {
				playerPoints[a] = playerPoints[a] + float32(s)
			}
		}
		day++
	}

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

	return lineup, bench
}
