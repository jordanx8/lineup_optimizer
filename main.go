package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	robotgo.Move(345, 51)
	robotgo.Click()
	robotgo.TypeStr("https://basketball.fantasysports.yahoo.com/")
	robotgo.KeyTap("enter")
	robotgo.Sleep(5)
	robotgo.Move(765, 522)       //my team
	robotgo.MoveSmooth(775, 586) //edit weekly lineup
	robotgo.Click()              //left click
	x, y := robotgo.GetMousePos()
	fmt.Println("pos: ", x, y)
}
