package src

import (
	"crypto/rand"
)

func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}


func emptyValidator(keys map[string]string) []string{
	var emptyKeys []string

	for k, v := range keys {
		if len(v) <= 0{
			emptyKeys = append(emptyKeys,k)
		}
	}

	return emptyKeys
}