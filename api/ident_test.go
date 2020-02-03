package api

import (
	hpf "github.com/d1str0/HPFBroker"
	auth "github.com/d1str0/HPFBroker/auth"

	"net/http"
	"testing"

	"github.com/d1str0/hpfeeds"
)

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

	t.Run("GET", func(t *testing.T) {
		db.SaveIdentity(id)
		db.SaveIdentity(id2)

		// Sanity Check FAIL
		testNoAuth(t, "User Found (No Auth)", router, "GET", "/api/ident/test-ident", nil, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		test(t, "User Found Invalid Token", router, "GET", "/api/ident/test-ident", nil, invalidToken, http.StatusUnauthorized, "token contains an invalid number of segments")

		// FAIL
		test(t, "User Not Found (HPF Reader)", router, "GET", "/api/ident/asdf", nil, hpfReaderToken, http.StatusNotFound, http.StatusText(http.StatusNotFound))

		// SUCCESS
		testObj(t, "User Found (HPF Reader)", router, "GET", "/api/ident/test-ident", nil, hpfReaderToken, http.StatusOK, id)

		// SUCCESS
		testObj(t, "Get All (HPF Reader)", router, "GET", "/api/ident/", nil, hpfReaderToken, http.StatusOK, []*hpfeeds.Identity{id, id2})

		db.DeleteIdentity(id.Ident)
		db.DeleteIdentity(id2.Ident)
	})

	t.Run("PUT", func(t *testing.T) {
		// FAIL
		r := encodeBody(t, id)

		testNoAuth(t, "Create Person (No Auth)", router, "PUT", "/api/ident/test-ident", r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

		test(t, "Missing Identifier (HPF Admin)", router, "PUT", "/api/ident/", r, hpfAdminToken, http.StatusBadRequest, ErrMissingID.Error())

		// FAIL
		test(t, "Missing Request Body (HPF Admin)", router, "PUT", "/api/ident/test-ident", nil, hpfAdminToken, http.StatusBadRequest, ErrBodyRequired.Error())

		// FAIL
		r = encodeBody(t, id)
		test(t, "Mismatched Identifier (HPF Admin)", router, "PUT", "/api/ident/asdf", r, hpfAdminToken, http.StatusBadRequest, ErrMismatchedID.Error())

		// FAIL
		r = encodeBody(t, id)
		test(t, "Create Ident (HPF Reader)", router, "PUT", "/api/ident/test-ident", r, hpfReaderToken, http.StatusForbidden, http.StatusText(http.StatusForbidden))

		// SUCCESS
		r = encodeBody(t, id)
		testObj(t, "Create Ident (HPF Admin)", router, "PUT", "/api/ident/test-ident", r, hpfAdminToken, http.StatusCreated, id)
		defer db.DeleteIdentity("test-ident")

		// SUCCESS
		id.Secret = "test-secret2"
		r = encodeBody(t, id)
		testObj(t, "Update Ident (HPF Admin)", router, "PUT", "/api/ident/test-ident", r, hpfAdminToken, http.StatusOK, id)

		// FAIL
		r = encodeBody(t, id2)
		test(t, "Update Mismatched Identifier (HPF Admin)", router, "PUT", "/api/ident/test-ident", r, hpfAdminToken, http.StatusBadRequest, ErrMismatchedID.Error())
	})

	t.Run("DELETE", func(t *testing.T) {
		// SUCCESS
		test(t, "Delete All (HPF Reader)", router, "DELETE", "/api/ident/", nil, hpfReaderToken, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		test(t, "Delete All (HPF Admin)", router, "DELETE", "/api/ident/", nil, hpfAdminToken, http.StatusNoContent, "")
		testNoAuth(t, "Delete All (No Auth)", router, "DELETE", "/api/ident/", nil, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		db.SaveIdentity(id)

		// SUCCESS
		test(t, "Delete One (HPF Admin)", router, "DELETE", "/api/ident/test-ident", nil, hpfAdminToken, http.StatusNoContent, "")
		testNoAuth(t, "Delete One (No Auth)", router, "DELETE", "/api/ident/test-ident", nil, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

		// FAIL
		test(t, "Delete One Not Found (HPF Admin)", router, "DELETE", "/api/ident/test-ident", nil, hpfAdminToken, http.StatusNotFound, http.StatusText(http.StatusNotFound))

	})
}
