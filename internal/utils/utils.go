package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

func UniqueId() string {
	uuId, _ := uuid.NewV7()

	u := uuId.String()
	hash := sha1.Sum([]byte(u))

	return base64.RawURLEncoding.EncodeToString(hash[:])[:8]
}

func JsonEncode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func JsonDecode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
