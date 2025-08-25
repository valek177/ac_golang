package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func generateUUID() UUID {
	// pseudo uuid
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return UUID(uuid)
}

func generateHash(request []byte) Hash {
	h := sha256.New()

	h.Write(request)

	return Hash(base64.URLEncoding.EncodeToString(h.Sum(nil)))
}
