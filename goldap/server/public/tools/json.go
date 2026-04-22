package tools

import (
	"encoding/json"
	"fmt"
)

// Struct2Json converts struct to JSON string
func Struct2Json(obj any) string {
	str, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("Struct2Json error: %v", err))
	}
	return string(str)
}

// Json2Struct converts JSON string to struct
func Json2Struct(str string, obj any) {
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		panic(fmt.Sprintf("Json2Struct error: %v", err))
	}
}

// JsonI2Struct converts JSON interface to struct
func JsonI2Struct(str any, obj any) {
	JsonStr := str.(string)
	Json2Struct(JsonStr, obj)
}

// JsonToMap converts JSON string to map
func JsonToMap(jsonStr string) (m map[string]string, err error) {
	err = json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}
	return
}

// MapToJson converts map to JSON string
func MapToJson(m map[string]string) (string, error) {
	result, err := json.Marshal(m)
	if err != nil {
		return "", nil
	}
	return string(result), nil
}
