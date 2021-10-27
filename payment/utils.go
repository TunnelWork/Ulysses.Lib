package payment

import (
	"encoding/json"
	"fmt"
)

// Pify() is used to recusrively convert a P{} into a
// recursive P: which contains only allowed types for value
// string, boolean, and P{}
func (params *P) Pify() *P {

	for key, val := range *params {
		switch val.(type) {
		case map[string]interface{}:
			copyVal := P(val.(map[string]interface{}))
			(*params)[key] = *(&copyVal).Pify()
		case string:
			continue
		case bool:
			continue
		default:
			copyVal := fmt.Sprintf("%v", val)
			(*params)[key] = copyVal
		}
	}

	return params
}

func (params *P) String() string {
	params.Pify()
	pByte, err := json.Marshal(*params)
	if err != nil {
		return ""
	}
	return string(pByte)
}

func Pify(save string) P {
	var recovered P = P{}
	err := json.Unmarshal([]byte(save), &recovered)
	if err != nil {
		return P{}
	}
	recovered.Pify()
	return recovered
}
