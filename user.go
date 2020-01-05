package main

type User struct {
	Name string
	Hash string // Will always be an encoding of a password hash
	Role string // RBAC role
}

