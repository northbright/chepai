package main

import (
	"fmt"
	"log"
	"net/http"

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

		log.Printf("m: %v", m)

		// Convert interface{} to string
		ID, ok := m["ID"].(string)
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

		c.JSON(200, gin.H{"success": success, "err": errMsg, "ID": ID, "time_info": myTimeInfo})
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
		jwthelper.NewClaim("ID", r.ID),
	)
	if err != nil {
		return
	}
	cookie := jwthelper.NewCookie(tokenString)
	http.SetCookie(c.Writer, cookie)

	success = true
	log.Printf("LoginPOST() OK: ID: %v", r.ID)
}
