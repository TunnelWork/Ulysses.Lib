package api

import (
	"github.com/gin-gonic/gin"
)

// API categories. Any package (no matter official or third party) except for main package must
// choose a category to register the API endpoints.
const (
	Assets uint8 = iota
	Auth
	Billing
	Internal
	Payment
	PaymentCallback
	Plugin
	Server
)

var availableCategories map[uint8]string = map[uint8]string{
	Assets:          "assets/",
	Auth:            "auth/",
	Billing:         "billing/",
	Internal:        "internal/",
	Payment:         "payment/",
	PaymentCallback: "payment/callback/",
	Plugin:          "plugin/",
	Server:          "server/",
}

func AuthedCGET(category uint8, relativePath, userGroup string, handler ...interface{}) error {
	acFuncs, acErr := getAccessControlFunc(userGroup)
	if acErr != nil {
		return acErr
	}

	// convert everything in acFuncs to interface{}
	var ac []interface{}
	for _, v := range acFuncs {
		ac = append(ac, v)
	}

	return CGET(category, relativePath, append(ac, handler...)...)
}

func AuthedCPOST(category uint8, relativePath, userGroup string, handler ...interface{}) error {
	acFuncs, acErr := getAccessControlFunc(userGroup)
	if acErr != nil {
		return acErr
	}

	// convert everything in acFuncs to interface{}
	var ac []interface{}
	for _, v := range acFuncs {
		ac = append(ac, v)
	}

	return CGET(category, relativePath, append(ac, handler...)...)
}

// CGET() stands for Categorized GET
// CGET(Payment, "dummy/test", f) will register f() as example.com/api/payment/dummy/test for GET method
// Not validating the authentication header.
func CGET(category uint8, relativePath string, handler ...interface{}) error {
	// convert everything in handler to *gin.HandlerFunc
	var h []*gin.HandlerFunc
	for _, v := range handler {
		tmp, ok := v.(*gin.HandlerFunc)
		if !ok {
			return ErrInvalidHandler
		}
		h = append(h, tmp)
	}

	if category, exist := availableCategories[category]; exist {
		return get(category+relativePath, h...)
	} else {
		return ErrInvalidCategory
	}
}

// CPOST() stands for Categorized POST
// CPOST(Payment, "dummy/test", f) will register f() as example.com/api/payment/dummy/test for POST method
// Not validating the authentication header.
func CPOST(category uint8, relativePath string, handler ...interface{}) error {
	// convert everything in handler to *gin.HandlerFunc
	var h []*gin.HandlerFunc
	for _, v := range handler {
		tmp, ok := v.(*gin.HandlerFunc)
		if !ok {
			return ErrInvalidHandler
		}
		h = append(h, tmp)
	}

	if category, exist := availableCategories[category]; exist {
		return post(category+relativePath, h...)
	} else {
		return ErrInvalidCategory
	}
}
