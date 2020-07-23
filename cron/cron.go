package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()

	callbackFunc := func() {
		fmt.Println("tick every 1 second")
	}
	c.AddFunc("@every 1s", callbackFunc)

	c.Start()
	time.Sleep(time.Second * 500)
}
