package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kawabatas/gomock/handler"
)

func TestDebugHandler(t *testing.T) {
	client := &http.Client{}

	t.Run("Not Use Mock", func(t *testing.T) {
		ts := httptest.NewServer(NewRouter())
		defer ts.Close()

		req, _ := http.NewRequest("GET", ts.URL+"/hello", nil)
		res, _ := client.Do(req)
		defer res.Body.Close()
		resBody, _ := ioutil.ReadAll(res.Body)
		resMessage := string(resBody)
		if resMessage != helloMessage {
			t.Errorf("%v want %v", resMessage, helloMessage)
		}
	})

	t.Run("Use Mock", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(handler.HandleStub))
		defer ts.Close()

		req, _ := http.NewRequest("GET", ts.URL+"/hello", nil)
		res, _ := client.Do(req)
		defer res.Body.Close()
		resBody, _ := ioutil.ReadAll(res.Body)
		resMessage := string(resBody)
		if resMessage == helloMessage {
			t.Errorf("%v not want %v", resMessage, helloMessage)
		}
	})
}
