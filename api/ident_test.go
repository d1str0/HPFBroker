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

func testObj(t *testing.T, name string, router *mux.Router, method string, uri string, r io.Reader, token string, expStatus int, expObj interface{}) {
	t.Run(name, func(t *testing.T) {
		req, err := http.NewRequest(method, uri, r)
		if err != nil {
			t.Fatal(err)
		}

		auth := fmt.Sprintf("Bearer %s", token)
		req.Header.Set("Authorization", auth)

		testRequestObj(t, router, req, expStatus, expObj)
	})
}

func testNoAuth(t *testing.T, name string, router *mux.Router, method string, uri string, r io.Reader, expStatus int, expResp string) {
	t.Run(name, func(t *testing.T) {
		req, err := http.NewRequest(method, uri, r)
		if err != nil {
			t.Fatal(err)
		}

		testRequest(t, router, req, expStatus, expResp)
	})
}

func TestIdentHandler(t *testing.T) {
	var secret = &auth.JWTSecret{}
	secret.SetSecret([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
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

	invalidToken := "totallynotvalid"

	hpfReaderToken, err := sc.JWTSecret.Sign(auth.RoleHPFReader)
	if err != nil {
		t.Fatal(err)
	}

	hpfAdminToken, err := sc.JWTSecret.Sign(auth.RoleHPFAdmin)
	if err != nil {
		t.Fatal(err)
	}

	superAdminToken, err := sc.JWTSecret.Sign(auth.RoleSuperAdmin)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GET", func(t *testing.T) {
		db.SaveIdentity(id)
		db.SaveIdentity(id2)

		// Sanity Check FAIL
		test(t, "User Found Invalid Token", router, "GET", "/api/ident/test-ident", nil, invalidToken, http.StatusUnauthorized, "token contains an invalid number of segments")

		// FAIL
		test(t, "User Not Found (HPF Reader)", router, "GET", "/api/ident/asdf", nil, hpfReaderToken, http.StatusNotFound, ErrNotFound.Error())
		test(t, "User Not Found (HPF Admin)", router, "GET", "/api/ident/asdf", nil, hpfAdminToken, http.StatusNotFound, ErrNotFound.Error())
		test(t, "User Not Found (Super Admin)", router, "GET", "/api/ident/asdf", nil, superAdminToken, http.StatusNotFound, ErrNotFound.Error())
		testNoAuth(t, "User Not Found (No Auth)", router, "GET", "/api/ident/asdf", nil, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

		// SUCCESS
		testObj(t, "User Found (HPF Reader)", router, "GET", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusOK, id)
		testObj(t, "User Found (HPF Admin)", router, "GET", "/api/ident/test-ident", nil, hpfAdminToken, http.StatusOK, id)
		testObj(t, "User Found (Super Admin)", router, "GET", "/api/ident/test-ident", nil, superAdminToken, http.StatusOK, id)
		testNoAuth(t, "User Found (No Auth)", router, "GET", "/api/ident/test-ident", nil, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

		// SUCCESS
		testObj(t, "Get All (HPF Reader)", router, "GET", "/api/ident/", nil, hpfReaderToken, http.StatusOK, []*hpfeeds.Identity{id, id2})
		testObj(t, "Get All (HPF Admin)", router, "GET", "/api/ident/", nil, hpfAdminToken, http.StatusOK, []*hpfeeds.Identity{id, id2})
		testObj(t, "Get All (Super Admin)", router, "GET", "/api/ident/", nil, superAdminToken, http.StatusOK, []*hpfeeds.Identity{id, id2})
		testNoAuth(t, "Get All (No Auth)", router, "GET", "/api/ident/", nil, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

		db.DeleteIdentity(id.Ident)
		db.DeleteIdentity(id2.Ident)
	})

	t.Run("PUT", func(t *testing.T) {
		// FAIL
		r := encodeBody(t, id)
		test(t, "Missing Identifier (HPF Reader)", router, "PUT", "/api/ident/", r, hpfReaderToken, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		test(t, "Missing Identifier (HPF Admin)", router, "PUT", "/api/ident/", r, hpfAdminToken, http.StatusBadRequest, ErrMissingID.Error())
		test(t, "Missing Identifier (Super Admin)", router, "PUT", "/api/ident/", r, superAdminToken, http.StatusBadRequest, ErrMissingID.Error())
		testNoAuth(t, "Missing Identifier (No Auth)", router, "PUT", "/api/ident/", nil, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())

		// FAIL
		test(t, "Missing Identifier (HPF Reader)", router, "PUT", "/api/ident/", nil, hpfReaderToken, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		test(t, "Missing Request Body (HPF Admin)", router, "PUT", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusBadRequest, ErrBodyRequired.Error())
		test(t, "Missing Request Body (Super Admin)", router, "PUT", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusBadRequest, ErrBodyRequired.Error())
		testNoAuth(t, "Missing Request Body (No Auth)", router, "PUT", "/api/ident/test-ident", nil, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())

		// FAIL
		test(t, "Mismatched Identifier (HPF Reader)", router, "PUT", "/api/ident/asdf", r, hpfReaderToken, http.StatusBadRequest, ErrMismatchedID.Error())
		test(t, "Mismatched Identifier (HPF Admin)", router, "PUT", "/api/ident/asdf", r, hpfReaderToken, http.StatusBadRequest, ErrMismatchedID.Error())
		test(t, "Mismatched Identifier (Super Admin)", router, "PUT", "/api/ident/asdf", r, hpfReaderToken, http.StatusBadRequest, ErrMismatchedID.Error())
		testNoAuth(t, "Mismatched Identifier (No Auth)", router, "PUT", "/api/ident/asdf", r, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())

		// SUCCESS
		testObj(t, "Create Ident (HPF Reader)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusCreated, id)
		testObj(t, "Create Ident (HPF Admin)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusCreated, id)
		testObj(t, "Create Ident (Super Admin)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusCreated, id)
		testNoAuth(t, "Create Ident (No Auth)", router, "PUT", "/api/ident/test-ident", r, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())
		defer db.DeleteIdentity("test-ident")

		// SUCCESS
		id.Secret = "test-secret2"
		r = encodeBody(t, id)
		testObj(t, "Update Ident (HPF Reader)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusOK, id)
		testObj(t, "Update Ident (HPF Admin)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusOK, id)
		testObj(t, "Update Ident (Super Admin)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusOK, id)
		testNoAuth(t, "Update Ident (No Auth)", router, "PUT", "/api/ident/test-ident", r, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())

		// FAIL
		r = encodeBody(t, id2)
		test(t, "Update Mismatched Identifier (HPF Reader)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusBadRequest, ErrMismatchedID.Error())
		test(t, "Update Mismatched Identifier (HPF Admin)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusBadRequest, ErrMismatchedID.Error())
		test(t, "Update Mismatched Identifier (Super Admin)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusBadRequest, ErrMismatchedID.Error())
		testNoAuth(t, "Update Mismatched Identifier (No Auth)", router, "PUT", "/api/ident/test-ident", r, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())
	})

	t.Run("DELETE", func(t *testing.T) {
		// SUCCESS
		test(t, "Delete All (HPF Reader)", router, "DELETE", "/api/ident/", nil, hpfReaderToken, http.StatusNoContent, "")
		test(t, "Delete All (HPF Admin)", router, "DELETE", "/api/ident/", nil, hpfReaderToken, http.StatusNoContent, "")
		test(t, "Delete All (Super Admin)", router, "DELETE", "/api/ident/", nil, hpfReaderToken, http.StatusNoContent, "")
		testNoAuth(t, "Delete All (No Auth)", router, "DELETE", "/api/ident/", nil, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())
		db.SaveIdentity(id)

		// SUCCESS
		test(t, "Delete One (HPF Reader)", router, "DELETE", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusNoContent, "")
		test(t, "Delete One (HPF Admin)", router, "DELETE", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusNoContent, "")
		test(t, "Delete One (Super Admin)", router, "DELETE", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusNoContent, "")
		testNoAuth(t, "Delete One (No Auth)", router, "DELETE", "/api/ident/test-ident", nil, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())

		// FAIL
		test(t, "Delete One Not Found (HPF Reader)", router, "DELETE", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusNotFound, ErrNotFound.Error())
		test(t, "Delete One Not Found (HPF Admin)", router, "DELETE", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusNotFound, ErrNotFound.Error())
		test(t, "Delete One Not Found (Super Admin)", router, "DELETE", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusNotFound, ErrNotFound.Error())
		testNoAuth(t, "Delete One Not Found (No Auth)", router, "DELETE", "/api/ident/test-ident", nil, http.StatusUnauthorized, ErrAuthInvalidCreds.Error())

	})
}
