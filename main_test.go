package main_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	main "github.com/ashayshub/tw-goodstuff"
	"github.com/ashayshub/tw-goodstuff/tw"
)

var (
	a = &tw.TwApp{
		"./conf.yaml.example",
		"",
		"",
		"",
		"",
		"",
	}
)

func TestFavPage(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8333/fav", nil)
	if ok := cr.FavPage(w, req, a); cr.Status != http.StatusOK || !ok {
		log.Println("Test failed")
		t.Fail()
	}
}

func TestRTPage(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8333/rt", nil)
	if ok := cr.RTPage(w, req, a); cr.Status != http.StatusOK || !ok {
		log.Println("Test failed")
		t.Fail()
	}
}

func TestHomePage(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	cr.TmplFile = "./tmpl/home.tmpl"
	cr.Hdr = w.Header()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8333/", nil)
	if ok := cr.HomePage(w, req, a); !(cr.Status == http.StatusFound || cr.Status == http.StatusOK) || !ok {
		log.Println("Test failed")
		t.Fail()
	}
}

func TestLoginPage(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	cr.Hdr = w.Header()
	a.LoadConfig()
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8333/login", nil)
	if ok := cr.LoginPage(w, req, a); !(cr.Status == http.StatusOK || cr.Status == http.StatusFound) || !ok {
		log.Println("Test failed")
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
		log.Println("Test failed")
		t.Fail()
	}
}

func TestSendRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}
	cr.Hdr = w.Header()

	if ok := cr.SendRedirect(w, "/rt"); !ok {
		log.Println("Test failed")
		t.Fail()
	}
}

func TestSendInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}

	if ok := cr.SendInternalError(w); ok {
		log.Println("Error: SendInternalError Should always fail")
		t.Fail()
	}
}

func TestSendNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	cr := &main.ContentResponse{}

	if ok := cr.SendNotFound(w, "/notFound"); !ok {
		log.Println("Test failed")
		t.Fail()
	}
}

func TestParseTmpl(t *testing.T) {
	data := []byte("Hello Testing")
	cr := &main.ContentResponse{}
	cr.Hdr = make(http.Header)

	if err := cr.ParseTmpl(data, nil); err != nil {
		log.Println(err)
		t.Fail()
	}
}

func TestReadParseTmpl(t *testing.T) {
	cr := &main.ContentResponse{}
	cr.Hdr = make(http.Header)
	cr.TmplFile = "./tmpl/home.tmpl"
	if err := cr.ReadParseTmpl(); err != nil {
		log.Println(err)
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
