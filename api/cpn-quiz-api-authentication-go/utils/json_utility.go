package utils

import (
	"encoding/json"
)

// =>แปลง Object เป็น Json String.
func JsonSerialize(data any) (string, error) {
	json, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	return string(json), nil
}

// =>แปลง Json String เป็น Object ตาม Struct Type.
func JsonDeserialize[T any](data string, response T) (T, error) {
	buf := []byte(data)

	if err := json.Unmarshal(buf, &response); err != nil {
		return response, err
	}

	return response, nil
}

// =>แปลง Struct A เป็น  Struct B
func Transform[T any](source any, response T) (T, error) {
	convertJson, err := JsonSerialize(source)
	if err != nil {
		return response, err
	}

	transformData, err := JsonDeserialize(convertJson, response)
	if err != nil {
		return transformData, err
	}

	return transformData, nil
}
