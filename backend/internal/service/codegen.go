package service

import (
	"crypto/rand"
	"encoding/hex"
)

const inviteCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandomInviteCode(n int) (string, error) {
	b := make([]byte, n)
	for i := range b {
		v := make([]byte, 1)
		if _, err := rand.Read(v); err != nil {
			return "", err
		}
		b[i] = inviteCharset[int(v[0])%len(inviteCharset)]
	}
	return string(b), nil
}

func RandomRedeemCode() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
