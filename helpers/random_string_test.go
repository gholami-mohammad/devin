package helpers

import (
	"strings"
	"testing"
)

func TestRandomString(t *testing.T) {
	s1 := RandomString(32)
	s2 := RandomString(32)
	t.Log(s1, s2)
	if strings.EqualFold(s1, s2) {
		t.Fatal("Bad random generated")
	}
}
