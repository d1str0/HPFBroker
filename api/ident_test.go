package api

import (
	hpf "github.com/d1str0/HPFBroker"
	auth "github.com/d1str0/HPFBroker/auth"

	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/d1str0/hpfeeds"
	"github.com/gorilla/mux"
)

func test(t *testing.T, name string, router *mux.Router, method string, uri string, r io.Reader, token string, expStatus int, expResp string) {
	t.Run(name, func(t *testing.T) {
		req, err := http.NewRequest(method, uri, r)
		if err != nil {
			t.Fatal(err)
		}

		auth := fmt.Sprintf("Bearer %s", token)
		req.Header.Set("Authorization", auth)

		testRequest(t, router, req, expStatus, expResp)
	})
}

func TestIdentHandler(t *testing.T) {
	var secret = &auth.JWTSecret{}
	secret.SetSecret([]byte{0x0000000000000000000000000000000000000000000000000000000000000000})
	db, err := hpf.OpenDB(TestDBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := auth.InitRBAC()

	sc := &ServerContext{Version: "69.420.80085", JWTSecret: secret, DB: db, RBAC: r}

	router := router(sc)

	id := &hpfeeds.Identity{Ident: "test-ident", Secret: "test-secret", SubChannels: []string{"asdf"}, PubChannels: []string{}}
	id2 := &hpfeeds.Identity{Ident: "test-ident1", Secret: "test-secret", SubChannels: []string{"asdf"}, PubChannels: []string{}}

	t.Run("GET", func(t *testing.T) {
		db.SaveIdentity(id)
		db.SaveIdentity(id2)

		token := "totallynotvalid"

		// FAIL
		test(t, "User Not Found", router, "GET", "/api/ident/asdf", nil, token, http.StatusNotFound, ErrNotFound.Error())

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

			testRequestObj(t, router, req, http.StatusOK, []*hpfeeds.Identity{id, id2})
		})
		db.DeleteIdentity(id.Ident)
		db.DeleteIdentity(id2.Ident)
	})

	t.Run("PUT", func(t *testing.T) {
		// FAIL
		t.Run("Missing Identifier", func(t *testing.T) {

			r := encodeBody(t, id)
			req, err := http.NewRequest("PUT", "/api/ident/", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMissingID.Error())
		})

		// FAIL
		t.Run("Missing Request Body", func(t *testing.T) {

			req, err := http.NewRequest("PUT", "/api/ident/test-ident", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrBodyRequired.Error())
		})

		// FAIL
		t.Run("Mismatched Identifier", func(t *testing.T) {

			r := encodeBody(t, id)
			req, err := http.NewRequest("PUT", "/api/ident/asdf", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMismatchedID.Error())
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
		defer db.DeleteIdentity("test-ident")

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

			testRequest(t, router, req, http.StatusBadRequest, ErrMismatchedID.Error())
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

		db.SaveIdentity(id)

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

			testRequest(t, router, req, http.StatusNotFound, ErrNotFound.Error())
		})

	})
}
