package api

import (
	hpf "github.com/d1str0/HPFBroker"
	auth "github.com/d1str0/HPFBroker/auth"

	"net/http"
	"strings"
)

var ErrAuthFailed = "Failed to authenticate"
var ErrAuthInvalidCreds = "Invalid credentials"

// AuthHandler parses a Basic Authentication request.
func AuthHandler(sc *hpf.ServerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		// Make sure this is of the format "Basic {credentials}"
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, ErrAuth, http.StatusBadRequest)
			return
		}

		credentials, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			// TODO: Provide better response.
			http.Error(w, ErrAuthFailed, http.StatusBadRequest)
			return
		}

		// SplitN because we only want to split on the first ":" as a password
		// may contain special characters.
		pair := strings.SplitN(credentials, ":", 2)

		// Validate we have a username and password
		if len(pair) != 2 {
			// TODO: Provide better response.
			http.Error(w, ErrAuthFailed, http.StatusBadRequest)
			return
		}

		name, pass := pair[0], pair[1]

		u, err := sc.db.GetUser(name)
		if err != nil {
			// TODO: Provide better response. (Here we might not want to be more
			// specific to avoid username enumeration)
			http.Error(w, ErrAuthFailed, http.StatusUnauthorized)
			return
		}

		match, err := u.Authenticate(pass)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
