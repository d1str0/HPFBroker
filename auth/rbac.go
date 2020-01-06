package auth

import (
	"github.com/mikespook/gorbac"
)

// Create separate read and write permissions
var (
	PermRead  = gorbac.NewStdPermission("read")
	PermWrite = gorbac.NewStdPermission("write")

	// Setup our main role
	roleAdmin = "admin"

	// These two are for testing primarily
	roleReader = "reader"
	roleWriter = "writer"
)

// rbac returns a new instance of gorbac.RBAC for Role-Based Access Controls.
func rbac() *gorbac.RBAC {
	r := gorbac.New()

	ra := gorbac.NewStdRole(roleAdmin)
	rr := gorbac.NewStdRole(roleReader)
	rw := gorbac.NewStdRole(roleWriter)

	// Assign both permissions to admin role
	ra.Assign(PermRead)
	ra.Assign(PermWrite)

	rr.Assign(PermRead)
	rw.Assign(PermWrite)

	// Add role to rbac instance
	r.Add(ra)
	r.Add(rr)
	r.Add(rw)

	return r
}
