package api

import (
	hpf "github.com/d1str0/HPFBroker"
	auth "github.com/d1str0/HPFBroker/auth"

	"net/http"
	"testing"
)

func TestUserHandler(t *testing.T) {
	var secret = &auth.JWTSecret{}
	secret.SetSecret([]byte{0x0000000000000000000000000000000000000000000000000000000000000000})
	db, err := hpf.OpenDB(TestDBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := auth.InitRBAC()

	sc := &hpf.ServerContext{Version: "69.420.80085", JWTSecret: secret, DB: db, RBAC: r}

	router := router(sc)

	u, err := hpf.NewUser("test-name", "test-password", auth.RoleUserReader)
	if err != nil {
		t.Fatal(err)
	}
	u2, err := hpf.NewUser("test2-name", "test-password", auth.RoleUserReader)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GET", func(t *testing.T) {
		db.SaveUser(u)
		db.SaveUser(u2)

		// FAIL
		t.Run("User Not Found", func(t *testing.T) {

			req, err := http.NewRequest("GET", "/api/user/asdf", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNotFound, ErrNotFound)
		})

		// SUCCESS
		t.Run("User Found", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/user/test-name", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Expected response structure
			ur := &UserResp{Name: u.Name, Role: u.Role}

			testRequestObj(t, router, req, http.StatusOK, ur)
		})

		// SUCCESS
		t.Run("Get All", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/user/", nil)
			if err != nil {
				t.Fatal(err)
			}

			ur := &UserResp{Name: u.Name, Role: u.Role}
			ur2 := &UserResp{Name: u2.Name, Role: u2.Role}

			testRequestObj(t, router, req, http.StatusOK, []*UserResp{ur, ur2})
		})
		db.DeleteUser(u.Name)
		db.DeleteUser(u2.Name)
	})

	t.Run("PUT", func(t *testing.T) {
		// FAIL
		t.Run("Missing Identifier", func(t *testing.T) {

			r := encodeBody(t, u)
			req, err := http.NewRequest("PUT", "/api/user/", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMissingIdentifier)
		})

		// FAIL
		t.Run("Missing Request Body", func(t *testing.T) {

			req, err := http.NewRequest("PUT", "/api/user/test-name", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrBodyRequired)
		})

		// FAIL
		t.Run("Mismatched Identifier", func(t *testing.T) {

			r := encodeBody(t, u)
			req, err := http.NewRequest("PUT", "/api/user/asdf", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMismatchedIdentifier)
		})

		// SUCCESS
		t.Run("Create User", func(t *testing.T) {

			r := encodeBody(t, u)
			req, err := http.NewRequest("PUT", "/api/user/test-name", r)
			if err != nil {
				t.Fatal(err)
			}

			ur := &UserResp{Name: u.Name, Role: u.Role}

			testRequestObj(t, router, req, http.StatusCreated, ur)
		})
		defer db.DeleteUser("test-name")

		// SUCCESS
		t.Run("Update User", func(t *testing.T) {

			r := encodeBody(t, u)
			req, err := http.NewRequest("PUT", "/api/user/test-name", r)
			if err != nil {
				t.Fatal(err)
			}

			ur := &UserResp{Name: u.Name, Role: u.Role}

			testRequestObj(t, router, req, http.StatusOK, ur)
		})

		// FAIL
		t.Run("Update Mismatched User", func(t *testing.T) {

			r := encodeBody(t, u2)
			req, err := http.NewRequest("PUT", "/api/user/test-name", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMismatchedIdentifier)
		})
	})

	t.Run("DELETE", func(t *testing.T) {
		// SUCCESS
		t.Run("Delete All", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/api/user/", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNoContent, "")
		})

		db.SaveUser(u)

		// SUCCESS
		t.Run("Delete One", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/api/user/test-name", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNoContent, "")
		})

		// SUCCESS
		t.Run("Delete One Not Found", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/api/user/test-name", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNotFound, ErrNotFound)
		})

	})
}
