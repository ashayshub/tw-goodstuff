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
		if resp.Code != http.StatusOK {
			fmt.Printf("%v: Status: %v, Body: %v, Result: %v\n", route, resp.Code, resp.Body.String(), resp.Result())
			t.Fail()
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
