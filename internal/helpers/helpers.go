package helpers

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
)

func GetSHA(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	hasher := sha1.New()
	hasher.Write(b)
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
