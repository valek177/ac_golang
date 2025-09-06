package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func generateIdFromUrl(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}
