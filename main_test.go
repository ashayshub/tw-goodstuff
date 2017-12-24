package main

import (
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	u := userDetails{name: "ashay"}
	u.parseArgs()
}

func TestValidateArgs(t *testing.T) {
	u := userDetails{name: "ashay"}
	if ok := u.validateArgs(); !ok {
		t.Fail()
	}
}

func TestOAuth(t *testing.T) {
	u := userDetails{name: "ashay"}
	if _, ok := u.oAuth(); !ok {
		t.Fail()
	}
}

func BenchmarkOAuth(b *testing.B) {
	u := userDetails{name: "ashay"}
	for i := 0; i < b.N; i++ {
		u.oAuth()
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
