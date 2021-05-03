package helpers

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func GetSHA(value interface{}) (string, error) {
	if value == nil {
		return "", fmt.Errorf("Cannot generate SHA from nil")
	}
	b1, err1 := json.Marshal(value)
	if err1 != nil {
		return "", err1
	}
	b2, err2 := json.Marshal(fmt.Sprintf("%t", value))
	if err2 != nil {
		return "", err2
	}
	hasher := sha1.New()
	hasher.Write(b1)
	hasher.Write(b2)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha, nil
}

/* Returns a valid index given the length of the list. */
func ValidIndex(i int, length int) int {
	if i < 0 {
		return 0
	} 
	if i > length {
		return length
	}
	return i
}
