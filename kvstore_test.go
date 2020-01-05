package main

import (
	"testing"

	"github.com/d1str0/hpfeeds"
)

func TestKvstore_BoltStore(t *testing.T) {
	bs := getTestDB(t)
	defer bs.Close()

	t.Run("Get Nonexistant", func(t *testing.T) {
		i, err := bs.GetIdentity("test")
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
		err := bs.SaveIdentity(id1)
		if err != nil {
			t.Fatal(err)
		}
		err = bs.SaveIdentity(id2)
		if err != nil {
			t.Fatal(err)
		}
		err = bs.SaveIdentity(id3)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get Existing", func(t *testing.T) {
		i, err := bs.GetIdentity("test-ident")
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
		i, err := bs.GetAllIdentities()
		if err != nil {
			t.Fatal(err)
		}
		if len(i) != 3 {
			t.Error("Expected 3 items at this point")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := bs.DeleteIdentity("test-ident")
		if err != nil {
			t.Fatal(err)
		}
		i, err := bs.GetIdentity("test-ident")
		if err != nil {
			t.Fatal(err)
		}
		if i != nil {
			t.Error("Expected nil returned after delete.")
		}

		// Should also work on non existent ident.
		err = bs.DeleteIdentity("test-ident4")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Delete All Identities", func(t *testing.T) {
		err := bs.DeleteAllIdentities()
		if err != nil {
			t.Fatal(err)
		}
		// Test by getting something that was there
		i, err := bs.GetIdentity("test-ident2")
		if err != nil {
			t.Fatal(err)
		}
		if i != nil {
			t.Error("Expected nil returned after delete all.")
		}
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
