package api

import (
	"github.com/gin-gonic/gin"
)

// API categories. Any package (no matter official or third party) except for main package must
// choose a category to register the API endpoints.
const (
	Assets uint8 = iota
	Internal
	Payment
	PaymentCallback
	Plugin
	Server
)

var availableCategories map[uint8]string = map[uint8]string{
	Assets:          "assets/",
	Internal:        "internal/",
	Payment:         "payment/",
	PaymentCallback: "payment/callback/",
	Plugin:          "plugin/",
	Server:          "server/",
}

func AuthedCGET(category uint8, relativePath string, userGroup string, handler ...*gin.HandlerFunc) error {
	acFuncs, acErr := getAccessControlFunc(userGroup)
	if acErr != nil {
		return acErr
	}

	return CGET(category, relativePath, append(acFuncs, handler...)...)
}

func AuthedCPOST(category uint8, relativePath string, userGroup string, handler ...*gin.HandlerFunc) error {
	acFuncs, acErr := getAccessControlFunc(userGroup)
	if acErr != nil {
		return acErr
	}

	return CPOST(category, relativePath, append(acFuncs, handler...)...)
}

// CGET() stands for Categorized GET
// CGET(Payment, "dummy/test", f) will register f() as example.com/api/payment/dummy/test for GET method
func CGET(category uint8, relativePath string, handler ...*gin.HandlerFunc) error {
	if category, exist := availableCategories[category]; exist {
		return get(category+relativePath, handler...)
	} else {
		return ErrInvalidCategory
	}
}

// CPOST() stands for Categorized POST
// CPOST(Payment, "dummy/test", f) will register f() as example.com/api/payment/dummy/test for POST method
func CPOST(category uint8, relativePath string, handler ...*gin.HandlerFunc) error {
	if category, exist := availableCategories[category]; exist {
		return post(category+relativePath, handler...)
	} else {
		return ErrInvalidCategory
	}
}
