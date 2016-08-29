package gmqnet

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

func verifyAuth(inuser, inpwd, user, password string) bool {
	return strings.EqualFold(user, inuser) && strings.EqualFold(password, inpwd)
}

func GenSha512Token(in string) string {
	return hex.EncodeToString(sha512.New().Sum([]byte(in)))
}
