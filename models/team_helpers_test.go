package models

import (
	"errors"
	"fmt"
	"testing"

	"github.com/taironas/gonawin/helpers"

	"appengine/aetest"
)

type testTeam struct {
	name        string
	description string
	adminID     int64
	private     bool
}

// checkTeam checks that the team passed has the same fields as the testTeam object.
//
func checkTeam(got *Team, want testTeam) error {
	var s string
	if got.Name != want.name {
		s = fmt.Sprintf("want name == %s, got %s", want.name, got.Name)
	} else if got.Description != want.description {
		s = fmt.Sprintf("want Description == %s, got %s", want.description, got.Description)
	} else if got.AdminIds[0] != want.adminID {
		s = fmt.Sprintf("want AdminId == %d, got %d", want.adminID, got.AdminIds[0])
	} else if got.Private != want.private {
		s = fmt.Sprintf("want Private == %t, got %t", want.private, got.Private)
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
		return fmt.Errorf("failed calling GetTeamInvertedIndexes %v", err)
	}
	for _, id := range ids {
		if id == got.Id {
			return nil
		}
	}

	return errors.New("team not found")
}

func createTeamsFromTestTeams(t *testing.T, c aetest.Context, testTeams []testTeam) (teamIDs []int64) {

	var err error
	for i, team := range testTeams {
		var got *Team
		if got, err = CreateTeam(c, team.name, team.description, team.adminID, team.private); err != nil {
			t.Errorf("team %v error: %v", i, err)
		}

		teamIDs = append(teamIDs, got.Id)
	}
	return
}

func createTestTeams(n int) (testTeams []testTeam) {
	for i := 0; i < n; i++ {
		newTeam := testTeam{
			name:        fmt.Sprintf("team %v", i),
			description: fmt.Sprintf("description %v", i),
			adminID:     10,
			private:     false,
		}
		testTeams = append(testTeams, newTeam)
	}
	return
}
