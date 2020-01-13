package main

import (
	"testing"
)

func TestUser(t *testing.T) {
	t.Run("NewUser", func(t *testing.T) {
		u, err := NewUser("name1", "pass1", "role1")
		if err != nil {
			t.Fatal(err)
		}

		if u == nil {
			t.Fatal("User should not be nil")
		}
	})

	t.Run("Authenticate", func(t *testing.T) {
		u, _ := NewUser("name1", "pass1", "role1")

		m, err := u.Authenticate("wrong")
		if err != nil {
			t.Fatal(err)
		}
		if m == true {
			t.Error("Authentication should have failed")
		}

		m, err = u.Authenticate("pass1")
		if err != nil {
			t.Fatal(err)
		}
		if m != true {
			t.Error("Authentication should have succeeded")
		}
	})
}
