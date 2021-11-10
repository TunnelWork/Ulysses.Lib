package main

import (
	"database/sql"
	"fmt"

	"github.com/TunnelWork/Ulysses.Lib/auth"
	utotp "github.com/TunnelWork/Ulysses.Lib/auth/mfa/totp"
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
	UserID uint64 `json:"userID" binding:"required"`
	Secret string `json:"secret" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

// type loginInit struct {
// 	UserID uint64 `json:"userID" binding:"required"`
// }

type loginFinish struct {
	UserID uint64 `json:"userID" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

func main() {
	wa := utotp.NewTOTP(map[string]string{
		"issuer": "Ulysses Example TOTP",
	})

	router := gin.Default()
	router.StaticFile("/", "./index.html")
	router.StaticFile("/scripts.js", "./scripts.js")
	router.StaticFile("/qrcode.min.js", "./qrcode.min.js")
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

		err = wa.CompleteSignUp(rf.UserID, map[string]string{
			"secret": rf.Secret,
			"code":   rf.Code,
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

	router.POST("/login/finish", func(c *gin.Context) {
		var lf *loginFinish = &loginFinish{}
		err := c.BindJSON(&lf)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = wa.SubmitChallenge(lf.UserID, map[string]string{
			"code": lf.Code,
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
