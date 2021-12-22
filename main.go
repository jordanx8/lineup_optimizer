package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/jordanx8/lineup_optimizer/player"
	"golang.org/x/term"
)

func getPassword() string {
	fmt.Println("\nPassword: ")
	// https://godoc.org/golang.org/x/crypto/ssh/terminal#ReadPassword
	// terminal.ReadPassword accepts file descriptor as argument, returns byte slice and error.
	passwd, e := term.ReadPassword(int(os.Stdin.Fd()))
	if e != nil {
		log.Fatal(e)
	}
	// Type cast byte slice to string.
	return string(passwd)
}

func main() {
	var username string

	fmt.Println("Enter email address/username: ")
	fmt.Scanln(&username)

	password := getPassword()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	fmt.Println("Received username/password; logging in")
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://login.yahoo.com/`),
		chromedp.WaitVisible(`username-challenge`),
		chromedp.SetValue(`phone-no`, username),
		chromedp.Click(`login-signin`, chromedp.NodeVisible),
		chromedp.WaitVisible(`password-container`),
		chromedp.SetValue(`login-passwd`, password),
		chromedp.Click(`login-signin`, chromedp.NodeVisible),
		chromedp.WaitVisible(`ybar-logo`),
		chromedp.Navigate(`https://basketball.fantasysports.yahoo.com/`),
		chromedp.Sleep(2*time.Second),
		chromedp.Nodes(`I Navtarget yfa-rapid-beacon`, &nodes),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in. Attempting to scrape players information.")

	var url = "https://basketball.fantasysports.yahoo.com/" + nodes[19].AttributeValue("href")

	var playerNames []string
	var playerData []string
	nodes = nil
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`[...document.querySelectorAll('#statTable0 a.Nowrap')].map((e) => e.innerText)`, &playerNames),
		chromedp.Evaluate(`[...document.querySelectorAll('#statTable0 span.Fz-xxs')].map((e) => e.innerText)`, &playerData),
		chromedp.Nodes(`ul.Nav-h.Nav-bot-pointers-north.No-bdr > li.Navitem.Mstart-xxl.Ta-c > a.Navtarget.yfa-rapid-beacon`, &nodes),
	)
	if err != nil {
		log.Fatal(err)
	}

	var playerPointsStrings []string
	var playerPoints []float32
	url = nodes[2].AttributeValue("href")
	fmt.Println("Scanning Day 1")
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`[...document.querySelectorAll('td > div > span.Fw-b')].map((e) => e.innerText)`, &playerPointsStrings),
		chromedp.Click(`Js-next Grid-u No-bdr-radius-start No-bdrstart Pstart-med Td-n Fz-xs`),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range playerPointsStrings {
		if s, err := strconv.ParseFloat(b, 32); err == nil {
			playerPoints = append(playerPoints, float32(s))
		}
	}

	day := 2
	for day < 8 {
		fmt.Printf("Scanning Day %d\n", day)
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`[...document.querySelectorAll('td > div > span.Fw-b')].map((e) => e.innerText)`, &playerPointsStrings),
			chromedp.Click(`Js-next Grid-u No-bdr-radius-start No-bdrstart Pstart-med Td-n Fz-xs `),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			log.Fatal(err)
		}
		for a, b := range playerPointsStrings {
			if s, err := strconv.ParseFloat(b, 32); err == nil {
				playerPoints[a] = playerPoints[a] + float32(s)
			}
		}
		day++
	}

	fmt.Println("Player scrape complete. Creating players array.")

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
		positions = strings.Split(playerData[a][6:], ",")
		positions = player.AddExtraPositions(positions)
		players = append(players, *player.NewPlayer(b, positions, playerData[a+1], playerPoints[c]))
		a = a + 2
		c++
	}

	lineup, players := player.OptimizeLineup(players)

	fmt.Println("Optimized Lineup:")
	fmt.Println("PG   -", lineup["PG"])
	fmt.Println("SG   -", lineup["SG"])
	fmt.Println("G    -", lineup["G"])
	fmt.Println("SF   -", lineup["SF"])
	fmt.Println("PF   -", lineup["PF"])
	fmt.Println("F    -", lineup["F"])
	fmt.Println("C    -", lineup["C"])
	fmt.Println("C    -", lineup["C2"])
	fmt.Println("Util -", lineup["Util"])
	fmt.Println("Util -", lineup["Util2"])
	fmt.Println("BN/IL Players: ", players, "\n")

	fmt.Println("Press Enter to Exit")
	fmt.Scanln()
	fmt.Println("exiting...")
}
