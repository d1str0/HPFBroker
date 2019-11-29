package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/d1str0/hpfeeds"
)

func statusHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", Version)
	}
}

//TODO return proper http codes per method
func apiIdentHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Fix this line, it will panic if the request is missing the
		// trailing /
		ident := r.URL.Path[len("/api/ident/"):]

		// Handle API requests depending on HTTP Method
		switch r.Method {
		case http.MethodGet:
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
		case http.MethodPut:
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
		case http.MethodDelete:
			// Delete user
			err := DeleteIdentity(bs, ident)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

}

func routes(bs BoltStore) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", statusHandler())
	mux.HandleFunc("/api/ident/", apiIdentHandler(bs))
	return mux
}
