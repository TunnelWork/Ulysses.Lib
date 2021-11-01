package api

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type ResponseStatus uint8

const (
	CANCELED ResponseStatus = iota
	ERROR
	SUCCESS
)

var (
	respMsgMapMutex sync.RWMutex      = sync.RWMutex{}
	respMsgMap      map[string]string = map[string]string{}
)

// MessageResponse() creates a simple response with only status and message properties
// and registers the message to a map. Frontend may request the map at anytime for
// application localization/user-friendlization purposes.
// Recommended for: Create/Update/Delete actions
func MessageResponse(status ResponseStatus, message string) gin.H {
	respMsgMapMutex.Lock()
	defer respMsgMapMutex.Unlock()

	switch status {
	case CANCELED:
		respMsgMap[message] = message
		return gin.H{
			"status":  "canceled",
			"message": message,
		}
	case ERROR:
		respMsgMap[message] = message
		return gin.H{
			"status":  "error",
			"message": message,
		}
	case SUCCESS:
		respMsgMap[message] = message
		return gin.H{
			"status":  "success",
			"message": message,
		}
	default:
		return nil
	}
}

// PayloadResponse() creates a response with status and payload properties.
// Where payload could be a random typed interface.
// Recommended for: Read actions
func PayloadResponse(status ResponseStatus, payload interface{}) gin.H {
	switch status {
	case CANCELED:
		return gin.H{
			"status":  "canceled",
			"payload": payload,
		}
	case ERROR:
		return gin.H{
			"status":  "error",
			"payload": payload,
		}
	case SUCCESS:
		return gin.H{
			"status":  "success",
			"payload": payload,
		}
	default:
		return nil
	}
}

func payloadResponseListAllMsg() gin.H {
	respMsgMapMutex.RLock()
	defer respMsgMapMutex.RUnlock()
	var allMsg []string = []string{}
	for msg := range respMsgMap {
		allMsg = append(allMsg, msg)
	}

	return PayloadResponse(SUCCESS, allMsg)
}
