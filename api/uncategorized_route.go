package api

import (
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthedGET(relativePath string, userGroup string, handler ...*gin.HandlerFunc) error {
	acFuncs, acErr := getAccessControlFunc(userGroup)
	if acErr != nil {
		return acErr
	}

	return get(relativePath, append(acFuncs, handler...)...)
}

func AuthedPOST(relativePath string, userGroup string, handler ...*gin.HandlerFunc) error {
	acFuncs, acErr := getAccessControlFunc(userGroup)
	if acErr != nil {
		return acErr
	}

	return post(relativePath, append(acFuncs, handler...)...)
}

// GET() is effectively like gin.Engine.GET()
// security measure: only main package can call GET(). For modules, refer to CGET()
// Warning: this function by-default registers routes which have no AUTHORIZATION. Use AuthedGET() instead!
func GET(relativePath string, handler ...*gin.HandlerFunc) error {
	var packageName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		packageName = strings.Split(details.Name(), ".")[0]
	}

	if packageName == "main" {
		return get(relativePath, handler...)
	} else {
		return ErrNotAllowDirectFuncReg
	}
}

// POST() is effectively like gin.Engine.POST(), but takes 1 handler function only
// security measure: only main package can call POST(). For modules, refer to CPOST()
// Warning: this function by-default registers routes which have no AUTHORIZATION. Use AuthedPOST() instead!
func POST(relativePath string, handler ...*gin.HandlerFunc) error {
	var packageName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		packageName = strings.Split(details.Name(), ".")[0]
	}

	if packageName == "main" {
		return post(relativePath, handler...)
	} else {
		return ErrNotAllowDirectFuncReg
	}
}
