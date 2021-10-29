package api

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mapMutex sync.RWMutex = sync.RWMutex{}

	apiGETMap  map[string]*gin.HandlerFunc = map[string]*gin.HandlerFunc{}
	apiPOSTMap map[string]*gin.HandlerFunc = map[string]*gin.HandlerFunc{}
)

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

	for path, handler := range apiGETMap {
		router.GET(pathPrefix+path, *handler)
	}

	for path, handler := range apiPOSTMap {
		router.POST(pathPrefix+path, *handler)
	}

}

// Disabled for now
// func ExportHandlerMaps() (getMap map[string]*gin.HandlerFunc, postMap map[string]*gin.HandlerFunc) {
// 	mapMutex.RLock()
// 	defer mapMutex.RUnlock()
// 	// Copy list
// 	getMap = apiGETMap
// 	postMap = apiPOSTMap

// 	return getMap, postMap
// }
