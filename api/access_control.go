package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mapAccessMutex       sync.RWMutex                  = sync.RWMutex{}
	mapAccessControlFunc map[string][]*gin.HandlerFunc = map[string][]*gin.HandlerFunc{}
)

func RegisterAccessControlFunc(userGroup string, acFuncs ...*gin.HandlerFunc) {
	mapAccessMutex.Lock()
	defer mapAccessMutex.Unlock()

	mapAccessControlFunc[userGroup] = acFuncs
}

func getAccessControlFunc(userGroup string) ([]*gin.HandlerFunc, error) {
	mapAccessMutex.RLock()
	defer mapAccessMutex.RUnlock()

	if acFuncs, ok := mapAccessControlFunc[userGroup]; ok {
		return acFuncs, nil
	}
	return nil, ErrUnknownUserGroup
}

func SampleAccessControlFunc(c *gin.Context) {
	authHeader := c.Request.Header["Authorization"]
	// Mock Authentication criteria: only one Auth header, and header length is 32
	if len(authHeader) != 1 || len(authHeader[0]) != 32 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, MessageResponse(ERROR, "AUTH_FAILED"))
	}
}
