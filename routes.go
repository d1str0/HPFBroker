package hpfbroker

import (
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

func router(db *DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler())
	r.HandleFunc("/api/ident/", apiIdentGETHandler(db)).Methods("GET")
	r.HandleFunc("/api/ident/", apiIdentPUTHandler(db)).Methods("PUT") // Funnel bad request for proper response.
	r.HandleFunc("/api/ident/", apiIdentDELETEHandler(db)).Methods("DELETE")
	r.HandleFunc("/api/ident/{id}", apiIdentGETHandler(db)).Methods("GET")
	r.HandleFunc("/api/ident/{id}", apiIdentPUTHandler(db)).Methods("PUT")
	r.HandleFunc("/api/ident/{id}", apiIdentDELETEHandler(db)).Methods("DELETE")
	return r
}

// NewMux returns a new http.ServeMux with established routes.
func NewMux(db *DB) *http.ServeMux {
	r := router(db)

	s := http.NewServeMux()
	s.Handle("/", r)
	return s
}
