package models

import (
	"errors"
	"fmt"
	"testing"

	"github.com/santiaago/gonawin/helpers"

	"appengine/aetest"
)

type testUser struct {
	title    string
	email    string
	username string
	name     string
	alias    string
	isAdmin  bool
	auth     string
	err      string
}

// TestCreateUser tests that you can create a user.
//
func TestCreateUser(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []testUser{
		{
			title:    "can create user",
			email:    "foo@bar.com",
			username: "john.snow",
			name:     "john snow",
			alias:    "crow",
			isAdmin:  false,
			auth:     "",
			err:      "",
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		var got *User
		if got, err = CreateUser(c, test.email, test.username, test.name, test.alias, test.isAdmin, test.auth); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = checkUser(got, test); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = checkUserInvertedIndex(t, c, got); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
	}
}

func checkUser(got *User, want testUser) error {
	var s string
	if got.Email != want.email {
		s = fmt.Sprintf("want Email == %s, got %s", want.email, got.Email)
	} else if got.Username != want.username {
		s = fmt.Sprintf("want Username == %s, got %s", want.username, got.Username)
	} else if got.Name != want.name {
		s = fmt.Sprintf("want Name == %s, got %s", want.name, got.Name)
	} else if got.Alias != want.alias {
		s = fmt.Sprintf("want Name == %s, got %s", want.alias, got.Alias)
	} else if got.IsAdmin != want.isAdmin {
		s = fmt.Sprintf("want isAdmin == %s, got %s", want.isAdmin, got.IsAdmin)
	} else if got.Auth != want.auth {
		s = fmt.Sprintf("want auth == %s, got %s", want.auth, got.Auth)
	} else {
		return nil
	}
	return errors.New(s)
}

func checkUserInvertedIndex(t *testing.T, c aetest.Context, got *User) error {

	var ids []int64
	var err error
	words := helpers.SetOfStrings("john")
	if ids, err = GetUserInvertedIndexes(c, words); err != nil {
		s := fmt.Sprintf("failed calling GetUserInvertedIndexes %v", err)
		return errors.New(s)
	}
	for _, id := range ids {
		if id == got.Id {
			return nil
		}
	}

	return errors.New("user not found")

}
