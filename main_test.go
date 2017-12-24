package main_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ashayshub/tw-goodstuff"
)

func TestHTTPEndPoint(t *testing.T) {
	ep := "http://" + main.HostAddr + ":" + main.HostPort
	
	if _, err := http.Get(ep); err != nil {
		fmt.Println("Error: ", err)
		t.Fail()
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
