package main

import (
	"testing"
)

func Test_RBAC(t *testing.T) {
	r := rbac()

	// Test "admin" role

	if !r.IsGranted("admin", pRead, nil) {
		t.Error("Admin must be able to read")
	}
	if !r.IsGranted("admin", pWrite, nil) {
		t.Error("Admin must be able to write")
	}

	if !r.IsGranted("reader", pRead, nil) {
		t.Error("Reader must be able to read")
	}
	if r.IsGranted("reader", pWrite, nil) {
		t.Error("Reader must NOT be able to write")
	}

	if r.IsGranted("writer", pRead, nil) {
		t.Error("Writer must NOT be able to read")
	}
	if !r.IsGranted("writer", pWrite, nil) {
		t.Error("Writer must be able to write")
	}
}
