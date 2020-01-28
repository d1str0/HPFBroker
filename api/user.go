package api

import (
	hpf "github.com/d1str0/HPFBroker"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: Add validations for usernames and password complexity.

// A struct for parsing API input
type UserReq struct {
	Name     string
	Password string
	Role     string
}

// A struct for exporting User data via API responses
type UserResp struct {
	Name string
	Role string
}

func UserDELETEHandler(sc *hpf.ServerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		// DELETE /api/user/
		// Delete all users
		if id == "" {
			err := sc.DB.DeleteAllUsers()
			if err != nil {
				log.Printf("apiIdentDELETEHandler, DeleteAllIdentities(), %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Delete user
		u, err := sc.DB.GetUser(id)
		if err != nil {
			log.Printf("apiUserDELETEHandler, GetUser(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If it doesn't already exist, return 404.
		if u == nil {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, ErrNotFound, http.StatusNotFound)
			return
		}

		err = sc.DB.DeleteUser(id)
		if err != nil {
			log.Printf("apiUserDELETEHandler, DeleteUser(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

}

func UserGETHandler(sc *hpf.ServerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		// TODO: Factor out this section into new handler perhaps.
		// "/api/user/"
		if id == "" {
			users, err := sc.DB.GetAllUsers()
			if err != nil {
				log.Printf("apiUserGETHandler, GetAllUsers(), %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var urs []*UserResp
			for _, u := range users {
				ur := &UserResp{Name: u.Name, Role: u.Role}
				urs = append(urs, ur)
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(urs)
			return
		}

		// Return user if found
		u, err := sc.DB.GetUser(id)
		if err != nil {
			log.Printf("apiUserGETHandler, GetUser(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if u == nil {
			http.Error(w, ErrNotFound, http.StatusNotFound)
			return
		}

		// Make an appropriate response object (ie. no hash returned)
		ur := &UserResp{Name: u.Name, Role: u.Role}
		buf, err := json.Marshal(ur)
		if err != nil {
			log.Printf("apiUserGETHandler, json.Marshal(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", buf)
	}
}

func UserPUTHandler(sc *hpf.ServerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var create bool

		vars := mux.Vars(r)
		id := vars["id"]

		// Can't PUT on /user/ without an identifier.
		if id == "" {
			http.Error(w, ErrMissingIdentifier, http.StatusBadRequest)
			return
		}

		// Check to see if this user name already exists
		u, err := sc.DB.GetUser(id)
		if err != nil {
			log.Printf("apiUserPUTHandler, GetIdentity(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Is this a new User being created?
		// We want to remember so we know to return 200 vs 201
		if u == nil {
			create = true
		}

		if r.Body == nil {
			http.Error(w, ErrBodyRequired, http.StatusBadRequest)
			return
		}

		ureq := &UserReq{}
		err = json.NewDecoder(r.Body).Decode(ureq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if id != ureq.Name {
			http.Error(w, ErrMismatchedIdentifier, http.StatusBadRequest)
			return
		}

		u, err = hpf.NewUser(ureq.Name, ureq.Password, ureq.Role)
		if err != nil {
			log.Printf("apiUserPUTHandler, NewUser(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = sc.DB.SaveUser(u)
		if err != nil {
			log.Printf("apiUserPUTHandler, SaveUser(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if create {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		uresp := &UserResp{Name: u.Name, Role: u.Role}
		out, err := json.Marshal(uresp)
		if err != nil {
			log.Printf("apiUserPUTHandler, json.Marshal(), %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", out)
	}
}
