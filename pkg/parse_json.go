package pkg

import (
	"encoding/json"
)

func ParseToJson(data any, errMsg string) []byte {
	json, err := json.Marshal(data)
	HandleError(errMsg, err, 1)

	return json
}

func ParseFromJson(body []byte, data any, errMsg string) {
	err := json.Unmarshal(body, data)

	HandleError(errMsg, err, 1)
}
