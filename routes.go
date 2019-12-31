package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/d1str0/hpfeeds"
	"github.com/gorilla/mux"
)

const ErrMissingIdentifier = "Missing identifier in URI"          // 400
const ErrMismatchedIdentifier = "URI doesn't match provided data" // 400
const ErrBodyRequired = "Body is required for this endpoint"      // 400

const ErrIdentNotFound = "Ident not found" // 404

func statusHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", Version)
	}
}

func apiIdentDELETEHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ident := vars["id"]

		// DELETE /api/ident/
		// Delete all identities
		if ident == "" {
			err := bs.DeleteAllIdentities()
			if err != nil {
				log.Printf("apiIdentDELETEHandler, DeleteAllIdentities(), %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Delete user
		i, err := GetIdentity(bs, ident)
		if err != nil {
			log.Printf("apiIdentDELETEHandler, GetIdentity(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If it doesn't already exist, return 404.
		if i == nil {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, ErrIdentNotFound, http.StatusNotFound)
			return
		}

		err = DeleteIdentity(bs, ident)
		if err != nil {
			log.Printf("apiIdentDELETEHandler, DeleteIdentity(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

}

func apiIdentGETHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ident := vars["id"]

		// "/api/ident/"
		if ident == "" {
			idents, err := bs.GetAllIdentities()
			if err != nil {
				log.Printf("apiIdentGETHandler, GetKeys(), %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(idents)
			return
		}

		// Return identity if found
		i, err := GetIdentity(bs, ident)
		buf, err := json.Marshal(i)
		if err != nil {
			log.Printf("apiIdentGETHandler, json.Marshal(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if i == nil {
				http.Error(w, "Ident not found", http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
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
	r.HandleFunc("/api/ident/", apiIdentPUTHandler(bs)).Methods("PUT") // Funnel bad request for proper response.
	r.HandleFunc("/api/ident/", apiIdentDELETEHandler(bs)).Methods("DELETE")
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
