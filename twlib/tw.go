package tw

import (
	"net/http"
)

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
