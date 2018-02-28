package chepai_test

import (
	"log"

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
