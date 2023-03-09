package utils

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestEncrypt(t *testing.T) {
	raw := []byte("255c2ae7f027edeee07d66522c36bdda4d3")
	key := []byte("0123456789abcdef")
	pwd, err := Encrypt(raw, key)
	if nil != err {
		t.Error(err)
	}

	data := base64.StdEncoding.EncodeToString(pwd)
	t.Log(data)
	crypted, err := base64.StdEncoding.DecodeString(data)
	if nil != err {
		t.Error(err)
	}

	rawData, err := Decrypt(crypted, key)
	if nil != err {
		t.Error(err)
	}
	if !bytes.Equal(raw, rawData) {
		t.Error("not equal")
	}
}
