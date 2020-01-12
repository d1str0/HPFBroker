package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/d1str0/hpfeeds"
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
	r.HandleFunc("/api/ident/", api.IdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/", api.IdentPUTHandler(bs)).Methods("PUT") // Funnel bad request for proper response.
	r.HandleFunc("/api/ident/", api.IdentDELETEHandler(bs)).Methods("DELETE")
	r.HandleFunc("/api/ident/{id}", api.IdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/{id}", api.IdentPUTHandler(bs)).Methods("PUT")
	r.HandleFunc("/api/ident/{id}", api.IdentDELETEHandler(bs)).Methods("DELETE")
	return r
}

func NewMux(bs BoltStore) *http.ServeMux {
	r := router(bs)

	s := http.NewServeMux()
	s.Handle("/", r)
	return s
}
