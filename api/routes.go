package api

import (
	auth "github.com/d1str0/HPFBroker/auth"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: Take a server object so we can display a version number.
func statusHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "I'm online D:")
	}
}

func router(sc *ServerContext) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/status", statusHandler())

	r.HandleFunc("/api/authenticate", AuthHandler(sc)).Methods("POST")

	r.HandleFunc("/api/ident/",
		sc.Permission(
			auth.PermHPFRead,
			IdentGETHandler(sc),
		)).Methods("GET")

	r.HandleFunc("/api/ident/",
		sc.Permission(
			auth.PermHPFWrite,
			IdentPUTHandler(sc),
		)).Methods("PUT") // Funnel bad request for proper response.

	r.HandleFunc("/api/ident/",
		sc.Permission(
			auth.PermHPFWrite,
			IdentDELETEHandler(sc),
		)).Methods("DELETE")

	r.HandleFunc("/api/ident/{id}",
		sc.Permission(
			auth.PermHPFRead,
			IdentGETHandler(sc),
		)).Methods("GET")

	r.HandleFunc("/api/ident/{id}",
		sc.Permission(
			auth.PermHPFWrite,
			IdentPUTHandler(sc),
		)).Methods("PUT")

	r.HandleFunc("/api/ident/{id}",
		sc.Permission(
			auth.PermHPFWrite,
			IdentDELETEHandler(sc),
		)).Methods("DELETE")

	r.HandleFunc("/api/user/",
		sc.Permission(
			auth.PermUserRead,
			UserGETHandler(sc),
		)).Methods("GET")

	r.HandleFunc("/api/user/",
		sc.Permission(
			auth.PermUserWrite,
			UserPUTHandler(sc),
		)).Methods("PUT") // Funnel bad request for proper response.

	r.HandleFunc("/api/user/",
		sc.Permission(
			auth.PermUserWrite,
			UserDELETEHandler(sc),
		)).Methods("DELETE")

	r.HandleFunc("/api/user/{id}",
		sc.Permission(
			auth.PermUserRead,
			UserGETHandler(sc),
		)).Methods("GET")

	r.HandleFunc("/api/user/{id}",
		sc.Permission(
			auth.PermUserWrite,
			UserPUTHandler(sc),
		)).Methods("PUT")

	r.HandleFunc("/api/user/{id}",
		sc.Permission(
			auth.PermUserWrite,
			UserDELETEHandler(sc),
		)).Methods("DELETE")

	return r
}

// NewMux returns a new http.ServeMux with established routes.
func NewMux(sc *ServerContext) *http.ServeMux {
	r := router(sc)

	s := http.NewServeMux()
	s.Handle("/", r)
	return s
}
