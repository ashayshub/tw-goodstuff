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
		tmplFile := "./tmpl/home.tmpl"
		cr.Hdr = w.Header()
		err := app.Auth()
		if err == nil {
			cr.Status = 302
			cr.Hdr.Set("Location", "/rt")
			cr.Body.Write([]byte{})
			return cr.WriteHTTPResponse(w)
		}

		err2 := cr.ParseTemplate(tmplFile)
		if err2 != nil {
			log.Printf("Errors: %v, %v", err2, err)
			cr.SendInternalError(w)
		}
		cr.Hdr.Set("Content-Type", "text/html")
		cr.Status = http.StatusOK
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) SendNotFound(w http.ResponseWriter, url string) (ok bool) {
	cr.Status = http.StatusNotFound
	cr.Body.Write([]byte("Not Found: " + url))
	cr.Hdr = w.Header()
	cr.WriteHTTPResponse(w)
	return false
}

func (cr *ContentResponse) SendInternalError(w http.ResponseWriter) (ok bool) {
	cr.Status = http.StatusInternalServerError
	cr.Body.Write([]byte("Internal Server Error"))
	cr.Hdr = w.Header()
	cr.WriteHTTPResponse(w)
	return true
}

func (cr *ContentResponse) WriteHTTPResponse(w http.ResponseWriter) (ok bool) {
	w.WriteHeader(cr.Status)
	if len(cr.Hdr) != 0 {
		if err := cr.Hdr.Write(w); err != nil {
			log.Printf("Could not write Header: %#v", cr.Hdr)
		}
	}
	if _, err := cr.Body.WriteTo(w); err != nil {
		log.Printf("Could not write Body: %#v", cr.Body)
	}
	return true
}

func (cr *ContentResponse) ParseTemplate(tmplFile string) (err error) {
	dat, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Method: Get, Error reading template file"))
	}

	tmpl, err := template.New("HTML").Parse(string(dat))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Method: Get, Error Parsing template file"))
	}

	err = tmpl.Execute(&cr.Body, nil)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Method: Get, Error Executing template file"))
	}

	return nil
}
