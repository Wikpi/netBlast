package pkg

import (
	"encoding/json"
	"fmt"
)

func ParseToJson(data any, errMsg string) []byte {
	json, err := json.Marshal(data)
	if err != nil {
		LogError(err)
		fmt.Println(errMsg)
	}

	return json
}

func ParseFromJson(body []byte, data any, errMsg string) {
	err := json.Unmarshal(body, data)
	if err != nil {
		LogError(err)
		fmt.Println(errMsg)
	}
}
