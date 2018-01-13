package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ashayshub/tw-goodstuff/tw"
	"github.com/dghubble/go-twitter/twitter"
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
	app.ConfigFile = tw.ConfigFile

	cr := &ContentResponse{}
	cr.Hdr = w.Header()
	cr.TmplFile = "./tmpl/home.tmpl"

	//Startup errors
	if err := app.LoadConfig(); err != nil {
		panic(err)
	}

	switch req.URL.Path {
	case "/fav":
		if ok := cr.FavPage(w, req, app); !ok {
			log.Println("Error Sending Fav Page")
		}

	case "/rt":
		if ok := cr.RTPage(w, req, app); !ok {
			log.Println("Error Sending RT Page")
		}

	case "/":
		if ok := cr.HomePage(w, req, app); !ok {
			log.Println("Error Sending Home Page")
		}

	case "/login":
		if ok := cr.LoginPage(w, req, app); !ok {
			log.Println("Error Sending Login Page")
		}

	default:
		cr.SendNotFound(w, req.URL.Path)
	}
}

type ContentResponse struct {
	Status   int
	Body     bytes.Buffer
	Hdr      http.Header
	TmplFile string
	TwUser   string
	TwFav    []twitter.Tweet
	TwRT     []twitter.Tweet
}

func (cr *ContentResponse) FavPage(w http.ResponseWriter, req *http.Request, app *tw.TwApp) (ok bool) {
	if req.Method == http.MethodGet {
		cr.Body.Write([]byte(req.URL.Path))
		cr.Status = http.StatusOK
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) RTPage(w http.ResponseWriter, req *http.Request, app *tw.TwApp) (ok bool) {
	if req.Method == http.MethodGet {
		cr.Body.Write([]byte(req.URL.Path))
		cr.Status = http.StatusOK
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) HomePage(w http.ResponseWriter, req *http.Request, app *tw.TwApp) (ok bool) {
	if req.Method == http.MethodGet {

		ok, err := app.IsLoggedIn(req)
		if err != nil {
			// ok will remain false
			log.Printf("Error: %v", err)
		}

		if ok {
			cr.TwFav, cr.TwRT, err = app.GetFavRT(req)
			if err != nil {
				log.Println(err)
				return cr.SendRedirect(w, "/")
			}

			cr.TwUser, err = app.GetTwUser(req)
			if err != nil {
				log.Println(err)
				return cr.SendRedirect(w, "/")
			}

			if err := cr.ReadParseTmpl(); err != nil {
				log.Println(err)
				return cr.SendInternalError(w)
			}

			return cr.WriteHTTPResponse(w)
		}

		if err := cr.ReadParseTmpl(); err != nil {
			log.Println(err)
			return cr.SendInternalError(w)
		}

		// write cr.Body buffer to the "wire"
		return cr.WriteHTTPResponse(w)
	}
	return cr.SendNotFound(w, req.URL.Path)
}

func (cr *ContentResponse) LoginPage(w http.ResponseWriter, req *http.Request, app *tw.TwApp) (ok bool) {
	if req.Method == http.MethodPost {
		ok, err := app.IsLoggedIn(req)
		if err != nil {
			// ok will remain false
			log.Printf("Error: %v", err)
		}

		if ok {
			cr.Status = http.StatusOK
			cr.Body.Write([]byte("/"))
			return cr.WriteHTTPResponse(w)
		}

		authURL, err := app.FetchRequestToken()
		if err != nil {
			log.Println(err)
			return cr.SendInternalError(w)
		}

		cr.Status = http.StatusOK
		cr.Body.Write([]byte(authURL))
		return cr.WriteHTTPResponse(w)

	} else if req.Method == http.MethodGet {
		if err := app.Auth(w, req); err != nil {
			log.Println(err)
			return cr.SendInternalError(w)
		}
		return cr.SendRedirect(w, "/")
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
	cr.WriteHTTPResponse(w)
	// To handle tests
	return false
}

func (cr *ContentResponse) SendRedirect(w http.ResponseWriter, url string) (ok bool) {
	cr.Status = 302
	cr.Hdr.Set("Location", url)
	cr.Body.Write([]byte{})
	return cr.WriteHTTPResponse(w)
}

func (cr *ContentResponse) WriteHTTPResponse(w http.ResponseWriter) (ok bool) {
	w.WriteHeader(cr.Status)
	cr.Body.WriteTo(w)
	return true
}

func (cr *ContentResponse) ReadParseTmpl() error {
	b, err := ioutil.ReadFile(cr.TmplFile)
	if err != nil {
		return errors.Wrap(err, "Method: Get, Error reading template file")
	}

	// Parsed template is copied into cr.Body buffer
	err2 := cr.ParseTmpl(b)
	if err2 != nil {
		return err2
	}
	return nil
}

func (cr *ContentResponse) ParseTmpl(b []byte) (err error) {

	tmpl, err := template.New("HTML").Parse(string(b))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Method: Get, Error Parsing template file"))
	}

	cr.Status = http.StatusOK
	cr.Hdr.Set("Content-Type", "text/html")
	err = tmpl.Execute(&cr.Body, cr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Method: Get, Error Executing template file"))
	}

	return nil
}
