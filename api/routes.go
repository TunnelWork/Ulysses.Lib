package api

import "github.com/gin-gonic/gin"

// Warning: this function by-default registers routes which have no AUTHORIZATION
func get(relativePath string, handler ...*gin.HandlerFunc) error {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if _, conflict := mapGet[relativePath]; conflict {
		return ErrRepeatGetPath
	} else {
		mapGet[relativePath] = handler
		return nil
	}
}

// Warning: this function by-default registers routes which have no AUTHORIZATION
func post(relativePath string, handler ...*gin.HandlerFunc) error {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if _, conflict := mapPost[relativePath]; conflict {
		return ErrRepeatGetPath
	} else {
		mapPost[relativePath] = handler
		return nil
	}
}
