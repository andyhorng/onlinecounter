package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleTouch(w, r)
}

func TestTouch(t *testing.T) {

	h := handler{}

	server := httptest.NewServer(h)
	defer server.Close()

	for i := 0; i < 1000; i++ {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fail()
		}

		if resp.StatusCode != 200 {
			t.Fail()
		}

		if p.count() != i+1 {
			t.Fail()
		}
	}
}
