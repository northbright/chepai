package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/northbright/chepai"
	"github.com/northbright/pathhelper"
	"github.com/northbright/redishelper"
)

// Config represents the app settings.
type Config struct {
	ServerAddr       string `json:"server_addr"`
	RedisServer      string `json:"redis_server"`
	RedisPassword    string `json:"redis_password"`
	StartAfter       int    `json:"start_after"`
	PhaseOneDuration int    `json:"phase_one_duration"`
	PhaseTwoDuration int    `json:"phase_two_duration"`
	StartPrice       int64  `json:"start_price"`
	LicensePlateNum  int64  `json:"license_plate_num"`
}

var (
	currentDir, configFile string
	config                 Config
	pool                   *redis.Pool
	cp                     *chepai.Chepai
)

func main() {
	var (
		err error
	)

	defer func() {
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}()

	if err = loadConfig(); err != nil {
		err = fmt.Errorf("loadConfig() error: %v", err)
		return
	}

	// New a redis pool
	pool = redishelper.NewRedisPool(":6379", "", 1000, 100, 60, true)
	defer pool.Close()

	// New a Chepai instance
	cp = chepai.New(pool,
		config.StartAfter,
		config.PhaseOneDuration,
		config.PhaseTwoDuration,
		config.StartPrice,
		config.LicensePlateNum,
	)

	log.Printf("cp: %v", cp)
	r := gin.Default()

	// Core APIs.
	r.POST("/login", loginPOST)
	r.GET("/time_info", getTimeInfo)
	r.GET("/start_price", getStartPrice)
	r.GET("/license_plate_num", getLicensePlateNum)
	r.GET("/lowest_price", getLowestPrice)
	r.POST("/bid", bid)
	// Get student names by phone num.
	/*
		r.GET("/api/get-names-by-phone-num/:phone_num", getNamesByPhoneNum)

		// Get classes by name and phone num.
		r.GET("/api/get-classes-by-name-and-phone-num/:name/:phone_num", getClassesByNameAndPhoneNum)

		// Get available periods for the category of the class.
		r.GET("/api/get-available-periods/:class", getAvailablePeriods)

		// Post request.
		r.POST("/api/request", postRequest)
	*/
	r.Run(config.ServerAddr)
}

// init initializes path variables.
func init() {
	currentDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(currentDir, "config.json")
}

// loadConfig loads app config.
func loadConfig() error {
	// Load Conifg
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("load config file error: %v", err)

	}

	if err = json.Unmarshal(buf, &config); err != nil {
		return fmt.Errorf("parse config err: %v", err)
	}

	return nil
}
