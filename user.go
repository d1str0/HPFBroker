package main

import (
	"github.com/alexedwards/argon2id"
)

type User struct {
	Name string
	Hash string // Will always be an encoding of a password hash
	Role string // RBAC role
}

// NewUser creates a user object with a hashed version of the passed in
// password.
func NewUser(name, password, role string) (*User, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}
	return &User{Name: name, Hash: hash, Role: role}, nil
}

// Authenticate takes a password as input, and compares the password hashes to
// determine if they should be authenticated.
func (u User) Authenticate(password string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, u.Hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
