package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	user := new(userDetails)
	user.parseArgs()
	if ok := user.validateArgs(); !ok {
		log.Fatal("Validation failed")
	}
	fmt.Println("Your username is", strings.ToLower(user.name))
}

func usage() {
	usageText := `Usage:
%v: --username <username> -p[enter]
<password>
`
	fmt.Printf(usageText, os.Args[0])
}

type userDetails struct {
	name     string
	password string
	token    string
	login    bool
}

func (u *userDetails) oAuth() (resp string, ok bool) {
	return "", true
}

func (u *userDetails) parseArgs() {
	flag.StringVar(&u.name, "username", "", "Enter your twitter username")
	flag.Parse()
}

func (u *userDetails) validateArgs() (ok bool) {
	if u.name == "" {
		usage()
		log.Print("--username flag cannot not be empty")
		return false
	}
	return true
}
