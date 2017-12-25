package main

import (
	"fmt"
	"net/http"
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
	fmt.Println("Listening on host:", HostAddr, "Port: ", HostPort)
	s.ListenAndServe()
}

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/fav":
		resp, statusCode, ok := FavPage(req)
		if !ok {
			sendInternalError(w)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	case "/rt":
		resp, statusCode, ok := RTPage(req)
		if !ok {
			sendInternalError(w)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	case "/":
		resp, statusCode, ok := HomePage(req)
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

func FavPage(req *http.Request) (resp string, statusCode int, ok bool) {
	if req.Method == http.MethodGet {
		return req.URL.Path, http.StatusOK, true
	}
	return "Not Found", http.StatusNotFound, true
}

func RTPage(req *http.Request) (resp string, statusCode int, ok bool) {
	if req.Method == http.MethodGet {
		return req.URL.Path, http.StatusOK, true
	}
	return "Not Found", http.StatusNotFound, true
}

func HomePage(req *http.Request) (resp string, statusCode int, ok bool) {
	if req.Method == http.MethodGet {
		return req.URL.Path, http.StatusOK, true
	}
	return "Not Found", http.StatusNotFound, true
}
