package tw

import (
	"io/ioutil"
	"net/http"

	"github.com/dghubble/oauth1"
	twitter2 "github.com/dghubble/oauth1/twitter"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const (
	ConfigFile string = "./conf.yaml.example"
)

var (
	SessionPath string = "./tmp/sessions"
	store              = sessions.NewFilesystemStore(SessionPath, []byte("something-very-secret"))
)

type TwApp struct {
	ConfigFile        string
	ConsumerKey       string `yaml:"consumerKey"`
	ConsumerSecret    string `yaml:"consumerSecret"`
	RedirectURL       string `yaml:"redirectURL"`
	AccessToken       string
	AccessTokenSecret string
}

func (t *TwApp) LoadConfig() error {
	confErrMsg := "Fatal error: Could not read app config or some/all params empty"
	data, err := ioutil.ReadFile(t.ConfigFile)
	if err != nil {
		return errors.Wrap(err, confErrMsg)
	}

	if err := yaml.Unmarshal(data, t); err != nil {
		return errors.Wrap(err, confErrMsg)
	}

	if t.ConsumerKey == "" || t.ConsumerSecret == "" {
		return errors.New(confErrMsg)
	}

	return nil
}

func (t *TwApp) CreateConfig() oauth1.Config {
	return oauth1.Config{
		ConsumerKey:    t.ConsumerKey,
		ConsumerSecret: t.ConsumerSecret,
		CallbackURL:    t.RedirectURL,
		Endpoint:       twitter2.AuthorizeEndpoint,
	}
}

func (t *TwApp) FetchRequestToken() (string, error) {
	c := t.CreateConfig()
	requestToken, _, err := c.RequestToken()
	if err != nil {
		return "", errors.Wrap(err, "Error during c.RequestToken()")
	}

	authorizationURL, err := c.AuthorizationURL(requestToken)
	if err != nil {
		return "", errors.Wrap(err, "Error during c.AuthorizationURL()")
	}

	return authorizationURL.String(), nil
}

func (t *TwApp) Auth(w http.ResponseWriter, req *http.Request) (err error) {
	c := t.CreateConfig()
	v := req.URL.Query()

	oauth_token := v.Get("oauth_token")
	oauth_verifier := v.Get("oauth_verifier")

	if oauth_token == "" || oauth_verifier == "" {
		return errors.New("Empty: oauth_token or oauth_verifier")
	}

	t.AccessToken, t.AccessTokenSecret, err = c.AccessToken(oauth_token, t.ConsumerSecret, oauth_verifier)
	token := oauth1.NewToken(t.AccessToken, t.AccessTokenSecret)

	// Save token to session
	session, err := store.Get(req, "tw-goodstuff")
	session.Values["IsLoggedIn"] = true
	session.Values["Token"] = token.Token
	session.Values["TokenSecret"] = token.TokenSecret

	if err := session.Save(req, w); err != nil {
		return errors.Wrap(err, "Error saving the session")
	}

	// httpClient := config.Client(oauth1.NoContext, token)
	// client := twitter.NewClient(httpClient)
	return nil
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
	if err := session.Save(req, w); err != nil {
		return false, errors.Wrap(err, "Error saving the session")
	}

	return true, nil
}
