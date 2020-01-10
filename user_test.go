package main

import (
	"testing"
)

func TestUser(t *testing.T) {
	u, err := NewUser("name1", "pass1", "role1")
	if err != nil {
		t.Fatal(err)
	}

	if u == nil {
		t.Fatal("User should not be nil")
	}

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
}
