package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/d1str0/hpfeeds"
	bolt "go.etcd.io/bbolt"
)

func statusHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page of %s!", Version)
	}
}

//TODO return proer http codes per method
func apiIdentHandler(bs BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ident := r.URL.Path[len("/api/ident/"):]

		// Handle API requests depending on HTTP Method
		switch r.Method {
		case http.MethodGet:
			// Return identity if found
			i, err := GetIdentity(bs, ident)
			buf, err := json.Marshal(i)
			if err != nil {
				http.Error(w, err.Error(), 500)
			} else {
				if i.Ident == "" {
					http.Error(w, "Ident not found", 404)
				} else {
					fmt.Fprintf(w, "%s", buf)
				}
			}
		case http.MethodPut:
			// Update user
			var id hpfeeds.Identity
			if r.Body == nil {
				http.Error(w, "Request body required", 400)
				return
			}

			err := json.NewDecoder(r.Body).Decode(&id)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if ident != id.Ident {
				http.Error(w, "Request body not valid for this URI. Ident mismatch.", 400)
				return
			}

			err = SaveIdentity(bs, id)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		case http.MethodDelete:
			// Delete user
			err := DeleteIdentity(bs, ident)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}

}

func routes(bs BoltStore) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", statusHandler(bs))
	mux.HandleFunc("/api/ident/", apiIdentHandler(bs))
	return mux
}

// Used to identify a user and their identity within hpfeeds broker.
func GetIdentity(bs BoltStore, ident string) (*hpfeeds.Identity, error) {
	var i hpfeeds.Identity
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		v := b.Get([]byte(ident))
		err := json.Unmarshal(v, &i)
		return err
	})
	return &i, err
}

func SaveIdentity(bs BoltStore, id hpfeeds.Identity) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		buf, err := json.Marshal(id)
		b.Put([]byte(id.Ident), buf)
		return err
	})
	return err
}

func DeleteIdentity(bs BoltStore, ident string) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		b.Put([]byte(ident), nil)
		return nil
	})
	return err
}
