package models

import (
	"errors"
	"fmt"
	"testing"

	"github.com/taironas/gonawin/helpers"

	"appengine/aetest"
)

type testUser struct {
	email    string
	username string
	name     string
	alias    string
	isAdmin  bool
	auth     string
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

	tests := []struct {
		title string
		user  testUser
	}{
		{"can create user", testUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""}},
	}

	for i, test := range tests {
		t.Log(test.title)
		var got *User
		if got, err = CreateUser(c, test.user.email, test.user.username, test.user.name, test.user.alias, test.user.isAdmin, test.user.auth); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = checkUser(got, test.user); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = checkUserInvertedIndex(t, c, got); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
	}
}

// TestUserById tests that you can get a user by its ID.
//
func TestUserById(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	test := struct {
		title string
		user  testUser
	}{
		"can get user by ID", testUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
	}

	t.Log(test.title)
	var got *User
	if got, err = CreateUser(c, test.user.email, test.user.username, test.user.name, test.user.alias, test.user.isAdmin, test.user.auth); err != nil {
		t.Errorf("Error: %v", err)
	}

	var u *User

	// Test non existing user
	if u, err = UserById(c, got.Id+50); u != nil {
		t.Errorf("Error: no user should have been found")
	}

	if err == nil {
		t.Errorf("Error: an error should have been returned in case of non existing user")
	}

	// Test existing user
	if u, err = UserById(c, got.Id); u == nil {
		t.Errorf("Error: user not found")
	}

	if err = checkUser(got, test.user); err != nil {
		t.Errorf("Error: want user == %v, got %v", test.user, got)
	}
}

// TestUsersByIds tests that you can get a list of users by their IDs.
//
func TestUsersByIds(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	test := struct {
		title string
		users []testUser
	}{
		"can get users by IDs",
		[]testUser{
			{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
			{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
			{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
		},
	}

	t.Log(test.title)
	var gotIDs []int64
	var got *User
	for _, user := range test.users {
		if got, err = CreateUser(c, user.email, user.username, user.name, user.alias, user.isAdmin, user.auth); err != nil {
			t.Errorf("Error: %v", err)
		}

		gotIDs = append(gotIDs, got.Id)
	}

	var users []*User

	// Test non existing users
	var nonExistingIDs []int64
	for _, ID := range gotIDs {
		nonExistingIDs = append(nonExistingIDs, ID+50)
	}

	if users, err = UsersByIds(c, nonExistingIDs); users != nil {
		t.Errorf("Error: no users should have been found")
	}

	// Test existing users
	if users, err = UsersByIds(c, gotIDs); users == nil {
		t.Errorf("Error: users not found")
	}

	for i, user := range test.users {
		if err = checkUser(users[i], user); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
	}
}

// TestDestroyUser tests that you can destroy a user.
//
func TestDestroyUser(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	test := struct {
		title string
		user  testUser
	}{
		"can destroy user", testUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
	}

	t.Log(test.title)
	var got *User
	if got, err = CreateUser(c, test.user.email, test.user.username, test.user.name, test.user.alias, test.user.isAdmin, test.user.auth); err != nil {
		t.Errorf("Error: %v", err)
	}

	if err = got.Destroy(c); err != nil {
		t.Errorf("Error: %v", err)
	}

	var u *User
	if u, err = UserById(c, got.Id); u != nil {
		t.Errorf("Error: user found, not properly destroyed")
	}
	if err = checkUserInvertedIndex(t, c, got); err == nil {
		t.Errorf("Error: user found in database")
	}
}

// TestFindUser tests that you can find a user.
//
func TestFindUser(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	test := struct {
		title string
		user  testUser
	}{
		"can find user", testUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
	}

	t.Log(test.title)

	if _, err = CreateUser(c, test.user.email, test.user.username, test.user.name, test.user.alias, test.user.isAdmin, test.user.auth); err != nil {
		t.Errorf("Error: %v", err)
	}

	var got *User
	if got = FindUser(c, "Username", "john.snow"); got == nil {
		t.Errorf("Error: user not found by Username")
	}

	if got = FindUser(c, "Name", "john snow"); got == nil {
		t.Errorf("Error: user not found by Name")
	}

	if got = FindUser(c, "Alias", "crow"); got == nil {
		t.Errorf("Error: user not found by Alias")
	}
}

// TestFindAllUsers tests that you can find all the users.
//
func TestFindAllUsers(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	test := struct {
		title string
		users []testUser
	}{
		"can find users",
		[]testUser{
			{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
			{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
			{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
		},
	}

	t.Log(test.title)

	for _, user := range test.users {
		if _, err = CreateUser(c, user.email, user.username, user.name, user.alias, user.isAdmin, user.auth); err != nil {
			t.Errorf("Error: %v", err)
		}
	}

	var got []*User
	if got = FindAllUsers(c); got == nil {
		t.Errorf("Error: users not found")
	}

	if len(got) != len(test.users) {
		t.Errorf("Error: want users count == %s, got %s", len(test.users), len(got))
	}

	for i, user := range test.users {
		if err = checkUser(got[i], user); err != nil {
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
		s = fmt.Sprintf("want isAdmin == %t, got %t", want.isAdmin, got.IsAdmin)
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
