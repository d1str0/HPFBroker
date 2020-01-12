package auth

import (
	"testing"
)

func Test_RBAC(t *testing.T) {
	r := rbac()

	// Test "admin" role

	t.Run("HPF Reader", func(t *testing.T) {
		if !r.IsGranted(RoleHPFReader, PermHPFRead, nil) {
			t.Error("HPF Reader must be able to read")
		}
		if r.IsGranted(RoleHPFReader, PermHPFWrite, nil) {
			t.Error("HPF Reader must not be able to write")
		}
	})

	t.Run("HPF Admin", func(t *testing.T) {
		if !r.IsGranted(RoleHPFAdmin, PermHPFRead, nil) {
			t.Error("HPF Admin must be able to read")
		}
		if !r.IsGranted(RoleHPFAdmin, PermHPFWrite, nil) {
			t.Error("HPF Admin must be able to write")
		}
	})

	t.Run("User Reader", func(t *testing.T) {
		if !r.IsGranted(RoleUserReader, PermUserRead, nil) {
			t.Error("User Reader must be able to read")
		}
		if r.IsGranted(RoleUserReader, PermUserWrite, nil) {
			t.Error("User Reader must not be able to write")
		}
	})

	t.Run("User Admin", func(t *testing.T) {
		if !r.IsGranted(RoleUserAdmin, PermUserRead, nil) {
			t.Error("User Admin must be able to read")
		}
		if !r.IsGranted(RoleUserAdmin, PermUserWrite, nil) {
			t.Error("User Admin must be able to write")
		}
	})

	t.Run("Super Admin", func(t *testing.T) {
		if !r.IsGranted(RoleSuperAdmin, PermHPFRead, nil) {
			t.Error("Super Admin must be able to read hpfeeds")
		}
		if !r.IsGranted(RoleSuperAdmin, PermHPFWrite, nil) {
			t.Error("Super Admin must be able to write hpfeeds")
		}
		if !r.IsGranted(RoleSuperAdmin, PermUserRead, nil) {
			t.Error("Super Admin must be able to read users")
		}
		if !r.IsGranted(RoleSuperAdmin, PermUserWrite, nil) {
			t.Error("Super Admin must be able to write users")
		}
	})
}
