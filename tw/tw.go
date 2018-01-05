package tw

import (
	"net/http"
)

type TwApp struct {
	ConsumerKey       string `yaml:"consumerKey"`
	ConsumerSecret    string `yaml:"consumerSecret"`
	AccessToken       string `yaml:"accessToken"`
	AccessTokenSecret string `yaml:"accessTokenSecret"`
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
