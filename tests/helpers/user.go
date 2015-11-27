package gonawintest

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"appengine/aetest"

	"github.com/taironas/gonawin/helpers"
	mdl "github.com/taironas/gonawin/models"
)

// TestUser represents a testing user
type TestUser struct {
	Email    string
	Username string
	Name     string
	Alias    string
	IsAdmin  bool
	Auth     string
}

// CheckUser checks that a given user is equivalent to a test user
//
func CheckUser(got *mdl.User, want TestUser) error {
	var s string
	if got.Email != want.Email {
		s = fmt.Sprintf("want Email == %s, got %s", want.Email, got.Email)
	} else if got.Username != want.Username {
		s = fmt.Sprintf("want Username == %s, got %s", want.Username, got.Username)
	} else if got.Name != want.Name {
		s = fmt.Sprintf("want Name == %s, got %s", want.Name, got.Name)
	} else if got.Alias != want.Alias {
		s = fmt.Sprintf("want Alias == %s, got %s", want.Alias, got.Alias)
	} else if got.IsAdmin != want.IsAdmin {
		s = fmt.Sprintf("want isAdmin == %t, got %t", want.IsAdmin, got.IsAdmin)
	} else {
		return nil
	}
	return errors.New(s)
}

// CheckUserInvertedIndex checks that the user is present in the datastore when
// performing a search.
//
func CheckUserInvertedIndex(t *testing.T, c aetest.Context, got *mdl.User, want TestUser) error {

	var ids []int64
	var err error
	words := helpers.SetOfStrings(want.Username)
	if ids, err = mdl.GetUserInvertedIndexes(c, words); err != nil {
		return fmt.Errorf("failed calling GetUserInvertedIndexes %v", err)
	}
	for _, id := range ids {
		if id == got.Id {
			return nil
		}
	}

	return errors.New("user not found")

}

// CreateNonSavedUser creates a non datastored user
//
func CreateNonSavedUser(email, username, name, alias string, isAdmin bool) mdl.User {
	return mdl.User{
		5,
		email,
		username,
		name,
		alias,
		isAdmin,
		"",
		[]int64{},
		[]int64{},
		[]int64{},
		[]int64{},
		[]int64{},
		0,
		[]mdl.ScoreOfTournament{},
		[]int64{},
		time.Now(),
	}
}
