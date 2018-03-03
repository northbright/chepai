package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/northbright/chepai"
	"github.com/northbright/jwthelper"
)

func getLoginID(c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie("jwt")
	switch err {
	case http.ErrNoCookie:
		return "", fmt.Errorf("no jwt found in cookie")
	case nil:
		tokenString := cookie.Value
		parser := jwthelper.NewRSASHAParser([]byte(rsaPubPEM))
		m, err := parser.Parse(tokenString)
		if err != nil {
			return "", fmt.Errorf("parser.Parse() error: %v", err)
		}

		// Convert interface{} to string
		ID, ok := m["id"].(string)
		if !ok {
			return "", fmt.Errorf("failed to convert interface{} to string")
		}
		return ID, nil
	default:
		return "", fmt.Errorf("get JWT cookie error: %v", err)
	}

}

func getTimeInfo(c *gin.Context) {
	var (
		err      error
		errMsg   string
		success  = false
		ID       string
		timeInfo chepai.TimeInfo
	)

	defer func() {
		type MyTimeInfo struct {
			BeginTime       int64
			PhaseOneEndTime int64
			PhaseTwoEndTime int64
		}

		var myTimeInfo MyTimeInfo

		if err != nil {
			errMsg = err.Error()
			log.Printf("getTimeInfo() error: %v", err)
		} else {
			myTimeInfo = MyTimeInfo{timeInfo.BeginTime.Unix(),
				timeInfo.PhaseOneEndTime.Unix(),
				timeInfo.PhaseTwoEndTime.Unix(),
			}
		}

		c.JSON(200, gin.H{"success": success, "err": errMsg, "id": ID, "time_info": myTimeInfo})
	}()

	if ID, err = getLoginID(c); err != nil {
		log.Printf("getLoginID() error: %v", ID)
		return
	}

	timeInfo = cp.GetTimeInfo()

	success = true
	log.Printf("getTimeInfo() OK, ID: %v, time info: %v, %v, %v", ID, timeInfo.BeginTime, timeInfo.PhaseOneEndTime, timeInfo.PhaseTwoEndTime)
}

func validLogin(ID, password string) bool {
	if ID == "" || password == "" {
		return false
	}

	if ID != password {
		return false
	}
	return true
}

func loginPOST(c *gin.Context) {
	type Req struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	var (
		err     error
		errMsg  string
		success = false
		r       Req
	)

	defer func() {
		if err != nil {
			errMsg = err.Error()
			log.Printf("LoginPOST() error: %v", err)
		}

		c.JSON(200, gin.H{"success": success, "err": errMsg})
	}()

	if err = c.BindJSON(&r); err != nil {
		err = fmt.Errorf("invalid request")
		return
	}

	if !validLogin(r.ID, r.Password) {
		err = fmt.Errorf("incorrect password")
		return
	}

	signer := jwthelper.NewRSASHASigner([]byte(rsaPrivPEM))
	tokenString, err := signer.SignedString(
		jwthelper.NewClaim("id", r.ID),
	)
	if err != nil {
		return
	}
	cookie := jwthelper.NewCookie(tokenString)
	http.SetCookie(c.Writer, cookie)

	success = true
	log.Printf("LoginPOST() OK: ID: %v", r.ID)
}

func getStartPrice(c *gin.Context) {
	var (
		err        error
		errMsg     string
		success    = false
		ID         string
		startPrice int64
	)

	defer func() {
		if err != nil {
			errMsg = err.Error()
			log.Printf("getStartPrice() error: %v", err)
		}

		c.JSON(200, gin.H{"success": success, "err": errMsg, "id": ID, "start_price": startPrice})
	}()

	if ID, err = getLoginID(c); err != nil {
		log.Printf("getLoginID() error: %v", ID)
		return
	}

	startPrice = cp.StartPrice
	success = true
	log.Printf("getStartPrice() OK, ID: %v, start price: %v", ID, startPrice)
}

func getLicensePlateNum(c *gin.Context) {
	var (
		err             error
		errMsg          string
		success         = false
		ID              string
		licensePlateNum int64
	)

	defer func() {
		if err != nil {
			errMsg = err.Error()
			log.Printf("getLicensePlateNum() error: %v", err)
		}

		c.JSON(200, gin.H{"success": success, "err": errMsg, "id": ID, "license_plate_num": licensePlateNum})
	}()

	if ID, err = getLoginID(c); err != nil {
		log.Printf("getLoginID() error: %v", ID)
		return
	}

	licensePlateNum = cp.LicensePlateNum
	success = true
	log.Printf("getLicensePlateNum() OK, ID: %v, license plate num: %v", ID, licensePlateNum)
}

func getLowestPrice(c *gin.Context) {
	var (
		err         error
		errMsg      string
		success     = false
		ID          string
		phase       int
		lowestPrice int64
	)

	defer func() {
		if err != nil {
			errMsg = err.Error()
			log.Printf("getLowestPrice() error: %v", err)
		}

		c.JSON(200, gin.H{"success": success, "err": errMsg, "id": ID, "phase": phase, "lowest_price": lowestPrice})
	}()

	if ID, err = getLoginID(c); err != nil {
		log.Printf("getLoginID() error: %v", ID)
		return
	}

	phase = cp.GetPhase(time.Now())
	if lowestPrice, err = cp.ComupteLowestPrice(phase); err != nil {
		return
	}

	success = true
	log.Printf("getLowestPrice() OK, ID: %v, phase: %v, lowest price: %v", ID, phase, lowestPrice)
}

func bid(c *gin.Context) {
	type Req struct {
		Price int64 `json:"price"`
	}

	var (
		err     error
		errMsg  string
		success = false
		r       Req
		ID      string
		phase   int
		price   int64
	)

	defer func() {
		if err != nil {
			errMsg = err.Error()
			log.Printf("bid() error: %v", err)
		}

		c.JSON(200, gin.H{"success": success, "err": errMsg, "id": ID, "phase": phase, "price": price})
	}()

	if err = c.BindJSON(&r); err != nil {
		err = fmt.Errorf("invalid request")
		return
	}

	price = r.Price

	if ID, err = getLoginID(c); err != nil {
		log.Printf("getLoginID() error: %v", ID)
		return
	}

	phase = cp.GetPhase(time.Now())

	if err = cp.Bid(ID, price); err != nil {
		return
	}

	success = true
	log.Printf("bid() OK: ID: %v, phase: %v, price: %v", ID, phase, r.Price)
}
