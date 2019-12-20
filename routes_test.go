package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/d1str0/hpfeeds"
	bolt "go.etcd.io/bbolt"
)

const TestDBPath = "test.db"

func TestRoutes_statusHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statusHandler())

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := Version
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRoutes_apiIdentPUTHandler(t *testing.T) {
	bs := getTestDB(t)

	// PUT
	id := hpfeeds.Identity{Ident: "test-ident", Secret: "test-secret", SubChannels: []string{"asdf"}, PubChannels: []string{}}
	buf, err := json.Marshal(id)
	if err != nil {
		t.Fatal(err)
	}

	// FAIL
	t.Run("Missing Identifier", func(t *testing.T) {

		r := bytes.NewReader(buf)
		req, err := http.NewRequest("PUT", "/api/ident/", r)
		if err != nil {
			t.Fatal(err)
		}

		testRequest(t, bs, req, http.StatusBadRequest, ErrMissingIdentifier)
	})

	// FAIL
	t.Run("Mismatched Identifier", func(t *testing.T) {

		r := bytes.NewReader(buf)
		req, err := http.NewRequest("PUT", "/api/ident/asdf", r)
		if err != nil {
			t.Fatal(err)
		}

		testRequest(t, bs, req, http.StatusBadRequest, ErrMismatchedIdentifier)
	})
}

func testRequest(t *testing.T, bs BoltStore, req *http.Request, expectedStatus int, expected string) {

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiIdentPUTHandler(bs))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code:\n\tgot %v \n\twant %v",
			status, expectedStatus)
	}

	respBody := strings.TrimSuffix(rr.Body.String(), "\n")
	if respBody != expected {
		t.Errorf("handler returned unexpected body:\n\tgot %v \n\twant %v",
			respBody, expected)
	}
}

func getTestDB(t *testing.T) BoltStore {
	// Open up the boltDB file
	db, err := bolt.Open(TestDBPath, 0666, nil)
	if err != nil {
		t.Fatalf("Error opening db: %s", err.Error())
	}
	defer db.Close()

	// For use with HTTP handlers
	bs := BoltStore{db: db}

	// Prepare DB to ensure we have the appropriate buckets ready
	err = initializeDB(bs)
	if err != nil {
		t.Fatalf("Error initializing db: %s", err.Error())
	}
	return bs
}
