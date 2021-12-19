package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	cdp "github.com/chromedp/chromedp"
)

func main() {
	var username string
	var password string

	fmt.Println("Enter email address/username: ")
	fmt.Scanln(&username)

	fmt.Println("Enter password: ")
	fmt.Scanln(&password)

	// create chrome instance
	ctx, cancel := cdp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte
	err := cdp.Run(ctx,
		cdp.Navigate(`https://login.yahoo.com/`),
		cdp.WaitVisible(`username-challenge`),
		cdp.SetValue(`phone-no`, username),
		cdp.Click(`login-signin`, cdp.NodeVisible),
		cdp.WaitVisible(`password-container`),
		cdp.SetValue(`login-passwd`, password),
		cdp.Click(`login-signin`, cdp.NodeVisible),
		cdp.WaitVisible(`ybar-logo`),
		cdp.Navigate(`https://basketball.fantasysports.yahoo.com/`),
		cdp.Sleep(2*time.Second),
		cdp.CaptureScreenshot(&buf),
	)

	ioutil.WriteFile("screenshot.png", buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
