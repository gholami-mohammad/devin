package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"devin/crypto/keys"
)

func Init() {
	log.SetFlags(log.Lshortfile)
}

func CBCEncrypter(str string) (string, error) {
	key, _ := hex.DecodeString(keys.AES_KEY)
	plainBytes := []byte(str)
	b64 := base64.RawStdEncoding.WithPadding('=').EncodeToString(plainBytes)
	b64Bytes := []byte(b64)
	if len(b64Bytes)%aes.BlockSize != 0 {
		rpt := aes.BlockSize - (len(b64Bytes) % aes.BlockSize)
		if rpt < 0 {
			rpt *= -1
		}
		b64 += strings.Repeat("=", rpt)
	}
	plaintext := []byte(b64)
	block, e := aes.NewCipher(key)
	if e != nil {
		return "", e
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, e := io.ReadFull(rand.Reader, iv); e != nil {
		return "", e
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return fmt.Sprintf("%x", ciphertext), nil
}

func CBCDecrypter(encodedString string) (string, error) {
	key, _ := hex.DecodeString(keys.AES_KEY)
	ciphertext, e := hex.DecodeString(encodedString)
	if e != nil {
		return "", e
	}

	block, e := aes.NewCipher(key)
	if e != nil {
		return "", e
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)
	b64 := strings.Replace(string(ciphertext), "=", "", -1)

	bts, e := base64.RawStdEncoding.DecodeString(b64)
	if e != nil {
		return "", e
	}
	return string(bts), nil
}
