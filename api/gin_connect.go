package api

import (
	"github.com/gin-gonic/gin"
)

// FinalizeGinEngine() binds all registered handlers to a gin.Engine.
// - pathPrefix should be in the format of aaa[/bbb[/ccc[/ddd]]] where [] encloses an optional portion
func FinalizeGinEngine(router *gin.Engine, pathPrefix string) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	// Safety measure: trim-off all leading/ending slashes (/)
	for pathPrefix[0] == '/' {
		pathPrefix = pathPrefix[1:] // trim off the first character
	}
	for pathPrefix[len(pathPrefix)-1] == '/' {
		pathPrefix = pathPrefix[:len(pathPrefix)-1] // trim off the last character
	}

	// For non empty pathPrefix, append ending slash to make it a path.
	if len(pathPrefix) > 0 {
		pathPrefix = pathPrefix + "/"
	}

	for path, handler := range mapGet {
		router.GET(pathPrefix+path, *handler)
	}

	for path, handler := range mapPost {
		router.POST(pathPrefix+path, *handler)
	}
}
