package chepai_test

import (
	"fmt"
	"log"
	"time"

	"github.com/northbright/chepai"
	"github.com/northbright/redishelper"
)

func ExampleChepai_GetTimeInfo() {
	pool := redishelper.NewRedisPool(":6379", "", 1000, 100, 60, true)
	defer pool.Close()

	cp := chepai.New(pool, 10, 30, 30, 83000, 10)

	timeInfo := cp.GetTimeInfo()
	log.Printf("Time Info:\nBegin Time: %v\nPhase One End Time: %v\nPhase Two End Time: %v",
		timeInfo.BeginTime.Unix(),
		timeInfo.PhaseOneEndTime.Unix(),
		timeInfo.PhaseTwoEndTime.Unix())
	// Output:
}

func ExampleChepai_GetPhase() {
	pool := redishelper.NewRedisPool(":6379", "", 1000, 100, 60, true)
	defer pool.Close()

	cp := chepai.New(pool, 10, 30, 30, 83000, 10)

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

func ExampleChepai_ComputePhaseTwoLowestPrice() {
	pool := redishelper.NewRedisPool(":6379", "", 1000, 100, 60, true)
	defer pool.Close()

	cp := chepai.New(pool, 0, 2, 5, 83000, 10)

	// Phase 1
	for i := 0; i < 10; i++ {
		ID := fmt.Sprintf("%v", i+1)
		cp.Bid(ID, 83000)
	}

	// Phase 2
	time.Sleep(time.Second * 2)
	for i := 0; i < 8; i++ {
		ID := fmt.Sprintf("%v", i+1)
		cp.Bid(ID, 82700)
	}
	cp.Bid("9", 83300)
	cp.Bid("10", 82400)

	price, err := cp.ComputePhaseTwoLowestPrice()
	if err != nil {
		log.Printf("ComputePhaseTwoLowestPrice() error: %v", err)
	}
	log.Printf("Phase Two Lowest Price: %v", price)
	// Output:
}

func ExampleChepai_ValidPhaseTwoPrice() {
	pool := redishelper.NewRedisPool(":6379", "", 1000, 100, 60, true)
	defer pool.Close()

	cp := chepai.New(pool, 0, 30, 30, 83000, 10)
	prices := []int64{82700, 82800, 82900, 83000, 83100, 83200, 83300, 83301, 82699, 84400}

	lowestPrice, err := cp.ComputePhaseTwoLowestPrice()
	if err != nil {
		log.Printf("ComputePhaseTwoLowestPrice() error: %v", err)
	}

	for _, price := range prices {
		valid := cp.ValidPhaseTwoPrice(lowestPrice, price)
		log.Printf("%v: %v", price, valid)
	}
	log.Printf("xx: Phase Two Lowest Price: %v", lowestPrice)
	// Output:
}

func ExampleChepai_Bid() {
	var price int64

	pool := redishelper.NewRedisPool(":6379", "", 1000, 100, 60, true)
	defer pool.Close()

	cp := chepai.New(pool, 1, 5, 5, 83000, 10)

	price = 83000
	log.Printf("1st bid: price: %v, %v", price, cp.Bid("1", price))

	time.Sleep(time.Second * 2)
	price = 83400
	log.Printf("2nd bid: price: %v, %v", price, cp.Bid("1", price))

	price = 83000
	log.Printf("3rd bid: price: %v, %v", price, cp.Bid("1", price))

	price = 83000
	log.Printf("4th bid: price: %v, %v", price, cp.Bid("1", price))

	time.Sleep(time.Second * 5)
	price = 83001
	log.Printf("5th bid: price: %v, %v", price, cp.Bid("1", price))

	price = 83100
	log.Printf("6th bid: price: %v, %v", price, cp.Bid("1", price))

	price = 83200
	log.Printf("7th bid: price: %v, %v", price, cp.Bid("1", price))

	price = 83200
	log.Printf("8th bid for ID 2: price: %v, %v", price, cp.Bid("2", price))

	time.Sleep(time.Second * 5)
	price = 83300
	log.Printf("5th bid: price: %v, %v", price, cp.Bid("1", price))

	// Output:
}
