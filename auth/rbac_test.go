package auth

import (
	"testing"
)

func Test_RBAC(t *testing.T) {
	r := rbac()

	// Test "admin" role

	if !r.IsGranted("admin", PermRead, nil) {
		t.Error("Admin must be able to read")
	}
	if !r.IsGranted("admin", PermWrite, nil) {
		t.Error("Admin must be able to write")
	}

	if !r.IsGranted("reader", PermRead, nil) {
		t.Error("Reader must be able to read")
	}
	if r.IsGranted("reader", PermWrite, nil) {
		t.Error("Reader must NOT be able to write")
	}

	if r.IsGranted("writer", PermRead, nil) {
		t.Error("Writer must NOT be able to read")
	}
	if !r.IsGranted("writer", PermWrite, nil) {
		t.Error("Writer must be able to write")
	}
}
