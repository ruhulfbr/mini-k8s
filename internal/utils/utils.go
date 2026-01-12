package utils

import "github.com/google/uuid"

func UniqueId() string {
	uuId, _ := uuid.NewV7()
	return uuId.String()
}
