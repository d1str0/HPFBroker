package api

import (
	hpf "github.com/d1str0/HPFBroker"
	auth "github.com/d1str0/HPFBroker/auth"

	"net/http"
	"testing"
)

// TODO: Rename
func newUser(t *testing.T, ur *UserReq) *hpf.User {
	u, err := hpf.NewUser(ur.Name, ur.Password, ur.Role)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

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

	u := &UserReq{Name: "test-name", Password: "test-password", Role: auth.RoleUserReader}
	u2 := &UserReq{Name: "test2-name", Password: "test-password", Role: auth.RoleUserReader}

	t.Run("GET", func(t *testing.T) {
		db.SaveUser(newUser(t, u))
		db.SaveUser(newUser(t, u2))

		// FAIL
		t.Run("User Not Found", func(t *testing.T) {

			req, err := http.NewRequest("GET", "/api/user/asdf", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusNotFound, ErrNotFound.Error())
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

			testRequest(t, router, req, http.StatusBadRequest, ErrMissingID.Error())
		})

		// FAIL
		t.Run("Missing Request Body", func(t *testing.T) {

			req, err := http.NewRequest("PUT", "/api/user/test-name", nil)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrBodyRequired.Error())
		})

		// FAIL
		t.Run("Mismatched Identifier", func(t *testing.T) {

			r := encodeBody(t, u)
			req, err := http.NewRequest("PUT", "/api/user/asdf", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrMismatchedID.Error())
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

			testRequest(t, router, req, http.StatusBadRequest, ErrMismatchedID.Error())
		})

		// FAIL
		t.Run("Invalid Username", func(t *testing.T) {
			bad_ur := &UserReq{Name: "name:with:colon", Password: "validPassW0rd!", Role: auth.RoleUserReader}

			r := encodeBody(t, bad_ur)
			req, err := http.NewRequest("PUT", "/api/user/name:with:colon", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrInvalidUserName.Error())
		})

		// FAIL
		t.Run("Invalid Password", func(t *testing.T) {
			bad_ur := &UserReq{Name: "validname", Password: "nope", Role: auth.RoleUserReader}

			r := encodeBody(t, bad_ur)
			req, err := http.NewRequest("PUT", "/api/user/validname", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrInvalidUserPassword.Error())
		})

		// FAIL
		t.Run("Invalid Role", func(t *testing.T) {
			bad_ur := &UserReq{Name: "validname", Password: "surething", Role: "doesnt_exist"}

			r := encodeBody(t, bad_ur)
			req, err := http.NewRequest("PUT", "/api/user/validname", r)
			if err != nil {
				t.Fatal(err)
			}

			testRequest(t, router, req, http.StatusBadRequest, ErrInvalidUserRole.Error())
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

		db.SaveUser(newUser(t, u))

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

			testRequest(t, router, req, http.StatusNotFound, ErrNotFound.Error())
		})

	})
}
