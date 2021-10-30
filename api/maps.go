package api

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mapMutex sync.RWMutex                = sync.RWMutex{}
	mapPost  map[string]*gin.HandlerFunc = map[string]*gin.HandlerFunc{}
	mapGet   map[string]*gin.HandlerFunc = map[string]*gin.HandlerFunc{}
)
