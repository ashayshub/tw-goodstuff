package main_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ashayshub/tw-goodstuff"
)

const ep string = "http://" + main.HostAddr + ":" + main.HostPort

func TestGetMostFavEp(t *testing.T){
	epfv := ep + "/fav"
	if _, er := http.Get(epfv); er != nil {
		fmt.Println("Error: ", er)
		t.Fail()
	}
}

func TestGetMostRtEp(t *testing.T){
	eprt := ep + "/rt"
	if _, er := http.Get(eprt); er != nil {
		fmt.Println("Error: ", er)
		t.Fail()
	}
}

func TestGetSlashEp(t *testing.T) {
	if _, er := http.Get(ep); er != nil {
		fmt.Println("Error: ", er)
		t.Fail()
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
