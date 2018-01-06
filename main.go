package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ashayshub/tw-goodstuff/tw"
	"github.com/pkg/errors"
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

	switch req.URL.Path {
	case "/fav":
		if ok := cr.FavPage(w, req, app); !ok {
			cr.SendInternalError(w)
			return
		}

	case "/rt":
		if ok := cr.RTPage(w, req, app); !ok {
			cr.SendInternalError(w)
			return
		}

	case "/":
		if ok := cr.HomePage(w, req, app); !ok {
			cr.SendInternalError(w)
			return
		}

	default:
		cr.SendNotFound(w, req.URL.Path)
	}
}

type ContentResponse struct {
	Status int
	Body   bytes.Buffer
	Hdr    http.Header
}

func (cr *ContentResponse) FavPage(w http.ResponseWriter, req *http.Request, app *tw.TwApp) (ok bool) {
	if req.Method == http.MethodGet {
		cr.Body.Write([]byte(req.URL.Path))
		cr.Status = http.StatusOK
		cr.Hdr = w.Header()
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) RTPage(w http.ResponseWriter, req *http.Request, app *tw.TwApp) (ok bool) {
	if req.Method == http.MethodGet {
		cr.Body.Write([]byte(req.URL.Path))
		cr.Status = http.StatusOK
		cr.Hdr = w.Header()
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) HomePage(w http.ResponseWriter, req *http.Request, app *tw.TwApp) (ok bool) {
	if req.Method == http.MethodGet {
		cr.Hdr = w.Header()
		tmplFile := "./tmpl/home.tmpl"
		err := app.Auth()
		if err == nil {
			cr.Status = 302
			cr.Hdr.Set("Location", "/rt")
			cr.Body.Write([]byte{})
			return cr.WriteHTTPResponse(w)
		}

		dat, err2 := ioutil.ReadFile(tmplFile)
		if err2 != nil {
			log.Printf("Method: Get, Error reading template file. Error: %v", err2)
			return cr.SendInternalError(w)
		}

		err3 := cr.ParseTemplate(dat, nil)
		if err3 != nil {
			log.Printf("Errors: %v, %v", err3, err2)
			cr.SendInternalError(w)
		}

		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) SendNotFound(w http.ResponseWriter, url string) (ok bool) {
	cr.Status = http.StatusNotFound
	cr.Body.Write([]byte("Not Found: " + url))
	return cr.WriteHTTPResponse(w)
}

func (cr *ContentResponse) SendInternalError(w http.ResponseWriter) (ok bool) {
	cr.Status = http.StatusInternalServerError
	cr.Body.Write([]byte("Internal Server Error"))
	return cr.WriteHTTPResponse(w)
}

func (cr *ContentResponse) WriteHTTPResponse(w http.ResponseWriter) (ok bool) {
	w.WriteHeader(cr.Status)
	cr.Body.WriteTo(w)
	return true
}

func (cr *ContentResponse) ParseTemplate(b []byte, data interface{}) (err error) {

	tmpl, err := template.New("HTML").Parse(string(b))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Method: Get, Error Parsing template file"))
	}

	cr.Status = http.StatusOK
	cr.Hdr.Set("Content-Type", "text/html")
	err = tmpl.Execute(&cr.Body, data)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Method: Get, Error Executing template file"))
	}

	return nil
}
