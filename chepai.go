package chepai

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type TimeInfo struct {
	BeginTime       time.Time
	PhaseOneEndTime time.Time
	PhaseTwoEndTime time.Time
}

type Chepai struct {
	pool            *redis.Pool
	timeInfo        TimeInfo
	StartPrice      int64
	LicensePlateNum int64
}

func New(pool *redis.Pool, startAfter, phaseOneDuration, phaseTwoDuration int, startPrice, licensePlateNum int64) *Chepai {
	beginTime := time.Now().Add(time.Duration(startAfter) * time.Second)
	phaseOneEndTime := beginTime.Add(time.Duration(phaseOneDuration) * time.Second)
	phaseTwoEndTime := phaseOneEndTime.Add(time.Duration(phaseTwoDuration) * time.Second)

	return &Chepai{
		pool: pool,
		timeInfo: TimeInfo{
			BeginTime:       beginTime,
			PhaseOneEndTime: phaseOneEndTime,
			PhaseTwoEndTime: phaseTwoEndTime,
		},
		StartPrice:      startPrice,
		LicensePlateNum: licensePlateNum,
	}
}

func (cp *Chepai) GetTimeInfo() *TimeInfo {
	return &(cp.timeInfo)
}

func (cp *Chepai) GetPhase(t time.Time) int {
	switch {
	case t.Before(cp.timeInfo.BeginTime):
		return 0
	case t.Equal(cp.timeInfo.BeginTime) || (t.After(cp.timeInfo.BeginTime) && t.Before(cp.timeInfo.PhaseOneEndTime)):
		return 1
	case t.Equal(cp.timeInfo.PhaseOneEndTime) || (t.After(cp.timeInfo.PhaseOneEndTime) && t.Before(cp.timeInfo.PhaseTwoEndTime)):
		return 2
	default:
		return 3
	}
}

func (cp *Chepai) ComputePhaseTwoLowestPrice() (int64, error) {
	conn := cp.pool.Get()
	defer conn.Close()

	k := "prices"
	prices, err := redis.Int64s(conn.Do("ZREVRANGE", k, 0, -1))
	if err != nil {
		return 0, err
	}

	sum := int64(0)
	for _, price := range prices {
		k := fmt.Sprintf("%v:count", price)
		num, err := redis.Int64(conn.Do("GET", k))
		if err != nil && err != redis.ErrNil {
			return 0, err
		}

		if num == 0 {
			return 0, fmt.Errorf("no member matches price: %v", price)
		}

		sum += num
		if sum >= cp.LicensePlateNum {
			return price, nil
		}
	}

	l := len(prices)
	if l > 0 {
		if prices[l-1] > cp.StartPrice {
			return cp.StartPrice, nil
		} else {
			return prices[l-1], nil
		}
	}
	return cp.StartPrice, nil
}

func (cp *Chepai) ComupteLowestPrice() (int64, error) {
	conn := cp.pool.Get()
	defer conn.Close()

	phase := cp.GetPhase(time.Now())
	switch phase {
	case 0, 1:
		return cp.StartPrice, nil
	case 2, 3:
		return cp.ComputePhaseTwoLowestPrice()
	default:
		return 0, fmt.Errorf("incorrect phase: %v", phase)
	}
}

func (cp *Chepai) Bid(ID string, price int64) error {
	phase := cp.GetPhase(time.Now())
	if phase != 1 && phase != 2 {
		return fmt.Errorf("incorrect phase: %v, should be 1,2", phase)
	}

	return nil
}

func (cp *Chepai) GetParticipantNum() (int64, error) {
	return 0, nil
}
