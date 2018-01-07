package tw_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashayshub/tw-goodstuff/tw"
)

func TestLoadConfig(t *testing.T) {
	a := &tw.TwApp{
		"../conf.yaml.example",
		"",
		"",
		"",
		"",
		"",
	}
	if err := a.LoadConfig(); err != nil {
		log.Println(err)
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

func TestIsLoggedIn(t *testing.T) {
	a := &tw.TwApp{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8333/", nil)
	ok, err := a.IsLoggedIn(req)
	if !(err == nil || ok == true) {
		log.Println(err)
		t.Fail()
	}
}

func TestLogout(t *testing.T) {
	a := &tw.TwApp{}
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8333/", nil)
	ok, err := a.IsLoggedIn(req)
	if !(err == nil || ok == true) {
		log.Println(err)
		t.Fail()
	}
}
