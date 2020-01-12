package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/d1str0/hpfeeds"
	"github.com/gorilla/mux"
)

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
		i, err := bs.GetIdentity(ident)
		if err != nil {
			log.Printf("apiIdentDELETEHandler, GetIdentity(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If it doesn't already exist, return 404.
		if i == nil {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, ErrNotFound, http.StatusNotFound)
			return
		}

		err = bs.DeleteIdentity(ident)
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
				log.Printf("apiIdentGETHandler, GetAllIdentities(), %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(idents)
			return
		}

		// Return identity if found
		i, err := bs.GetIdentity(ident)
		buf, err := json.Marshal(i)
		if err != nil {
			log.Printf("apiIdentGETHandler, json.Marshal(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if i == nil {
				http.Error(w, ErrNotFound, http.StatusNotFound)
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
		i, err := bs.GetIdentity(ident)
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

		err = bs.SaveIdentity(id)
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
