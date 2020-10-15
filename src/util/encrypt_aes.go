package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime/debug"
)

func TestAes() {
	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	key := []byte("hundsun@12345678")
	result, err := AES_Encrypt([]byte("polaris@studygolang"), key)
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	origData, err := AES_Decrypt(result, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}

func AES_Encrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding_aes(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}
func AES_EncryptBase64(origData, key []byte) string {
	result, err := AES_Encrypt(origData, key)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(result)
}

func AES_Decrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding_aes(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func AES_DecryptBase64(crypted string, key []byte) (value string, cerr error) {
	defer func() {
		if err := recover(); err != nil {
			cerr =errors.New("解密失败："+ fmt.Sprintf("%s",err))
			LogError("解密失败：", cerr, string(debug.Stack()))
		}
	}()
	value = ""

	tmp, err := base64.StdEncoding.DecodeString(crypted)
	cerr = err
	if err != nil {
		return
	}

	result, err := AES_Decrypt(tmp, key)
	cerr = err
	if err != nil {
		return
	}
	value = string(result)
	cerr = nil
	return
}
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding_aes(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding_aes(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
