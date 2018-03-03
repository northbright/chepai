package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/northbright/pathhelper"
)

// Config represents the app settings.
type Config struct {
	ServerAddr  string `json:"server_addr"`
	BidderNum   int64  `json:"bidder_num"`
	Concurrency int64  `json:"concurrency"`
}

var (
	currentDir, configFile string
	config                 Config
)

type Session struct {
	ServerURLStr string
	ServerURL    *url.URL
	client       *http.Client
}

var ()

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

	Emu(config.ServerAddr, config.BidderNum, config.Concurrency)
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

func Emu(serverURL string, bidderNum, concurrency int64) {
	sem := make(chan struct{}, concurrency)

	for i := int64(0); i < bidderNum; i++ {
		// After first "concurrency" amount of goroutines started,
		// It'll block starting new goroutines until one running goroutine finishs.
		sem <- struct{}{}
		go EmuBid(sem, serverURL, i)

	}

	// After last goroutine is started,
	// there're still "concurrency" amount of goroutines running.
	// Make sure wait all goroutines to finish.
	for j := 0; j < cap(sem); j++ {
		sem <- struct{}{}
		log.Printf("----- j: %v\n", j)
	}
}

func EmuBid(sem chan struct{}, serverURL string, i int64) {
	defer func() { <-sem }()
	// New session
	s, err := NewSession(serverURL)
	if err != nil {
		log.Printf("NewSession() error: %v", err)
		return
	}

	// Login
	ID := strconv.FormatInt(i, 10)
	if err = s.Login(ID, ID); err != nil {
		log.Printf("Login() error: %v", err)
		return
	}

	// Get time info
	info, err := s.GetTimeInfo()
	if err != nil {
		log.Printf("GetTimeInfo() error: %v", err)
		return
	}

	log.Printf("time info: %v", info)

	log.Printf("i: %v\n", i)
}
