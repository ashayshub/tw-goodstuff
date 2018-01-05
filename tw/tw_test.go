package tw_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashayshub/tw-goodstuff/tw"
)

func TestFavPage(t *testing.T) {
	// a := tw.TwApp{}
	req := httptest.NewRequest(http.MethodGet, "http://locahost:8333/fav", nil)
	_, status, ok := tw.FavPage(req)
	if status != http.StatusOK || !ok {
		t.Fail()
	}
}

func TestRTPage(t *testing.T) {
	// a := tw.TwApp{}
	req := httptest.NewRequest(http.MethodGet, "http://locahost:8333/rt", nil)
	_, status, ok := tw.RTPage(req)
	if status != http.StatusOK || !ok {
		t.Fail()
	}
}

func TestHomePage(t *testing.T) {
	// a := tw.TwApp{}
	req := httptest.NewRequest(http.MethodGet, "http://locahost:8333/", nil)
	_, status, ok := tw.HomePage(req)
	if status != http.StatusOK || !ok {
		t.Fail()
	}
}
