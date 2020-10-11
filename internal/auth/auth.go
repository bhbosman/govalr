package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

type ValrAuth struct {
}

func (self ValrAuth) Hash(secretApiKey string, timeStamp string, urlPath string, verb string, body string) string {
	h := hmac.New(sha512.New, []byte(secretApiKey))
	h.Write([]byte(timeStamp))
	h.Write([]byte(verb))
	h.Write([]byte(urlPath))
	h.Write([]byte(body))
	return hex.EncodeToString(h.Sum(nil))
}
