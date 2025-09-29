package formatter

import "encoding/json"

func ToJsonString(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return `{"error":"internal: failed to marshal data"}`
	}
	return string(jsonBytes)
}
