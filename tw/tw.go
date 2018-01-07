package tw

import (
	"io/ioutil"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	yaml "gopkg.in/yaml.v2"
)

const (
	ConfigFile string = "./conf.yaml.example"
)

type TwApp struct {
	ConsumerKey       string `yaml:"consumerKey"`
	ConsumerSecret    string `yaml:"consumerSecret"`
	AccessToken       string `yaml:"accessToken"`
	AccessTokenSecret string `yaml:"accessTokenSecret"`
	RedirectURL       string `yaml:"redirectURL"`
}

func (t *TwApp) LoadConfig() error {
	data, err := ioutil.ReadFile(ConfigFile)
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
