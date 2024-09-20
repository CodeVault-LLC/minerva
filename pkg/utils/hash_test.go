package utils

import (
	"testing"
)

func TestSHA256(t *testing.T) {
	data := "hello world"
	hash := SHA256(data)

	if hash != "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9" {
		t.Error("Expected SHA256 hash of 'hello world'")
	}
}

func TestSHA1(t *testing.T) {
	data := "hello world"
	hash := SHA1(data)

	if hash != "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed" {
		t.Error("Expected SHA1 hash of 'hello world'")
	}
}

func TestMD5(t *testing.T) {
	data := "hello world"
	hash := MD5(data)

	if hash != "5eb63bbbe01eeed093cb22bb8f5acdc3" {
		t.Error("Expected MD5 hash of 'hello world'")
	}
}
