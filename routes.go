package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func statusHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", Version)
	}
}

func router(bs BoltStore) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler())
	r.HandleFunc("/api/ident/", apiIdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/", apiIdentPUTHandler(bs)).Methods("PUT") // Funnel bad request for proper response.
	r.HandleFunc("/api/ident/", apiIdentDELETEHandler(bs)).Methods("DELETE")
	r.HandleFunc("/api/ident/{id}", apiIdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/{id}", apiIdentPUTHandler(bs)).Methods("PUT")
	r.HandleFunc("/api/ident/{id}", apiIdentDELETEHandler(bs)).Methods("DELETE")
	return r
}

// NewMux returns a new http.ServeMux with established routes.
func NewMux(bs BoltStore) *http.ServeMux {
	r := router(bs)

	s := http.NewServeMux()
	s.Handle("/", r)
	return s
}
