package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FinalizeGinEngine() binds all registered handlers to a gin.Engine.
// - pathPrefix should be in the format of aaa[/bbb[/ccc[/ddd]]] where [] encloses an optional portion
func FinalizeGinEngine(router *gin.Engine, pathPrefix string) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	if len(pathPrefix) > 0 {
		// Safety measure: trim-off all leading/ending slashes (/)
		for pathPrefix[0] == '/' {
			pathPrefix = pathPrefix[1:] // trim off the first character
		}
		for pathPrefix[len(pathPrefix)-1] == '/' {
			pathPrefix = pathPrefix[:len(pathPrefix)-1] // trim off the last character
		}
	}

	// For non empty pathPrefix, append ending slash to make it a path.
	if len(pathPrefix) > 0 {
		pathPrefix = pathPrefix + "/"
	}

	for path, handlers := range mapGet {
		sliceHandler := []gin.HandlerFunc{}
		for _, handler := range handlers {
			sliceHandler = append(sliceHandler, *handler)
		}
		router.GET(pathPrefix+path, sliceHandler...)
	}

	for path, handlers := range mapPost {
		sliceHandler := []gin.HandlerFunc{}
		for _, handler := range handlers {
			sliceHandler = append(sliceHandler, *handler)
		}
		router.POST(pathPrefix+path, sliceHandler...)
	}

	// TODO: Clean up
	router.GET(pathPrefix+"internal/response", func(c *gin.Context) {
		cmd := c.Query("cmd")
		switch cmd {
		case "listAllMsg":
			c.JSON(http.StatusOK, payloadResponseListAllMsg())
		}
	})
}
