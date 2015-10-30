package models

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/test"

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
			} else if !strings.Contains(gonawintest.ErrorString(err), test.err) {
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

// TestFindTeams tests that you can find teams.
//
func TestFindTeams(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title string
		teams []testTeam
		query string
		want  int
	}{
		{
			title: "can find team",
			teams: []testTeam{
				testTeam{"my team", "description", 10, false},
				testTeam{"my other team", "description", 10, false},
			},
			query: "my team",
			want:  1,
		},
		{
			title: "cannot find teams",
			teams: []testTeam{
				testTeam{"real", "description", 10, false},
				testTeam{"barça", "description", 10, false},
			},
			query: "something else",
			want:  0,
		},
		{
			title: "can find multiple teams",
			teams: []testTeam{
				testTeam{"lakers", "description", 10, false},
				testTeam{"lakers", "description", 10, false},
				testTeam{"lakers", "description", 10, false},
			},
			query: "lakers",
			want:  3,
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		for _, team := range test.teams {
			if _, err = CreateTeam(c, team.name, team.description, team.adminId, team.private); err != nil {
				t.Errorf("test %v - Error: %v", i, err)
			}
		}

		var got []*Team
		if got = FindTeams(c, "Name", test.query); len(got) != test.want {
			t.Errorf("test %v - found %v teams expected %v with query %v by Name", i, test.want, len(got), test.query)
		}
	}
}

// TestFindAllTeams tests that you can find all teams in the datastore.
//
func TestFindAllTeams(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	testTeams := createTestTeams(10)
	teamIDs := createTeamsFromTestTeams(t, c, testTeams)

	got := FindAllTeams(c)
	if len(got) != len(teamIDs) {
		t.Errorf("lenght of expected(%v) and actual(%v) teams are different", len(teamIDs), len(got))
	}

}

// TestTeamById tests TeamById function.
//
func TestTeamById(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tTeam := testTeam{"my team", "description", 10, false}

	var team *Team
	if team, err = CreateTeam(c, tTeam.name, tTeam.description, tTeam.adminId, tTeam.private); err != nil {
		t.Errorf("Error: %v", err)
	}

	tests := []struct {
		title  string
		Id     int64
		wanted testTeam
		err    string
	}{
		{
			title:  "can get team by Id",
			Id:     team.Id,
			wanted: testTeam{team.Name, team.Description, team.AdminIds[0], team.Private},
		},
		{
			title: "cannot get team by Id",
			Id:    -1,
			err:   "no such entity",
		},
	}

	for i, test := range tests {
		t.Log(test.title)

		var got *Team
		if got, err = TeamById(c, test.Id); err != nil {
			if len(test.err) == 0 {
				t.Errorf("test %v - Error: %v", i, err)
			} else if !strings.Contains(gonawintest.ErrorString(err), test.err) {
				t.Errorf("test %v - Error: %v expected %v", i, err, test.err)
			}
		} else if err = checkTeam(got, test.wanted); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
	}
}

// TestTeamKeyById tests TeamKeyById function.
//
func TestTeamKeyById(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title string
		id    int64
	}{
		{
			title: "can get team Key by Id",
			id:    0,
		},
	}

	for i, test := range tests {
		t.Log(test.title)

		if got := TeamKeyById(c, test.id); got == nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
	}
}

// TestTeamUpdate tests team.Update function.
//
func TestTeamUpdate(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tTeam := testTeam{"my team", "description", 10, false}

	var newTeam *Team
	if newTeam, err = CreateTeam(c, tTeam.name, tTeam.description, tTeam.adminId, tTeam.private); err != nil {
		t.Errorf("Error: %v", err)
	}

	tests := []struct {
		title      string
		id         int64
		updateTeam testTeam
		overrideId bool
		newId      int64
		err        string
	}{
		{
			title:      "can update team",
			id:         newTeam.Id,
			updateTeam: testTeam{name: "updated team 1", description: "updated description 1"},
		},
		{
			title:      "cannot update, team not found",
			id:         newTeam.Id,
			updateTeam: testTeam{name: "updated team 2", description: "updated description 2"},
			overrideId: true,
			newId:      -1,
			err:        "no such entity",
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		var team *Team
		if team, err = TeamById(c, test.id); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}

		team.Name = test.updateTeam.name
		team.Description = test.updateTeam.description
		team.AdminIds[0] = test.updateTeam.adminId
		team.Private = test.updateTeam.private

		if test.overrideId {
			team.Id = test.newId
		}

		if err = team.Update(c); err != nil {
			if len(test.err) == 0 {
				t.Errorf("test %v - Error: %v", i, err)
			} else if !strings.Contains(gonawintest.ErrorString(err), test.err) {
				t.Errorf("test %v - Error: %v expected %v", i, err, test.err)
			}
			continue
		}

		var got *Team
		if got, err = TeamById(c, team.Id); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = checkTeam(got, test.updateTeam); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = checkTeamInvertedIndex(t, c, got, test.updateTeam); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
	}
}

// TestTeamsByIds tests that you can get a list of teams by their IDs.
//
func TestTeamsByIds(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	testTeams := createTestTeams(3)
	teamIDs := createTeamsFromTestTeams(t, c, testTeams)

	// Test data: only one bad team ID
	teamIDsWithOneBadID := make([]int64, len(teamIDs))
	copy(teamIDsWithOneBadID, teamIDs)
	teamIDsWithOneBadID[0] = teamIDsWithOneBadID[0] + 50

	// Test data: bad team IDs
	teamIDsWithBadIDs := make([]int64, len(teamIDs))
	copy(teamIDsWithBadIDs, teamIDs)
	teamIDsWithBadIDs[0] = teamIDsWithBadIDs[0] + 50
	teamIDsWithBadIDs[1] = teamIDsWithBadIDs[1] + 50
	teamIDsWithBadIDs[2] = teamIDsWithBadIDs[2] + 50

	tests := []struct {
		title   string
		teamIDs []int64
		teams   []testTeam
		err     string
	}{
		{
			"can get teams by IDs",
			teamIDs,
			testTeams,
			"",
		},
		{
			"can get all teams by IDs except one",
			teamIDsWithOneBadID,
			createTestTeams(3)[1:],
			"",
		},
		{
			"non existing teams for given IDs",
			teamIDsWithBadIDs,
			createTestTeams(0),
			"",
		},
	}

	for i, test := range tests {
		t.Log(test.title)

		var teams []*Team
		teams, err = TeamsByIds(c, test.teamIDs)

		if gonawintest.ErrorString(err) != test.err {
			t.Errorf("test %v error: want err: %s, got: %q", i, test.err, err)
		} else if test.err == "" && teams != nil {
			for i, team := range test.teams {
				if err = checkTeam(teams[i], team); err != nil {
					t.Errorf("test %v error: want team: %v, got: %v", i, team, teams[i])
				}
			}
		}
	}

}

// TestTeamsKeysByIds tests team.TeamsKeysByIds function.
//
func TestTeamsKeysByIds(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		ids []int64
	}{
		{[]int64{}},
		{[]int64{1, 2, 3, 4}},
	}

	for i, test := range tests {
		keys := TeamsKeysByIds(c, test.ids)
		if len(keys) != len(test.ids) {
			t.Errorf("test %v: keys lenght does not match, expected: %v, got: %v", i, len(test.ids), len(keys))
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
		if got, err = CreateTeam(c, team.name, team.description, team.adminId, team.private); err != nil {
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
			adminId:     10,
			private:     false,
		}
		testTeams = append(testTeams, newTeam)
	}
	return
}
