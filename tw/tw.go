package tw

import (
	"io/ioutil"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	yaml "gopkg.in/yaml.v2"
)

var (
	store = sessions.NewCookieStore([]byte("something-very-secret"))
)

const (
	ConfigFile string = "./conf.yaml.example"
)

type TwApp struct {
	ConfigFile        string
	ConsumerKey       string `yaml:"consumerKey"`
	ConsumerSecret    string `yaml:"consumerSecret"`
	AccessToken       string `yaml:"accessToken"`
	AccessTokenSecret string `yaml:"accessTokenSecret"`
	RedirectURL       string `yaml:"redirectURL"`
}

func (t *TwApp) LoadConfig() error {
	data, err := ioutil.ReadFile(t.ConfigFile)
	confErrMsg := "Fatal error: Could not read app config"
	if err != nil {
		return errors.Wrap(err, confErrMsg)
	}

	if err := yaml.Unmarshal(data, t); err != nil {
		return errors.Wrap(err, confErrMsg)
	}
	return nil
}

func (t *TwApp) Auth() error {
	config := &clientcredentials.Config{
		ClientID:       t.ConsumerKey,
		ClientSecret:   t.ConsumerSecret,
		TokenURL:       "https://api.twitter.com/oauth2/token",
		Scopes:         nil,
		EndpointParams: nil,
	}

	httpClient := config.Client(oauth2.NoContext)
	twitter.NewClient(httpClient)

	return errors.New("twitter  error")
	// return nil
}

func (t *TwApp) IsLoggedIn(req *http.Request) (bool, error) {
	session, err := store.Get(req, "tw-goodstuff")
	if err != nil {
		return false, errors.Wrap(err, "Error retrieving session")
	}

	if _, ok := session.Values["IsLoggedIn"]; ok {
		return true, nil
	}

	return false, nil
}

func (t *TwApp) Logout(w http.ResponseWriter, req *http.Request) (bool, error) {

	session, err := store.Get(req, "tw-goodstuff")
	if err != nil {
		return false, errors.Wrap(err, "Error retrieving/creating session")
	}

	session.Options.MaxAge = -1
	session.Save(req, w)
	return true, nil
}
