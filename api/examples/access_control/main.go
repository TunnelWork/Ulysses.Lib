package main

import (
	"net/http"

	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/gin-gonic/gin"
)

var (
	userPwdMap = map[string]string{
		"user1": "pass1",
		"user2": "pass2",
	}
)

func main() {
	var userFunc gin.HandlerFunc = func(c *gin.Context) {
		userName := c.Query("user")
		pwd := c.Query("pwd")
		// Check if user exists
		if _, ok := userPwdMap[userName]; !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.MessageResponse(api.ERROR, "NOT_USER"))
		}
		if userPwdMap[userName] != pwd {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.MessageResponse(api.ERROR, "NOT_USER"))
		}
	}

	var adminFunc gin.HandlerFunc = func(c *gin.Context) {
		authHeader := c.Request.Header["Authorization"]
		// Mock Authentication criteria: only one Auth header, and header length is 32
		if len(authHeader) != 1 || len(authHeader[0]) != 32 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.MessageResponse(api.ERROR, "NOT_ADMIN"))
		}
	}

	var superAdminFunc gin.HandlerFunc = func(c *gin.Context) {
		superHeader := c.Request.Header["Superadmin"]
		if len(superHeader) != 1 || superHeader[0] != "True" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.MessageResponse(api.ERROR, "NOT_SUPER_ADMIN"))
		}
	}

	// Register the user group authentications
	api.RegisterAccessControlFuncs("user", &userFunc)            // Register the user group
	api.AppendAccessControlFuncs("admin", &userFunc, &adminFunc) // Register the admin group using AppendAccessControlFunc()
	api.InheritAccessControlFuncs("super_admin", "admin")        // Register the super admin group using InheritAccessControlFunc()
	api.AppendAccessControlFuncs("super_admin", &superAdminFunc) // Append the additional Access Control Funcs using AppendAccessControlFunc()

	// Create the router
	r := gin.Default()

	var renderSuccessPage gin.HandlerFunc = func(c *gin.Context) {
		c.JSON(http.StatusOK, api.MessageResponse(api.SUCCESS, "SUCCESS"))
	}
	api.GET("/", &renderSuccessPage)

	api.AuthedGET("/user", "user", &renderSuccessPage)
	api.AuthedGET("/admin", "admin", &renderSuccessPage)
	api.AuthedGET("/super_admin", "super_admin", &renderSuccessPage)

	// Import routes
	api.FinalizeGinEngine(r, "")

	// Run the server
	r.Run(":8081")
}
