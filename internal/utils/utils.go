package utils

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

func UniqueId() string {
	id, _ := uuid.NewV7()
	sum := sha256.Sum256([]byte(id.String()))

	return strings.ToLower(
		base32.StdEncoding.WithPadding(base32.NoPadding).
			EncodeToString(sum[:])[:10],
	)
}

func JsonEncode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func JsonDecode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
