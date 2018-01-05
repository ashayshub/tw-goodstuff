package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ashayshub/tw-goodstuff/tw"
	yaml "gopkg.in/yaml.v2"
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
	configFile := "./conf.yaml.example"
	app := &tw.TwApp{}
	data, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}

	if err := yaml.Unmarshal(data, app); err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}

	switch req.URL.Path {
	case "/fav":
		resp, statusCode, ok := tw.FavPage(req)
		if !ok {
			sendInternalError(w)
			return
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	case "/rt":
		resp, statusCode, ok := tw.RTPage(req)
		if !ok {
			sendInternalError(w)
			return
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))

	case "/":
		resp, statusCode, ok := tw.HomePage(req)
		if !ok {
			sendInternalError(w)
			return
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
