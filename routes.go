package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/d1str0/hpfeeds"
	"github.com/gorilla/mux"
)

const ErrMissingIdentifier = "Missing identifier in URI"
const ErrMismatchedIdentifier = "URI doesn't match provided data"
const ErrBodyRequired = "Body is required for this endpoint"

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
			log.Printf("apiIdentDELETEHandler, DeleteIdentity(), %s", err.Error())
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
				log.Printf("apiIdentGETHandler, GetKeys(), %s", err.Error())
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
			log.Printf("apiIdentGETHandler, json.Marshal(), %s", err.Error())
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

		// Can't PUT on /ident/ without an identifier.
		if ident == "" {
			http.Error(w, ErrMissingIdentifier, http.StatusBadRequest)
			return
		}

		// Check to see if this ident already exists
		i, err := GetIdentity(bs, ident)
		if err != nil {
			log.Printf("apiIdentPUTHandler, GetIdentity(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Is this a new Ident being created?
		// We want to remember so we know to return 200 vs 201
		if i == nil {
			create = true
		}

		// Update user
		var id hpfeeds.Identity
		if r.Body == nil {
			http.Error(w, ErrBodyRequired, http.StatusBadRequest)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ident != id.Ident {
			http.Error(w, ErrMismatchedIdentifier, http.StatusBadRequest)
			return
		}

		err = SaveIdentity(bs, id)
		if err != nil {
			log.Printf("apiIdentPUTHandler, SaveIdentity(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if create {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		out, err := json.Marshal(id)
		if err != nil {
			log.Printf("apiIdentPUTHandler, json.Marshal(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", out)
	}
}

func router(bs BoltStore) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler())
	r.HandleFunc("/api/ident/", apiIdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/", apiIdentPUTHandler(bs)).Methods("PUT")
	r.HandleFunc("/api/ident/{id}", apiIdentGETHandler(bs)).Methods("GET")
	r.HandleFunc("/api/ident/{id}", apiIdentPUTHandler(bs)).Methods("PUT")
	r.HandleFunc("/api/ident/{id}", apiIdentDELETEHandler(bs)).Methods("DELETE")
	return r
}

func NewMux(bs BoltStore) *http.ServeMux {
	r := router(bs)

	s := http.NewServeMux()
	s.Handle("/", r)
	return s
}
