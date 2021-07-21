package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"github.com/YeHeng/go-web-api/pkg/config"
)

const (
	defaultPassword = "123456"
)

func GeneratePassword(str string) (password string) {
	// md5
	m := md5.New()
	m.Write([]byte(str))
	mByte := m.Sum(nil)

	// hmac
	h := hmac.New(sha256.New, []byte(config.Get().Secure.Salt))
	h.Write(mByte)
	password = hex.EncodeToString(h.Sum(nil))

	return
}

func ResetPassword() (password string) {
	m := md5.New()
	m.Write([]byte(defaultPassword))
	mStr := hex.EncodeToString(m.Sum(nil))

	password = GeneratePassword(mStr)

	return
}
