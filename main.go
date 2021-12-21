package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
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
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte
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

	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("screenshot.png", buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
