package chepai_test

import (
	"log"
	"time"

	"github.com/northbright/chepai"
	"github.com/northbright/redishelper"
)

func ExampleGetTimeInfo() {
	pool := redishelper.NewRedisPool(":6379", "", 1000, 100, 60, true)
	cp := chepai.New(pool, 10, 30, 30, 83000, 10000)

	timeInfo := cp.GetTimeInfo()
	log.Printf("Time Info:\nBegin Time: %v\nPhase One End Time: %v\nPhase Two End Time: %v",
		timeInfo.BeginTime.Unix(),
		timeInfo.PhaseOneEndTime.Unix(),
		timeInfo.PhaseTwoEndTime.Unix())
	// Output:
}

func ExampleGetPhase() {
	pool := redishelper.NewRedisPool(":6379", "", 1000, 100, 60, true)
	cp := chepai.New(pool, 10, 30, 30, 83000, 10000)

	t := time.Now()
	times := []time.Time{
		t,
		t.Add(10 * time.Second),
		t.Add(15 * time.Second),
		t.Add(30 * time.Second),
		t.Add(45 * time.Second),
		t.Add(70 * time.Second),
		t.Add(80 * time.Second),
	}

	for _, t := range times {
		phase := cp.GetPhase(t)
		log.Printf("Phase: %v\n", phase)
	}
	// Output:
}
