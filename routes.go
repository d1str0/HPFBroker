package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/d1str0/hpfeeds"
	"github.com/gorilla/mux"
)

func statusHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", Version)
	}
}

//TODO return proper http codes per method
func apiIdentDELETEHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ident := vars["id"]

		// Delete user
		err := DeleteIdentity(bs, ident)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func apiIdentGETHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ident := vars["id"]

		// "/api/ident/"
		if ident == "" {
			http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Return identity if found
		i, err := GetIdentity(bs, ident)
		buf, err := json.Marshal(i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if i.Ident == "" {
				http.Error(w, "Ident not found", http.StatusNotFound)
			} else {
				fmt.Fprintf(w, "%s", buf)
			}
		}
	}
}

const ErrorMissingIdentifier = "400 - Bad Request (Are you missing an identifier in your URI?)"

func apiIdentPUTHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ident := vars["id"]

		if ident == "" {
			http.Error(w, ErrorMissingIdentifier, http.StatusBadRequest)
			return
		}

		// Update user
		var id hpfeeds.Identity
		if r.Body == nil {
			http.Error(w, "Request body required", http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ident != id.Ident {
			http.Error(w, "Request body not valid for this URI. Ident mismatch.", 400)
			return
		}

		err = SaveIdentity(bs, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", r.Body)

	}
}

func routes(bs BoltStore) *http.ServeMux {
	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler())
	r.HandleFunc("/api/ident/{id}", apiIdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/{id}", apiIdentPUTHandler(bs)).Methods("PUT")
	r.HandleFunc("/api/ident/{id}", apiIdentDELETEHandler(bs)).Methods("DELETE")

	s := http.NewServeMux()
	s.Handle("/", r)
	return s
}
