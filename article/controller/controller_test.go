package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {

	srv := httptest.NewServer(New(nil))
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/article/", srv.URL))
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
	_, err = strconv.Atoi(string(b))
	if err != nil {
		t.Fatalf("expected an integer; got %s,",err.Error())
	}

}