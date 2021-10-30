package api

import (
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

// GET() is effectively like gin.Engine.GET(), but takes 1 handler function only
// security measure: only main package can call GET(). For module, refer to CGET()
func GET(relativePath string, handler *gin.HandlerFunc) error {
	var packageName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		packageName = strings.Split(details.Name(), ".")[0]
	}

	if packageName == "main" {
		return get(relativePath, handler)
	} else {
		return ErrNotAllowDirectFuncReg
	}
}

// POST() is effectively like gin.Engine.POST(), but takes 1 handler function only
// security measure: only main package can call POST(). For module, refer to CPOST()
func POST(relativePath string, handler *gin.HandlerFunc) error {
	var packageName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		packageName = strings.Split(details.Name(), ".")[0]
	}

	if packageName == "main" {
		return post(relativePath, handler)
	} else {
		return ErrNotAllowDirectFuncReg
	}
}

func get(relativePath string, handler *gin.HandlerFunc) error {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if _, conflict := mapGet[relativePath]; conflict {
		return ErrRepeatGetPath
	} else {
		mapGet[relativePath] = handler
		return nil
	}
}

func post(relativePath string, handler *gin.HandlerFunc) error {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if _, conflict := mapPost[relativePath]; conflict {
		return ErrRepeatGetPath
	} else {
		mapPost[relativePath] = handler
		return nil
	}
}
