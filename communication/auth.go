package gmqnet

import (
	"crypto/sha512"
	"strings"
)

func verifyAuth(inuser, inpwd, user, password string) bool {
	return strings.EqualFold(user, inuser) && strings.EqualFold(password, inpwd)
}

func GenToken(in string) string {
	return string(sha512.New().Sum([]byte(in)))
}
