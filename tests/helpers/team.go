package gonawintest

import (
	"errors"
	"fmt"
	"testing"

	"appengine/aetest"

	"github.com/taironas/gonawin/helpers"
	mdl "github.com/taironas/gonawin/models"
)

// TestTeam represents a testing team
type TestTeam struct {
	Name        string
	Description string
	AdminID     int64
	Private     bool
}

// CheckTeam checks that the team passed has the same fields as the testTeam object.
//
func CheckTeam(got *mdl.Team, want TestTeam) error {
	var s string
	if got.Name != want.Name {
		s = fmt.Sprintf("want name == %s, got %s", want.Name, got.Name)
	} else if got.Description != want.Description {
		s = fmt.Sprintf("want Description == %s, got %s", want.Description, got.Description)
	} else if got.AdminIds[0] != want.AdminID {
		s = fmt.Sprintf("want AdminId == %d, got %d", want.AdminID, got.AdminIds[0])
	} else if got.Private != want.Private {
		s = fmt.Sprintf("want Private == %t, got %t", want.Private, got.Private)
	} else {
		return nil
	}
	return errors.New(s)
}

// CheckTeamInvertedIndex checks that the team is present in the datastore when
// performing a search.
//
func CheckTeamInvertedIndex(t *testing.T, c aetest.Context, got *mdl.Team, want TestTeam) error {

	var ids []int64
	var err error
	words := helpers.SetOfStrings(want.Name)
	if ids, err = mdl.GetTeamInvertedIndexes(c, words); err != nil {
		return fmt.Errorf("failed calling GetTeamInvertedIndexes %v", err)
	}
	for _, id := range ids {
		if id == got.Id {
			return nil
		}
	}

	return errors.New("team not found")
}

// CreateTeamsFromTestTeams creates saved teams from test teams
//
func CreateTeamsFromTestTeams(t *testing.T, c aetest.Context, testTeams []TestTeam) (teamIDs []int64) {

	var err error
	for i, team := range testTeams {
		var got *mdl.Team
		if got, err = mdl.CreateTeam(c, team.Name, team.Description, team.AdminID, team.Private); err != nil {
			t.Errorf("team %v error: %v", i, err)
		}

		teamIDs = append(teamIDs, got.Id)
	}
	return
}

// CreateTestTeams creates test teams
//
func CreateTestTeams(n int) (testTeams []TestTeam) {
	for i := 0; i < n; i++ {
		newTeam := TestTeam{
			Name:        fmt.Sprintf("team %v", i),
			Description: fmt.Sprintf("description %v", i),
			AdminID:     10,
			Private:     false,
		}
		testTeams = append(testTeams, newTeam)
	}
	return
}
