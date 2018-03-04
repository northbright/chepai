package main

import (
	"math/rand"
	"time"

	"github.com/northbright/random"
)

func generatePhaseOneSleepTime(beginTime, endTime time.Time) (time.Duration, error) {
	pad := time.Millisecond * 100
	wakeupTime, err := random.RandTime(beginTime, endTime, pad)
	if err != nil {
		return 0, err
	}
	return wakeupTime.Sub(time.Now()), nil
}

func generatePhaseTwoSleepTime(beginTime, endTime time.Time) (time.Duration, error) {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(100)

	d := 0
	switch {
	case x < 60: // 60%: wakeup at end time - 3 seconds
		d = -3
	case x >= 60 && x < 80: // 20%: wakeup at end time - 5 seconds
		d = -5
	default: // 10%: wakeup at end time - 8 seconds
		d = -8
	}

	min := endTime.Add(time.Duration(d) * time.Second)
	max := endTime
	pad := time.Millisecond * 500

	wakeupTime, err := random.RandTime(min, max, pad)
	if err != nil {
		return 0, err
	}
	return wakeupTime.Sub(time.Now()), nil
}

func generatePhaseTwoPrice(startPrice int64) int64 {
	var price int64

	t := time.Now()

	rand.Seed(t.UnixNano())
	x := rand.Intn(100)

	switch {
	case x < 10:
		price = startPrice
	case x >= 10 && x < 30:
		price = startPrice + 100
	case x >= 30 && x < 65:
		price = startPrice + 200
	default:
		price = startPrice + 300
	}
	return price
}
