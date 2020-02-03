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

	u := &UserReq{Name: "test-name", Password: "test-password", Role: auth.RoleUserReader}
	u2 := &UserReq{Name: "test2-name", Password: "test-password", Role: auth.RoleUserReader}

	ur := &UserResp{Name: u.Name, Role: u.Role}
	ur2 := &UserResp{Name: u2.Name, Role: u2.Role}

	invalidToken := "totallynotvalid"

	userReaderToken, err := sc.JWTSecret.Sign(auth.RoleUserReader)
	if err != nil {
		t.Fatal(err)
	}

	userAdminToken, err := sc.JWTSecret.Sign(auth.RoleUserAdmin)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GET", func(t *testing.T) {
		db.SaveUser(newUser(t, u))
		db.SaveUser(newUser(t, u2))

		// FAIL
		testNoAuth(t, "User Not Found (No Auth)", router, "GET", "/api/user/asdf", nil, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		test(t, "User Found (Invalid Token)", router, "GET", "/api/user/asdf", nil, invalidToken, http.StatusUnauthorized, "token contains an invalid number of segments")

		test(t, "User Not Found (User Reader)", router, "GET", "/api/user/asdf", nil, userReaderToken, http.StatusNotFound, http.StatusText(http.StatusNotFound))

		testObj(t, "User Found (User Reader)", router, "GET", "/api/user/test-name", nil, userReaderToken, http.StatusOK, ur)
		testObj(t, "Get All Users (User Reader)", router, "GET", "/api/user/", nil, userReaderToken, http.StatusOK, []*UserResp{ur, ur2})

		db.DeleteUser(u.Name)
		db.DeleteUser(u2.Name)
	})

	t.Run("PUT", func(t *testing.T) {
		// FAIL
		r := encodeBody(t, u)
		testNoAuth(t, "Create Person (No Auth)", router, "PUT", "/api/user/test-name", r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

		r = encodeBody(t, u)
		test(t, "Missing Identifier (User Admin)", router, "PUT", "/api/user/", r, userAdminToken, http.StatusBadRequest, ErrMissingID.Error())

		// FAIL
		test(t, "Missing Request Body (User Admin)", router, "PUT", "/api/user/test-name", nil, userAdminToken, http.StatusBadRequest, ErrBodyRequired.Error())

		r = encodeBody(t, u)
		test(t, "Missing Identifier (User Admin)", router, "PUT", "/api/user/asdf", r, userAdminToken, http.StatusBadRequest, ErrMismatchedID.Error())

		// SUCCESS
		r = encodeBody(t, u)
		testObj(t, "Create User (User Admin)", router, "PUT", "/api/user/test-name", r, userAdminToken, http.StatusCreated, ur)
		defer db.DeleteUser("test-name")

		// SUCCESS
		r = encodeBody(t, u)
		testObj(t, "Update User (User Admin)", router, "PUT", "/api/user/test-name", r, userAdminToken, http.StatusOK, ur)

		// FAIL
		r = encodeBody(t, u2)
		test(t, "Update Mismatched User (User Admin)", router, "PUT", "/api/user/test-name", r, userAdminToken, http.StatusBadRequest, ErrMismatchedID.Error())

		// FAIL
		bad_ur := &UserReq{Name: "name:with:colon", Password: "validPassW0rd!", Role: auth.RoleUserReader}
		r = encodeBody(t, bad_ur)
		test(t, "Invalid Username (User Admin)", router, "PUT", "/api/user/name:with:colon", r, userAdminToken, http.StatusBadRequest, ErrInvalidUserName.Error())

		// FAIL
		bad_ur = &UserReq{Name: "validname", Password: "nope", Role: auth.RoleUserReader}
		r = encodeBody(t, bad_ur)
		test(t, "Invalid Password (User Admin)", router, "PUT", "/api/user/validname", r, userAdminToken, http.StatusBadRequest, ErrInvalidUserPassword.Error())

		// FAIL
		bad_ur = &UserReq{Name: "validname", Password: "surething", Role: "doesnt_exist"}
		r = encodeBody(t, bad_ur)
		test(t, "Invalid Password (User Admin)", router, "PUT", "/api/user/validname", r, userAdminToken, http.StatusBadRequest, ErrInvalidUserRole.Error())

	})

	t.Run("DELETE", func(t *testing.T) {
		// SUCCESS
		test(t, "Delete All (User Admin)", router, "DELETE", "/api/user/", nil, userAdminToken, http.StatusNoContent, "")

		db.SaveUser(newUser(t, u))

		// SUCCESS
		test(t, "Delete One (User Admin)", router, "DELETE", "/api/user/test-name", nil, userAdminToken, http.StatusNoContent, "")

		// FAIL
		test(t, "Delete One Not Found (User Admin)", router, "DELETE", "/api/user/test-name", nil, userAdminToken, http.StatusNotFound, http.StatusText(http.StatusNotFound))

	})
}
