package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	main "github.com/ashayshub/tw-goodstuff"
	"github.com/ashayshub/tw-goodstuff/tw"
)

func TestFavPage(t *testing.T) {
	a := &tw.TwApp{}
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8333/fav", nil)
	if ok := cr.FavPage(w, req, a); cr.Status != http.StatusOK || !ok {
		t.Fail()
	}
}

func TestRTPage(t *testing.T) {
	a := &tw.TwApp{}
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8333/rt", nil)
	if ok := cr.RTPage(w, req, a); cr.Status != http.StatusOK || !ok {
		t.Fail()
	}
}

func TestHomePage(t *testing.T) {
	a := &tw.TwApp{}
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8333/", nil)
	if ok := cr.HomePage(w, req, a); !(cr.Status == http.StatusFound || cr.Status == http.StatusOK) || !ok {
		t.Fail()
	}
}

func TestWriteHTTPResponse(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	cr.Hdr = make(http.Header)

	cr.Hdr["TestKey"] = []string{"TestValue"}
	cr.Status = 200
	cr.Body.Write([]byte("200 Ok"))

	if ok := cr.WriteHTTPResponse(w); !ok {
		t.Fail()
	}
}

func TestSendInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}

	if ok := cr.SendInternalError(w); !ok {
		t.Fail()
	}
}

func TestSendNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}

	if ok := cr.SendInternalError(w); !ok {
		t.Fail()
	}
}

func TestServeHTTP(t *testing.T) {
	hd := main.Handler{}
	for _, route := range main.ActiveRoute {
		req := httptest.NewRequest(http.MethodGet, main.EndPoint+route, nil)
		resp := httptest.NewRecorder()
		hd.ServeHTTP(resp, req)
		if !(resp.Code == http.StatusOK || resp.Code == http.StatusFound) {
			fmt.Printf("%v: Status: %v, Body: %v, Result: %v\n", route, resp.Code, resp.Body.String(), resp.Result())
			t.Fail()
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
