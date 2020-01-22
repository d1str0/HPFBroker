package main

import (
	"testing"

	"github.com/d1str0/hpfeeds"
)

func TestKvstore_BoltStore(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	t.Run("IDENTITIES", func(t *testing.T) {
		t.Run("Get Nonexistant", func(t *testing.T) {
			i, err := db.GetIdentity("test")
			if err != nil {
				t.Fatal(err)
			}

			if i != nil {
				t.Error("Expected nil identity returned.")
			}
		})

		id1 := hpfeeds.Identity{Ident: "test-ident", Secret: "test-secret", SubChannels: []string{}, PubChannels: []string{}}
		id2 := hpfeeds.Identity{Ident: "test-ident2", Secret: "test-secret", SubChannels: []string{}, PubChannels: []string{}}
		id3 := hpfeeds.Identity{Ident: "test-ident3", Secret: "test-secret", SubChannels: []string{}, PubChannels: []string{}}

		t.Run("Save Identity", func(t *testing.T) {
			err := db.SaveIdentity(id1)
			if err != nil {
				t.Fatal(err)
			}
			err = db.SaveIdentity(id2)
			if err != nil {
				t.Fatal(err)
			}
			err = db.SaveIdentity(id3)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Get Existing", func(t *testing.T) {
			i, err := db.GetIdentity("test-ident")
			if err != nil {
				t.Fatal(err)
			}

			if i == nil {
				t.Error("Unexpected nil identity returned.")
			}

			// expected, got
			assertEqualIdentity(t, id1, *i)
		})

		t.Run("Get All Identities", func(t *testing.T) {
			i, err := db.GetAllIdentities()
			if err != nil {
				t.Fatal(err)
			}
			if len(i) != 3 {
				t.Error("Expected 3 items at this point")
			}
		})

		t.Run("Delete", func(t *testing.T) {
			err := db.DeleteIdentity("test-ident")
			if err != nil {
				t.Fatal(err)
			}
			i, err := db.GetIdentity("test-ident")
			if err != nil {
				t.Fatal(err)
			}
			if i != nil {
				t.Error("Expected nil returned after delete.")
			}

			// Should also work on non existent ident.
			err = db.DeleteIdentity("test-ident4")
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Delete All Identities", func(t *testing.T) {
			err := db.DeleteAllIdentities()
			if err != nil {
				t.Fatal(err)
			}
			// Test by getting something that was there
			i, err := db.GetIdentity("test-ident2")
			if err != nil {
				t.Fatal(err)
			}
			if i != nil {
				t.Error("Expected nil returned after delete all.")
			}
		})
	})

	t.Run("USERS", func(t *testing.T) {
		t.Run("Get Nonexistant", func(t *testing.T) {
			u, err := db.GetUser("test")
			if err != nil {
				t.Fatal(err)
			}

			if u != nil {
				t.Error("Expected nil identity returned.")
			}
		})

		u1 := User{Name: "test-name", Hash: "test-hash", Role: "admin"}
		u2 := User{Name: "test-name2", Hash: "test-hash", Role: "admin"}
		u3 := User{Name: "test-name3", Hash: "test-hash", Role: "admin"}

		t.Run("Save User", func(t *testing.T) {
			err := db.SaveUser(u1)
			if err != nil {
				t.Fatal(err)
			}
			err = db.SaveUser(u2)
			if err != nil {
				t.Fatal(err)
			}
			err = db.SaveUser(u3)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Get Existing", func(t *testing.T) {
			u, err := db.GetUser("test-name")
			if err != nil {
				t.Fatal(err)
			}

			if u == nil {
				t.Error("Unexpected nil user returned.")
			}

			// expected, got
			assertEqualUser(t, u1, *u)
		})

		t.Run("Get All Users", func(t *testing.T) {
			u, err := db.GetAllUsers()
			if err != nil {
				t.Fatal(err)
			}
			if len(u) != 3 {
				t.Error("Expected 3 items at this point")
			}
		})

		t.Run("Delete", func(t *testing.T) {
			err := db.DeleteIdentity("test-user")
			if err != nil {
				t.Fatal(err)
			}
			u, err := db.GetIdentity("test-user")
			if err != nil {
				t.Fatal(err)
			}
			if u != nil {
				t.Error("Expected nil returned after delete.")
			}

			// Should also work on non existent ident.
			err = db.DeleteIdentity("test-user4")
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Delete All Users", func(t *testing.T) {
			err := db.DeleteAllUsers()
			if err != nil {
				t.Fatal(err)
			}
			// Test by getting something that was there
			u, err := db.GetIdentity("test-user2")
			if err != nil {
				t.Fatal(err)
			}
			if u != nil {
				t.Error("Expected nil returned after delete all.")
			}
		})
	})

}

func assertEqualIdentity(t *testing.T, expect hpfeeds.Identity, got hpfeeds.Identity) {
	if expect.Ident != got.Ident {
		t.Errorf("Mismatched Idents:\n\tgot %s \n\twant %s", got.Ident, expect.Ident)
	}
	if expect.Secret != got.Secret {
		t.Errorf("Mismatched Secrets:\n\tgot %s \n\twant %s", got.Secret, expect.Secret)
	}
	if !testEq(expect.SubChannels, got.SubChannels) {
		t.Errorf("Mismatched SubChannels:\n\tgot %v \n\twant %v", got.SubChannels, expect.SubChannels)
	}
	if !testEq(expect.PubChannels, got.PubChannels) {
		t.Errorf("Mismatched PubChannels:\n\tgot %v \n\twant %v", got.PubChannels, expect.PubChannels)
	}
}

func assertEqualUser(t *testing.T, expect User, got User) {
	if expect.Name != got.Name {
		t.Errorf("Mismatched Names:\n\tgot %s \n\twant %s", got.Name, expect.Name)
	}
	if expect.Hash != got.Hash {
		t.Errorf("Mismatched Hashes:\n\tgot %s \n\twant %s", got.Hash, expect.Hash)
	}
	if expect.Role != got.Role {
		t.Errorf("Mismatched Roles:\n\tgot %s \n\twant %s", got.Role, expect.Role)
	}
}

func testEq(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
