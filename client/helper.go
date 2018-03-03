package main

import (
	"math/rand"
	"time"
)

func generateSleepTime(endTime time.Time) time.Duration {
	var d int

	t := time.Now()

	rand.Seed(t.UnixNano())
	x := rand.Intn(100)

	switch {
	case x < 60: // 60%: wakeup at end time - 3 seconds
		d = -3
	case x >= 60 && x < 80: // 20%: wakeup at end time - 5 seconds
		d = -5
	default: // 10%: wakeup at end time - 8 seconds
		d = -8
	}

	wakeupTime := endTime.Add(time.Duration(d) * time.Second)
	return wakeupTime.Sub(t)
}

func generatePhaseTwoPrice(startPrice int64) int64 {
	var price int64

	t := time.Now()

	rand.Seed(t.UnixNano())
	x := rand.Intn(100)

	switch {
	case x < 25:
		price = startPrice
	case x >= 25 && x < 50:
		price = startPrice + 100
	case x <= 50 && x < 75:
		price = startPrice + 200
	default:
		price = startPrice + 300
	}
	return price
}
