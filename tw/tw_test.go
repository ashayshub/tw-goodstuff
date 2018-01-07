package tw_test

import (
	"testing"

	"github.com/ashayshub/tw-goodstuff/tw"
)

func TestLoadConfig(t *testing.T) {
	a := &tw.TwApp{
		"",
		"",
		"",
		"",
		"",
	}
	if err := a.LoadConfig(); err != nil {
		t.Fail()
	}
}

func TestAuth(t *testing.T) {
	// Dummy
	// a := &tw.TwApp{"./conf.yaml.example"}
	// if err := a.LoadConfig(); err != nil {
	// 	t.Fail()
	// }
}
