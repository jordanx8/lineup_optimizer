package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/jordanx8/lineup_optimizer/player"
)

func main() {
	var username string
	var password string

	fmt.Println("Enter email address/username: ")
	fmt.Scanln(&username)

	fmt.Println("Enter password: ")
	fmt.Scanln(&password)

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

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
	var url = "https://basketball.fantasysports.yahoo.com/" + nodes[19].AttributeValue("href")

	var playerNames []string
	var playerData []string
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`[...document.querySelectorAll('#statTable0 a.Nowrap')].map((e) => e.innerText)`, &playerNames),
		chromedp.Evaluate(`[...document.querySelectorAll('#statTable0 span.Fz-xxs')].map((e) => e.innerText)`, &playerData),
	)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(playerNames)
	// fmt.Println(playerData)

	var players []player.Player
	var positions []string
	a := 0

	// TODO: add player projected weekly points total
	for _, b := range playerNames {
		positions = strings.Split(playerData[a][6:], ",")
		positions = player.AddExtraPositions(positions)
		players = append(players, *player.NewPlayer(b, positions, playerData[a+1], 0))
		a = a + 2
	}
	fmt.Println("Players: ", players)
}
