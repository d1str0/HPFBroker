package auth

import (
	"testing"
)

func Test_RBAC(t *testing.T) {
	r := rbac()

	// Test "admin" role

	if !r.IsGranted(roleHPFReader, PermHPFRead, nil) {
		t.Error("HPF Reader must be able to read")
	}
	if r.IsGranted(roleHPFReader, PermHPFWrite, nil) {
		t.Error("HPF Reader must not be able to write")
	}

	if !r.IsGranted(roleHPFAdmin, PermHPFRead, nil) {
		t.Error("HPF Admin must be able to read")
	}
	if !r.IsGranted(roleHPFAdmin, PermHPFWrite, nil) {
		t.Error("HPF Admin must be able to write")
	}

	if !r.IsGranted(roleUserReader, PermUserRead, nil) {
		t.Error("User Reader must be able to read")
	}
	if r.IsGranted(roleUserReader, PermUserWrite, nil) {
		t.Error("User Reader must not be able to write")
	}

	if !r.IsGranted(roleUserAdmin, PermUserRead, nil) {
		t.Error("User Admin must be able to read")
	}
	if !r.IsGranted(roleUserAdmin, PermUserWrite, nil) {
		t.Error("User Admin must be able to write")
	}

	if !r.IsGranted(roleSuperAdmin, PermHPFRead, nil) {
		t.Error("Super Admin must be able to read hpfeeds")
	}
	if !r.IsGranted(roleSuperAdmin, PermHPFWrite, nil) {
		t.Error("Super Admin must be able to write hpfeeds")
	}
	if !r.IsGranted(roleSuperAdmin, PermUserRead, nil) {
		t.Error("Super Admin must be able to read users")
	}
	if !r.IsGranted(roleSuperAdmin, PermUserWrite, nil) {
		t.Error("Super Admin must be able to write users")
	}
}
