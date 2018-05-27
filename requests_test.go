package requests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlainTextResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	cli := NewClient(ts.URL).Accept("text/plain")

	var text string
	err := cli.NewRequest("GET", "/").
		Into(&text).
		Run()
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}

	if text != "Hello, client" {
		t.Errorf("Failed reading plain text body: got %+q, expected %+q", text, "Hello, client")
	}
}

func TestJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "encoding/json")
		fmt.Fprintln(w, `{ "message": "Hello, client" }`)
	}))
	defer ts.Close()

	cli := NewClient(ts.URL)

	var response struct {
		Message string `json:"message"`
	}

	err := cli.NewRequest("GET", "/").
		Into(&response).
		Run()
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}

	if response.Message != "Hello, client" {
		t.Errorf("Failed reading JSON body: got %#v, expected %+q", response, "Hello, client")
	}
}

func TestResponseHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Custom-Header", "bla")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	cli := NewClient(ts.URL)

	var customHeader string
	err := cli.NewRequest("GET", "/").
		HeaderInto("Custom-Header", &customHeader).
		ExpectedStatus(http.StatusNoContent).
		Run()
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}

	if customHeader != "bla" {
		t.Errorf("Failed reading custom header: got %+q, expected %+q", customHeader, "bla")
	}
}
