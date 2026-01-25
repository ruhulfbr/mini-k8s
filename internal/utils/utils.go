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

func NormalizeEventMetadata(metadata any) string {
	if metadata == nil {
		return ""
	}

	switch v := metadata.(type) {

	case error:
		return MustJSON(map[string]any{
			"type":    "error",
			"message": v.Error(),
		})

	case string:
		return MustJSON(map[string]any{
			"type":  "string",
			"value": v,
		})

	default:
		return MustJSON(map[string]any{
			"type":  "object",
			"value": v,
		})
	}
}

func MustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return `{"type":"invalid","error":"metadata marshal failed"}`
	}
	return string(b)
}
