package api

import (
	"github.com/gin-gonic/gin"
)

// API categories. Any package (no matter official or third party) except for main package must
// choose a category to register the API endpoints.
const (
	Payment uint8 = iota
	PaymentCallback
	Server
)

var availableCategories map[uint8]string = map[uint8]string{
	Payment:         "payment/",
	PaymentCallback: "payment/callback/",
	Server:          "server/",
}

// CGET() stands for Categorized GET
// CGET(Payment, "dummy/test", f) will register f() as example.com/api/payment/dummy/test for GET method
func CGET(category uint8, relativePath string, handler *gin.HandlerFunc) error {
	if category, exist := availableCategories[category]; exist {
		return get(category+relativePath, handler)
	} else {
		return ErrInvalidCategory
	}
}

// CPOST() stands for Categorized POST
// CPOST(Payment, "dummy/test", f) will register f() as example.com/api/payment/dummy/test for POST method
func CPOST(category uint8, relativePath string, handler *gin.HandlerFunc) error {
	if category, exist := availableCategories[category]; exist {
		return post(category+relativePath, handler)
	} else {
		return ErrInvalidCategory
	}
}
