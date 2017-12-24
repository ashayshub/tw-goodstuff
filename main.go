package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
)

func main() {

	user := new(UserProperty)
	usage := user.ParseArgs()
	if ok := user.ValidateArgs(&usage); !ok {
		log.Fatal("Validation failed")
	}
	user.OAuth()
	fmt.Println("Your username is", strings.ToLower(user.Name))
}

type UserProperty struct {
	Name        string
	AccessToken string
}

func (u *UserProperty) OAuth() (resp string, ok bool) {
	config := &oauth2.Config{}
	token := &oauth2.Token{AccessToken: u.AccessToken}
	httpClient := config.Client(oauth2.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	userShowParams := &twitter.UserShowParams{ScreenName: u.Name}
	user, err, code := client.Users.Show(userShowParams)
	if code != nil {
		fmt.Println("Error: ", err)
		return "", false
	}
	fmt.Printf("USERS SHOW:\n%+v, %#v, %#v\n", user, err, code)
	return "", true
}

func (u *UserProperty) ParseArgs() (f func()) {
	flags := flag.NewFlagSet("tw-goodstuff", flag.ExitOnError)
	flags.StringVar(&u.Name, "username", "", "Enter your twitter username")
	flags.StringVar(&u.AccessToken, "app-access-token", "", "Twitter Application Access Token")
	flags.Parse(os.Args[1:])
	return flags.Usage
}

func (u *UserProperty) ValidateArgs(usage *func()) (ok bool) {
	if u.Name == "" {
		(*usage)()
		fmt.Println("Error: --username flag cannot not be empty")
		return false
	}
	if u.AccessToken == "" {
		fmt.Println("Error: --app-access-token flag cannot not be empty")
		return false
	}
	return true
}
