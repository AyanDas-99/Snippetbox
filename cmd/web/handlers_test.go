package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AyanDas-99/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	rs, er := ts.Client().Get(ts.URL + "/ping")
	if er != nil {
		t.Fatal(er)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	defer ts.Close()
	body, er := io.ReadAll(rs.Body)

	if er != nil {
		t.Fatal(er)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
