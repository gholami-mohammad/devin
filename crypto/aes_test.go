package crypto

import (
	"strings"
	"testing"
)

func TestCBCEncrypter(t *testing.T) {
	str, e := CBCEncrypter("MGH")
	if e != nil {
		t.Fatal(e)
	}

	t.Log(str)
}

func TestCBCDecrypter(t *testing.T) {
	input := "MGH"
	str, _ := CBCEncrypter(input)
	ret, e := CBCDecrypter(str)
	if e != nil {
		t.Fatal(e)
	}

	if !strings.EqualFold(ret, input) {
		t.Fatal("Decrypted string not match, retured string is: ", ret)
	}
	t.Log("Decoded string is: ", ret)
}
