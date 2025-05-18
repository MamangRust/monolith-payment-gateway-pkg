package apikey

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateApiKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(key)
}
