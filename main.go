package main

import (
	"fmt"
	"net/http"

	"github.com/ashayshub/tw-goodstuff/twlib"
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
		Addr: HostAddr + ":" + HostPort,
		Handler: hd,
	}
	fmt.Printf("Starting on host: %v:%v\n", HostAddr, HostPort)
	s.ListenAndServe()
}

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/fav":
		resp, statusCode, ok := tw.FavPage(req)
		if !ok {
			sendInternalError(w)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	case "/rt":
		resp, statusCode, ok := tw.RTPage(req)
		if !ok {
			sendInternalError(w)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	case "/":
		resp, statusCode, ok := tw.HomePage(req)
		if !ok {
			sendInternalError(w)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found: " + req.URL.Path))
	}
}

func sendInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}
