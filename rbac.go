package main

import (
	"github.com/mikespook/gorbac"
)

// Create separate read and write permissions
var (
	pRead  = gorbac.NewStdPermission("read")
	pWrite = gorbac.NewStdPermission("write")
)

// rbac returns a new instance of gorbac.RBAC for Role-Based Access Controls.
func rbac() *gorbac.RBAC {
	r := gorbac.New()

	// Setup our main role
	roleAdmin := gorbac.NewStdRole("admin")

	// These two are for testing primarily
	roleReader := gorbac.NewStdRole("reader")
	roleWriter := gorbac.NewStdRole("writer")

	// Create separate read and write permissions
	pRead := gorbac.NewStdPermission("read")
	pWrite := gorbac.NewStdPermission("write")

	// Assign both permissions to admin role
	roleAdmin.Assign(pRead)
	roleAdmin.Assign(pWrite)

	roleReader.Assign(pRead)
	roleWriter.Assign(pWrite)

	// Add role to rbac instance
	r.Add(roleAdmin)
	r.Add(roleReader)
	r.Add(roleWriter)

	return r
}
