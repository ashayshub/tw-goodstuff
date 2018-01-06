package tw

import (
	"io/ioutil"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	yaml "gopkg.in/yaml.v2"
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
	client := twitter.NewClient(httpClient)
	params := &twitter.UserShowParams{ScreenName: "golang"}
	if _, _, err := client.Users.Show(params); err != nil {
		return errors.Wrap(err, "Fatal error: Could not return User Show")
	}

	client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{})
	return nil
}
