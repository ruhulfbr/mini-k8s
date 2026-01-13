package utils

import (
	"encoding/json"

	"github.com/google/uuid"
)

func UniqueId() string {
	uuId, _ := uuid.NewV7()
	return uuId.String()
}

func JsonEncode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func JsonDecode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
