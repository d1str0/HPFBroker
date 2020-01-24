package api

import (
	hpf "github.com/d1str0/HPFBroker"

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

func router(sc *hpf.ServerContext) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler())
	r.HandleFunc("/api/ident/", IdentGETHandler(sc)).Methods("GET")
	r.HandleFunc("/api/ident/", IdentPUTHandler(sc)).Methods("PUT") // Funnel bad request for proper response.
	r.HandleFunc("/api/ident/", IdentDELETEHandler(sc)).Methods("DELETE")
	r.HandleFunc("/api/ident/{id}", IdentGETHandler(sc)).Methods("GET")
	r.HandleFunc("/api/ident/{id}", IdentPUTHandler(sc)).Methods("PUT")
	r.HandleFunc("/api/ident/{id}", IdentDELETEHandler(sc)).Methods("DELETE")
	return r
}

// NewMux returns a new http.ServeMux with established routes.
func NewMux(sc *hpf.ServerContext) *http.ServeMux {
	r := router(sc)

	s := http.NewServeMux()
	s.Handle("/", r)
	return s
}
