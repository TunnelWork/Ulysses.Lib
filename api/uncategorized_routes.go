package api

import (
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthedGET(relativePath, userGroup string, handler ...interface{}) error {
	acFuncs, acErr := getAccessControlFunc(userGroup)
	if acErr != nil {
		return acErr
	}

	// convert everything in handler to *gin.HandlerFunc
	var h []*gin.HandlerFunc
	for _, v := range handler {
		tmp, ok := v.(*gin.HandlerFunc)
		if !ok {
			return ErrInvalidHandler
		}
		h = append(h, tmp)
	}

	return get(relativePath, append(acFuncs, h...)...)
}

func AuthedPOST(relativePath, userGroup string, handler ...interface{}) error {
	acFuncs, acErr := getAccessControlFunc(userGroup)
	if acErr != nil {
		return acErr
	}

	// convert everything in handler to *gin.HandlerFunc
	var h []*gin.HandlerFunc
	for _, v := range handler {
		tmp, ok := v.(*gin.HandlerFunc)
		if !ok {
			return ErrInvalidHandler
		}
		h = append(h, tmp)
	}

	return post(relativePath, append(acFuncs, h...)...)
}

// GET() is effectively like gin.Engine.GET()
// security measure: only main package can call GET(). For modules, refer to CGET()
// Warning: this function by-default registers routes which have no AUTHORIZATION. Use AuthedGET() instead!
func GET(relativePath string, handler ...interface{}) error {
	var packageName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		packageName = strings.Split(details.Name(), ".")[0]
	}

	// convert everything in handler to *gin.HandlerFunc
	var h []*gin.HandlerFunc
	for _, v := range handler {
		tmp, ok := v.(*gin.HandlerFunc)
		if !ok {
			return ErrInvalidHandler
		}
		h = append(h, tmp)
	}

	if packageName == "main" {
		return get(relativePath, h...)
	} else {
		return ErrNotAllowDirectFuncReg
	}
}

// POST() is effectively like gin.Engine.POST(), but takes 1 handler function only
// security measure: only main package can call POST(). For modules, refer to CPOST()
// Warning: this function by-default registers routes which have no AUTHORIZATION. Use AuthedPOST() instead!
func POST(relativePath string, handler ...interface{}) error {
	var packageName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		packageName = strings.Split(details.Name(), ".")[0]
	}

	// convert everything in handler to *gin.HandlerFunc
	var h []*gin.HandlerFunc
	for _, v := range handler {
		tmp, ok := v.(*gin.HandlerFunc)
		if !ok {
			return ErrInvalidHandler
		}
		h = append(h, tmp)
	}

	if packageName == "main" {
		return post(relativePath, h...)
	} else {
		return ErrNotAllowDirectFuncReg
	}
}
