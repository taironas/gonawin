package models

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"appengine/aetest"

	"github.com/taironas/gonawin/helpers"
)

type testUser struct {
	email    string
	username string
	name     string
	alias    string
	isAdmin  bool
	auth     string
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
		s = fmt.Sprintf("want Alias == %s, got %s", want.alias, got.Alias)
	} else if got.IsAdmin != want.isAdmin {
		s = fmt.Sprintf("want isAdmin == %t, got %t", want.isAdmin, got.IsAdmin)
	} else {
		return nil
	}
	return errors.New(s)
}

// checkUserInvertedIndex checks that the user is present in the datastore when
// performing a search.
//
func checkUserInvertedIndex(t *testing.T, c aetest.Context, got *User, want testUser) error {

	var ids []int64
	var err error
	words := helpers.SetOfStrings(want.username)
	if ids, err = GetUserInvertedIndexes(c, words); err != nil {
		return fmt.Errorf("failed calling GetUserInvertedIndexes %v", err)
	}
	for _, id := range ids {
		if id == got.Id {
			return nil
		}
	}

	return errors.New("user not found")

}

func createNonSavedUser(email, username, name, alias string, isAdmin bool) User {
	return User{
		Id:                    5,
		Email:                 email,
		Username:              username,
		Name:                  name,
		Alias:                 alias,
		IsAdmin:               isAdmin,
		Auth:                  "",
		PredictIds:            []int64{},
		ArchivedPredictInds:   []int64{},
		TournamentIds:         []int64{},
		ArchivedTournamentIds: []int64{},
		TeamIds:               []int64{},
		Score:                 0,
		ScoreOfTournaments:    []ScoreOfTournament{},
		ActivityIds:           []int64{},
		Created:               time.Now(),
	}
}

func createUsersFromTestUsers(t *testing.T, c aetest.Context, testUsers []testUser) (userIDs []int64) {

	var err error
	for i, user := range testUsers {
		var got *User
		if got, err = CreateUser(c, user.email, user.username, user.name, user.alias, user.isAdmin, user.auth); err != nil {
			t.Errorf("user %d error: %v", i, err)
		}

		userIDs = append(userIDs, got.Id)
	}
	return
}

func createTestUsers(n int) (testUsers []testUser) {
	for i := 0; i < n; i++ {
		newUser := testUser{
			email:    fmt.Sprintf("foo%d@foo.com", i),
			username: fmt.Sprintf("foo_%d", i),
			name:     fmt.Sprintf("foo %d", i),
			alias:    fmt.Sprintf("alias foo %d", i),
			isAdmin:  false,
			auth:     "",
		}
		testUsers = append(testUsers, newUser)
	}
	return
}
