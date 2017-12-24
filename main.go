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

type handler struct {
}

func (h handler)  ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/fav":
		resp, ok := FavPage(req); 
		if !ok {
			sendInternalError(w)
		}
		w.Write([]byte(resp))

	case "/rt":
		resp, ok := RTPage(req); 
		if !ok {
			sendInternalError(w)
		}
		w.Write([]byte(resp))

	case "/":
		resp, ok := HomePage(req);
		if !ok {
			sendInternalError(w)
		}
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

func FavPage(req *http.Request) (resp string, ok bool) {
	if req.Method == http.MethodGet {
		return  "Fav Page: " + req.URL.Path, true
	}
	return "", false
}

func RTPage(req *http.Request) (resp string, ok bool) {
	if req.Method == http.MethodGet {
		return  "RT Page: " + req.URL.Path, true
	}
	return "", false
}

func HomePage(req *http.Request) (resp string, ok bool) {
	if req.Method == http.MethodGet {
		return  "Home Page: " + req.URL.Path, true
	}
	return "", false
}
