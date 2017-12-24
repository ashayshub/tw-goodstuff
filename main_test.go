package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ashayshub/tw-goodstuff"
)

func TestServeHTTP(t *testing.T) {
	hd := main.Handler{}
	for _, route := range main.ActiveRoute {
		req := httptest.NewRequest(http.MethodGet, main.EndPoint+route, nil)
		resp := httptest.NewRecorder()

		hd.ServeHTTP(resp, req)
		resp.Flush()

		fmt.Println(resp.Result())
		if resp.Code != http.StatusOK {
			t.Fail()
		}
	}
}

func TestFavPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, main.EndPoint+"/fav", nil)
	_, status, ok := main.FavPage(req)
	if status != http.StatusOK || !ok {
		t.Fail()
	}
}

func TestRTPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, main.EndPoint+"/rt", nil)
	_, status, ok := main.RTPage(req)
	if status != http.StatusOK || !ok {
		t.Fail()
	}
}

func TestHomePage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, main.EndPoint+"/rt", nil)
	_, status, ok := main.RTPage(req)
	if status != http.StatusOK || !ok {
		t.Fail()
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
