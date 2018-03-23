package crypto

import (
	"strings"
	"testing"
)

func TestCBCEncrypter(t *testing.T) {
	t.Run("short string", func(t *testing.T) {
		str, e := CBCEncrypter(`mgh`)
		if e != nil {
			t.Fatal(e)
		}
		t.Log(str)
	})
	t.Run("long string", func(t *testing.T) {
		str, e := CBCEncrypter(`{"id":103,"username":"success_token","email":"success_token@gmail.com"}`)
		if e != nil {
			t.Fatal(e)
		}
		t.Log(str)
	})
	t.Run("long string", func(t *testing.T) {
		str, e := CBCEncrypter(`{"id":103,"username":"success_token","email":"success"}`)
		if e != nil {
			t.Fatal(e)
		}
		t.Log(str)
	})
}

func TestCBCDecrypter(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		input := `mgh`
		str, _ := CBCEncrypter(input)
		ret, e := CBCDecrypter(str)
		if e != nil {
			t.Fatal(e)
		}

		if !strings.EqualFold(ret, input) {
			t.Fatal("Decrypted string not match, retured string is: ", ret)
		}
		t.Log("Decoded string is: ", ret)
	})
	t.Run("2", func(t *testing.T) {
		input := `{"id":103,"username":"success_token","email":"success_token@gmail.com"}`
		str, _ := CBCEncrypter(input)
		ret, e := CBCDecrypter(str)
		if e != nil {
			t.Fatal(e)
		}

		if !strings.EqualFold(ret, input) {
			t.Fatal("Decrypted string not match, retured string is: ", ret)
		}
		t.Log("Decoded string is: ", ret)
	})
	t.Run("3", func(t *testing.T) {
		input := `{"id":103,"username":"success_token","email":"success"}`
		str, _ := CBCEncrypter(input)
		ret, e := CBCDecrypter(str)
		if e != nil {
			t.Fatal(e)
		}

		if !strings.EqualFold(ret, input) {
			t.Fatal("Decrypted string not match, retured string is: ", ret)
		}
		t.Log("Decoded string is: ", ret)
	})

}

// func TestCBCDecrypter_CustomInput(t *testing.T) {
// 	ret, e := CBCDecrypter("3fdceb6f72b71597c26e301c91be6a55e59c93ea11829b8eec79b2a93a07b9ce")
//
// 	t.Log(ret, e)
// }
