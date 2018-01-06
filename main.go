package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ashayshub/tw-goodstuff/tw"
)

const (
	HostAddr string = "localhost"
	HostPort string = "8333"
	EndPoint string = "http://" + HostAddr + ":" + HostPort
)

var (
	ActiveRoute [3]string = [3]string{"/fav", "/rt", "/"}
)

func main() {
	var hd = Handler{}
	s := &http.Server{
		Addr:    HostAddr + ":" + HostPort,
		Handler: hd,
	}
	fmt.Printf("Starting on host: %v:%v\n", HostAddr, HostPort)
	s.ListenAndServe()
}

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app := &tw.TwApp{}
	app.ConfigFile = "./conf.yaml.example"
	cr := &ContentResponse{}

	//Startup errors
	if err := app.LoadConfig(); err != nil {
		log.Fatal(err)
	}
	if err := app.Auth(); err != nil {
		log.Fatal(err)
	}

	switch req.URL.Path {
	case "/fav":
		if ok := cr.FavPage(w, req); !ok {
			cr.SendInternalError(w)
			return
		}

	case "/rt":
		if ok := cr.RTPage(w, req); !ok {
			cr.SendInternalError(w)
			return
		}

	case "/":
		if ok := cr.HomePage(w, req); !ok {
			cr.SendInternalError(w)
			return
		}

	default:
		cr.SendNotFound(w, req.URL.Path)
	}
}

type ContentResponse struct {
	Status int
	Body   string
	Hdr    http.Header
}

func (cr *ContentResponse) FavPage(w http.ResponseWriter, req *http.Request) (ok bool) {
	if req.Method == http.MethodGet {
		cr.Body = req.URL.Path
		cr.Status = http.StatusOK
		cr.Hdr = w.Header()
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) RTPage(w http.ResponseWriter, req *http.Request) (ok bool) {
	if req.Method == http.MethodGet {
		cr.Body = req.URL.Path
		cr.Status = http.StatusOK
		cr.Hdr = w.Header()
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) HomePage(w http.ResponseWriter, req *http.Request) (ok bool) {
	if req.Method == http.MethodGet {
		cr.Status = 302
		cr.Hdr = w.Header()
		cr.Hdr.Add("Location", "/rt")
		cr.Body = ""
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) SendNotFound(w http.ResponseWriter, url string) (ok bool) {
	cr.Status = http.StatusNotFound
	cr.Body = "Not Found: " + url
	cr.Hdr = w.Header()
	cr.WriteHTTPResponse(w)
	return false
}

func (cr *ContentResponse) SendInternalError(w http.ResponseWriter) (ok bool) {
	cr.Status = http.StatusInternalServerError
	cr.Body = "Internal Server Error"
	cr.Hdr = w.Header()
	cr.WriteHTTPResponse(w)
	return true
}

func (cr *ContentResponse) WriteHTTPResponse(w http.ResponseWriter) (ok bool) {
	w.WriteHeader(cr.Status)
	if len(cr.Hdr) != 0 {
		if err := cr.Hdr.Write(w); err != nil {
			log.Printf("Could not write Header: %#v", cr.Hdr)
			return cr.SendInternalError(w)
		}
	}
	w.Write([]byte(cr.Body))
	return true
}
