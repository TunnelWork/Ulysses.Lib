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

// RegisterAccessControlFunc wipes all existing Access Control Funcs and replace it with the new ones
func RegisterAccessControlFuncs(userGroup string, acFuncs ...*gin.HandlerFunc) {
	mapAccessMutex.Lock()
	defer mapAccessMutex.Unlock()

	mapAccessControlFunc[userGroup] = acFuncs
}

// AppendAccessControlFunc adds new Access Control Funcs after the existing ones
func AppendAccessControlFuncs(userGroup string, acFuncs ...*gin.HandlerFunc) {
	mapAccessMutex.Lock()
	defer mapAccessMutex.Unlock()

	if _, ok := mapAccessControlFunc[userGroup]; !ok {
		mapAccessControlFunc[userGroup] = acFuncs
	}

	mapAccessControlFunc[userGroup] = append(mapAccessControlFunc[userGroup], acFuncs...)
}

// InheritAccessControlFunc inherits Access Control Funcs from parent user group
// if the userGroup exists, it appends parent's Access Control Funcs after the existing ones
func InheritAccessControlFuncs(userGroup, parentUserGroup string) error {
	mapAccessMutex.Lock()
	defer mapAccessMutex.Unlock()

	// Check if parent exists
	if _, ok := mapAccessControlFunc[parentUserGroup]; !ok {
		return ErrAccessControlFuncNotFound
	}

	if _, ok := mapAccessControlFunc[userGroup]; !ok {
		mapAccessControlFunc[userGroup] = mapAccessControlFunc[parentUserGroup]
	} else {
		mapAccessControlFunc[userGroup] = append(mapAccessControlFunc[userGroup], mapAccessControlFunc[parentUserGroup]...)
	}

	return nil
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
