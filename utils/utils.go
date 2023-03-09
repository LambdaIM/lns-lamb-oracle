package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func Decrypt(cipherData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}
	size := block.BlockSize()
	mode := cipher.NewCBCDecrypter(block, key[:size])
	rawData := make([]byte, len(cipherData))
	mode.CryptBlocks(rawData, cipherData)
	rawData = unPadding(rawData)

	return rawData, nil
}

func unPadding(rawData []byte) []byte {
	length := len(rawData)
	end := int(rawData[length-1])
	return rawData[:(length - end)]
}

func Encrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}
	size := block.BlockSize()
	rawData = padding(rawData, size)
	mode := cipher.NewCBCEncrypter(block, key[:size])
	cipherData := make([]byte, len(rawData))
	mode.CryptBlocks(cipherData, rawData)

	return cipherData, nil
}

func padding(rawData []byte, size int) []byte {
	content := size - len(rawData)%size
	text := bytes.Repeat([]byte{byte(content)}, content)
	return append(rawData, text...)
}
