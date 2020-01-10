package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/d1str0/hpfeeds"
	"github.com/gorilla/mux"
)

func TestIdentHandler(t *testing.T) {
	bs := getTestDB(t)
	router := router(bs)

	id := hpfeeds.Identity{Ident: "test-ident", Secret: "test-secret", SubChannels: []string{"asdf"}, PubChannels: []string{}}
	id2 := hpfeeds.Identity{Ident: "test-ident1", Secret: "test-secret", SubChannels: []string{"asdf"}, PubChannels: []string{}}

	t.Run("GET", func(t *testing.T) {
		bs.SaveIdentity(id)
		bs.SaveIdentity(id2)

		// FAIL
		t.Run("User Not Found", func(t *testing.T) {

			req, err := http.NewRequest("GET", "/api/ident/asdf", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNotFound, ErrIdentNotFound)
		})

		// SUCCESS
		t.Run("User Found", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/ident/test-ident", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequestObj(t, router, req, http.StatusOK, id)
		})

		// SUCCESS
		t.Run("Get All", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/ident/", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequestObj(t, router, req, http.StatusOK, []hpfeeds.Identity{id, id2})
		})
		bs.DeleteIdentity(id.Ident)
		bs.DeleteIdentity(id2.Ident)
	})

	t.Run("PUT", func(t *testing.T) {
		// FAIL
		t.Run("Missing Identifier", func(t *testing.T) {

			r := encodeBody(t, id)
			req, err := http.NewRequest("PUT", "/api/ident/", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMissingIdentifier)
		})

		// FAIL
		t.Run("Missing Request Body", func(t *testing.T) {

			req, err := http.NewRequest("PUT", "/api/ident/test-ident", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrBodyRequired)
		})

		// FAIL
		t.Run("Mismatched Identifier", func(t *testing.T) {

			r := encodeBody(t, id)
			req, err := http.NewRequest("PUT", "/api/ident/asdf", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMismatchedIdentifier)
		})

		// SUCCESS
		t.Run("Create Ident", func(t *testing.T) {

			r := encodeBody(t, id)
			req, err := http.NewRequest("PUT", "/api/ident/test-ident", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequestObj(t, router, req, http.StatusCreated, id)
		})
		defer bs.DeleteIdentity("test-ident")

		// SUCCESS
		t.Run("Update Ident", func(t *testing.T) {

			r := encodeBody(t, id)
			req, err := http.NewRequest("PUT", "/api/ident/test-ident", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequestObj(t, router, req, http.StatusOK, id)
		})

		// FAIL
		t.Run("Update Mismatched Ident", func(t *testing.T) {

			r := encodeBody(t, id2)
			req, err := http.NewRequest("PUT", "/api/ident/test-ident", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMismatchedIdentifier)
		})
	})

	t.Run("DELETE", func(t *testing.T) {
		// SUCCESS
		t.Run("Delete All", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/api/ident/", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNoContent, "")
		})

		bs.SaveIdentity(id)

		// SUCCESS
		t.Run("Delete One", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/api/ident/test-ident", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNoContent, "")
		})

		// SUCCESS
		t.Run("Delete One Not Found", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/api/ident/test-ident", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNotFound, ErrIdentNotFound)
		})

	})
}

// encodeBody is used to encode a request body
func encodeBody(t *testing.T, obj interface{}) io.Reader {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		t.Fatalf("error encoding obj: %#v", err)
	}
	return buf
}

func testRequest(t *testing.T, router *mux.Router, req *http.Request, expectedStatus int, expected string) {
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code:\n\tgot %v \n\twant %v",
			status, expectedStatus)
	}

	respBody := strings.TrimSuffix(rr.Body.String(), "\n")
	if respBody != expected {
		t.Errorf("handler returned unexpected body:\n\tgot %s \n\twant %s",
			respBody, expected)
	}
}

func testRequestObj(t *testing.T, router *mux.Router, req *http.Request, expectedStatus int, obj interface{}) {
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code:\n\tgot %v \n\twant %v",
			status, expectedStatus)
	}

	s, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("Error marshaling: %#v", err)
	}
	expected := string(s)

	respBody := strings.TrimSuffix(rr.Body.String(), "\n")
	if respBody != expected {
		t.Errorf("handler returned unexpected body:\n\tgot %s \n\twant %s",
			respBody, expected)
	}
}
