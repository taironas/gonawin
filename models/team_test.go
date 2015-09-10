package models

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/taironas/gonawin/helpers"

	"appengine/aetest"
)

type testTeam struct {
	name        string
	description string
	adminId     int64
	private     bool
}

// TestCreateTeam tests that you can create a team.
//
func TestCreateTeam(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title string
		team  testTeam
	}{
		{
			title: "can create public team",
			team:  testTeam{"my team", "description", 10, false},
		},
		{
			title: "can create private team",
			team:  testTeam{"my other team", "description", 0, true},
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		var got *Team
		if got, err = CreateTeam(c, test.team.name, test.team.description, test.team.adminId, test.team.private); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = checkTeam(got, test.team); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = checkTeamInvertedIndex(t, c, got, test.team); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
	}
}

// TestDestroyTeam test that you can destroy a team.
//
func TestDestroyTeam(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title      string
		team       testTeam
		overrideId bool
		newId      int64
		err        string
	}{
		{
			title: "can destroy team",
			team:  testTeam{"my team", "description", 10, false},
		},
		{
			title:      "cannot destroy team",
			team:       testTeam{"my team other team", "description", 10, false},
			overrideId: true,
			newId:      11,
			err:        "Cannot find team with Id",
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		var got *Team
		if got, err = CreateTeam(c, test.team.name, test.team.description, test.team.adminId, test.team.private); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}

		if test.overrideId {
			got.Id = test.newId
		}

		if err = got.Destroy(c); err != nil {
			if len(test.err) == 0 {
				t.Errorf("test %v - Error: %v", i, err)
			} else if !strings.Contains(errString(err), test.err) {
				t.Errorf("test %v - Error: %v expected %v", i, err, test.err)
			}
		}

		var team *Team
		if team, err = TeamById(c, got.Id); team != nil {
			t.Errorf("test %v - Error: team found, not properly destroyed - %v", i, err)
		}

		if err = checkTeamInvertedIndex(t, c, got, test.team); err == nil {
			t.Errorf("test %v - Error: team found in database", i)
		}
	}
}

// checkTeam checks that the team passed has the same fields as the testTeam object.
//
func checkTeam(got *Team, want testTeam) error {
	var s string
	if got.Name != want.name {
		s = fmt.Sprintf("want name == %s, got %s", want.name, got.Name)
	} else if got.Description != want.description {
		s = fmt.Sprintf("want Description == %s, got %s", want.description, got.Description)
	} else if got.AdminIds[0] != want.adminId {
		s = fmt.Sprintf("want AdminId == %s, got %s", want.adminId, got.AdminIds[0])
	} else if got.Private != want.private {
		s = fmt.Sprintf("want Private == %s, got %s", want.private, got.Private)
	} else {
		return nil
	}
	return errors.New(s)
}

// checkTeamInvertedIndex checks that the team is present in the datastore when
// performing a search.
//
func checkTeamInvertedIndex(t *testing.T, c aetest.Context, got *Team, want testTeam) error {

	var ids []int64
	var err error
	words := helpers.SetOfStrings(want.name)
	if ids, err = GetTeamInvertedIndexes(c, words); err != nil {
		return errors.New(fmt.Sprintf("failed calling GetTeamInvertedIndexes %v", err))
	}
	for _, id := range ids {
		if id == got.Id {
			return nil
		}
	}

	return errors.New("team not found")
}

// errString returns the string representation of an error.
func errString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
