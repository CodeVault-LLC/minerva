package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
)

func SHA256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func SHA1(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func MD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
