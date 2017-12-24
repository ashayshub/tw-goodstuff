package main_test

import (
	"os"
	"testing"

	"github.com/ashayshub/tw-goodstuff"
)

func TestParseArgs(t *testing.T) {
	u := main.UserProperty{Name: "ashay"}
	u.ParseArgs()
}

func TestValidateArgs(t *testing.T) {
	u := main.UserProperty{Name: "ashay"}
	usage := u.ParseArgs()
	if ok := u.ValidateArgs(&usage); !ok {
		t.Fail()
	}
}

func TestOAuth(t *testing.T) {
	u := main.UserProperty{Name: "ashay"}
	if _, ok := u.OAuth(); !ok {
		t.Fail()
	}
}

func BenchmarkOAuth(b *testing.B) {
	u := main.UserProperty{Name: "ashay"}
	for i := 0; i < b.N; i++ {
		u.OAuth()
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
