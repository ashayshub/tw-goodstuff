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

func main() {
	var hd = handler{}
	var ar = [3]string {"/fav", "/rt", "/"}

	fmt.Println(HostAddr, HostPort)
	for _, route := range ar {
		http.Handle(route, hd)
	}
	http.ListenAndServe(HostAddr + ":" + HostPort, nil)
}

type handler struct {}

func (h handler)  ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/fav":
		resp, statusCode, ok := FavPage(req); 
		if !ok {
			sendInternalError(w)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	case "/rt":
		resp, statusCode, ok := RTPage(req); 
		if !ok {
			sendInternalError(w)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	case "/":
		resp, statusCode, ok := HomePage(req);
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

func sendInternalError(w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}

func FavPage(req *http.Request) (resp string, statusCode int, ok bool) {
	if req.Method == http.MethodGet {
		return  "Fav Page: " + req.URL.Path, http.StatusOK, true
	}
	return "Not Found", http.StatusNotFound , true
}

func RTPage(req *http.Request) (resp string, statusCode int, ok bool) {
	if req.Method == http.MethodGet {
		return  "RT Page: " + req.URL.Path, http.StatusOK, true
	}
	return "Not Found", http.StatusNotFound , true
}

func HomePage(req *http.Request) (resp string, statusCode int, ok bool) {
	if req.Method == http.MethodGet {
		return  "Home Page: " + req.URL.Path, http.StatusOK, true
	}
	return "Not Found", http.StatusNotFound , true
}
