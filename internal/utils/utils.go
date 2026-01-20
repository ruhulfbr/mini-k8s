package utils

import (
	"crypto/sha256"
	"encoding/base32"
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
