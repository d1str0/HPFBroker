package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/d1str0/hpfeeds"
	"github.com/gorilla/mux"
)

const ErrorMissingIdentifier = "400 - Bad Request (Are you missing an identifier in your URI?)"
const ErrorMethodNotAllowed = "405 - Method Not Allowed"

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
			fmt.Printf("DELETE error")
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
			keys, err := bs.GetKeys()
			if err != nil {
				fmt.Printf("GET error")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			for _, v := range keys {
				fmt.Fprintf(w, "%s\n", v)
			}
			return
		}

		// Return identity if found
		i, err := GetIdentity(bs, ident)
		buf, err := json.Marshal(i)
		if err != nil {
			fmt.Printf("GET error 2")
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

func apiIdentPUTHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var create bool

		vars := mux.Vars(r)
		ident := vars["id"]

		if ident == "" {
			http.Error(w, ErrorMissingIdentifier, http.StatusBadRequest)
			return
		}

		i, err := GetIdentity(bs, ident)
		if err != nil {
			fmt.Printf("PUT error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if i == nil {
			create = true
		}

		// Update user
		var id hpfeeds.Identity
		if r.Body == nil {
			http.Error(w, "Request body required", http.StatusBadRequest)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&id)
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
			fmt.Print("Error decoding json")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if create {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		fmt.Fprintf(w, "%s", r.Body)
	}
}

func routes(bs BoltStore) *http.ServeMux {
	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler())
	r.HandleFunc("/api/ident/", apiIdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/{id}", apiIdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/{id}", apiIdentPUTHandler(bs)).Methods("PUT")
	r.HandleFunc("/api/ident/{id}", apiIdentDELETEHandler(bs)).Methods("DELETE")

	s := http.NewServeMux()
	s.Handle("/", r)
	return s
}
