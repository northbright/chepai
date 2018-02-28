package chepai

import (
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
	return 0
}

func (cp *Chepai) Bid(ID string, price int64) error {
	return nil
}

func (cp *Chepai) GetParticipantNum() (int64, error) {
	return 0, nil
}

func (cp *Chepai) ComupteLowestPrice() (int64, error) {
	return 0, nil
}
