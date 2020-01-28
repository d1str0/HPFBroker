package api

import (
	hpf "github.com/d1str0/HPFBroker"
	auth "github.com/d1str0/HPFBroker/auth"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

const TestDBPath = ".test.db"

func testRouter(t *testing.T, db *hpf.DB) *mux.Router {
	var secret = &auth.JWTSecret{}
	secret.SetSecret([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	r := auth.InitRBAC()

	sc := &hpf.ServerContext{Version: "69.420.80085", JWTSecret: secret, DB: db, RBAC: r}

	router := router(sc)
	return router
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
