package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/TunnelWork/Ulysses.Lib/auth/mfa/webauthn"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	Username = "staging"
	Password = "staging"
	Host     = "127.0.0.1"
	Port     = 3306
	Database = "tmp"
)

func stagingDB() (*sql.DB, error) {
	driverName := "mysql"
	// dsn = fmt.Sprintf("user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local", Username, Password, Host, Port, Database)

	dsn += "&autocommit=true"

	var db *sql.DB
	var err error
	db, err = sql.Open(driverName, dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func init() {
	db, err := stagingDB()
	if err != nil {
		panic(err)
	}

	auth.Setup(db, "ulysses_")
}

type registerInit struct {
	UserID   uint64 `json:"userID" binding:"required"`
	UserName string `json:"userName" binding:"required"`
}

type registerFinish struct {
	UserID     uint64      `json:"userID" binding:"required"`
	SessionKey string      `json:"sessionKey" binding:"required"`
	Response   interface{} `json:"response" binding:"required"`
}

type loginInit struct {
	UserID uint64 `json:"userID" binding:"required"`
}

type loginFinish struct {
	UserID     uint64      `json:"userID" binding:"required"`
	SessionKey string      `json:"sessionKey" binding:"required"`
	Response   interface{} `json:"response" binding:"required"`
}

func main() {
	wa := webauthn.NewWebAuthn(map[string]string{
		"RPDisplayName": "Ulysses Example WebAuthn",
		"RPID":          "localhost",
		"RPOriginURL":   "http://localhost:8081",
		// "RPIconURL":  "http://localhost/icon.png",
	})

	router := gin.Default()
	router.StaticFile("/", "./index.html")
	router.StaticFile("/scripts.js", "./scripts.js")
	router.POST("/register/init", func(c *gin.Context) {
		var ri *registerInit = &registerInit{}
		err := c.BindJSON(&ri)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("ri.userID: %d, ri.userName: %s\n", ri.UserID, ri.UserName)
		WebAuthnSignUp, err := wa.InitSignUp(ri.UserID, ri.UserName)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, WebAuthnSignUp)
	})

	router.POST("/register/finish", func(c *gin.Context) {
		var rf *registerFinish = &registerFinish{}
		err := c.BindJSON(&rf)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		responseStr, err := json.Marshal(rf.Response)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = wa.CompleteSignUp(rf.UserID, map[string]string{
			"sessionKey": rf.SessionKey,
			"response":   string(responseStr),
		})
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"success": true,
		})
	})

	router.POST("/login/init", func(c *gin.Context) {
		var li *loginInit = &loginInit{}
		err := c.BindJSON(&li)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("li.userID: %d\n", li.UserID)
		WebAuthnLogin, err := wa.NewChallenge(li.UserID)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, WebAuthnLogin)
	})

	router.POST("/login/finish", func(c *gin.Context) {
		var lf *loginFinish = &loginFinish{}
		err := c.BindJSON(&lf)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		responseStr, err := json.Marshal(lf.Response)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = wa.SubmitChallenge(lf.UserID, map[string]string{
			"sessionKey": lf.SessionKey,
			"response":   string(responseStr),
		})
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"success": true,
		})
	})

	router.Run(":8081")
}
