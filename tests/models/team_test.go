package gonawintest

import (
	"strings"
	"testing"

	"appengine/aetest"

	mdl "github.com/taironas/gonawin/models"
	"github.com/taironas/gonawin/tests/helpers"
)

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
		team  gonawintest.TestTeam
	}{
		{
			title: "can create public team",
			team:  gonawintest.TestTeam{"my team", "description", 10, false},
		},
		{
			title: "can create private team",
			team:  gonawintest.TestTeam{"my other team", "description", 0, true},
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		var got *mdl.Team
		if got, err = mdl.CreateTeam(c, test.team.Name, test.team.Description, test.team.AdminID, test.team.Private); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = gonawintest.CheckTeam(got, test.team); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = gonawintest.CheckTeamInvertedIndex(t, c, got, test.team); err != nil {
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
		team       gonawintest.TestTeam
		overrideID bool
		newID      int64
		err        string
	}{
		{
			title: "can destroy team",
			team:  gonawintest.TestTeam{"my team", "description", 10, false},
		},
		{
			title:      "cannot destroy team",
			team:       gonawintest.TestTeam{"my team other team", "description", 10, false},
			overrideID: true,
			newID:      11,
			err:        "Cannot find team with Id",
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		var got *mdl.Team
		if got, err = mdl.CreateTeam(c, test.team.Name, test.team.Description, test.team.AdminID, test.team.Private); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}

		if test.overrideID {
			got.Id = test.newID
		}

		if err = got.Destroy(c); err != nil {
			if len(test.err) == 0 {
				t.Errorf("test %v - Error: %v", i, err)
			} else if !strings.Contains(gonawintest.ErrorString(err), test.err) {
				t.Errorf("test %v - Error: %v expected %v", i, err, test.err)
			}
		}

		var team *mdl.Team
		if team, err = mdl.TeamById(c, got.Id); team != nil {
			t.Errorf("test %v - Error: team found, not properly destroyed - %v", i, err)
		}

		if err = gonawintest.CheckTeamInvertedIndex(t, c, got, test.team); err == nil {
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
		teams []gonawintest.TestTeam
		query string
		want  int
	}{
		{
			title: "can find team",
			teams: []gonawintest.TestTeam{
				gonawintest.TestTeam{"my team", "description", 10, false},
				gonawintest.TestTeam{"my other team", "description", 10, false},
			},
			query: "my team",
			want:  1,
		},
		{
			title: "cannot find teams",
			teams: []gonawintest.TestTeam{
				gonawintest.TestTeam{"real", "description", 10, false},
				gonawintest.TestTeam{"barça", "description", 10, false},
			},
			query: "something else",
			want:  0,
		},
		{
			title: "can find multiple teams",
			teams: []gonawintest.TestTeam{
				gonawintest.TestTeam{"lakers", "description", 10, false},
				gonawintest.TestTeam{"lakers", "description", 10, false},
				gonawintest.TestTeam{"lakers", "description", 10, false},
			},
			query: "lakers",
			want:  3,
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		for _, team := range test.teams {
			if _, err = mdl.CreateTeam(c, team.Name, team.Description, team.AdminID, team.Private); err != nil {
				t.Errorf("test %v - Error: %v", i, err)
			}
		}

		var got []*mdl.Team
		if got = mdl.FindTeams(c, "Name", test.query); len(got) != test.want {
			t.Errorf("test %v - found %v teams expected %v with query %v by Name", i, test.want, len(got), test.query)
		}
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

	tTeam := gonawintest.TestTeam{"my team", "description", 10, false}

	var team *mdl.Team
	if team, err = mdl.CreateTeam(c, tTeam.Name, tTeam.Description, tTeam.AdminID, tTeam.Private); err != nil {
		t.Errorf("Error: %v", err)
	}

	tests := []struct {
		title  string
		id     int64
		wanted gonawintest.TestTeam
		err    string
	}{
		{
			title:  "can get team by Id",
			id:     team.Id,
			wanted: gonawintest.TestTeam{team.Name, team.Description, team.AdminIds[0], team.Private},
		},
		{
			title: "cannot get team by Id",
			id:    -1,
			err:   "no such entity",
		},
	}

	for i, test := range tests {
		t.Log(test.title)

		var got *mdl.Team
		if got, err = mdl.TeamById(c, test.id); err != nil {
			if len(test.err) == 0 {
				t.Errorf("test %v - Error: %v", i, err)
			} else if !strings.Contains(gonawintest.ErrorString(err), test.err) {
				t.Errorf("test %v - Error: %v expected %v", i, err, test.err)
			}
		} else if err = gonawintest.CheckTeam(got, test.wanted); err != nil {
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

		if got := mdl.TeamKeyById(c, test.id); got == nil {
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

	tTeam := gonawintest.TestTeam{"my team", "description", 10, false}

	var newTeam *mdl.Team
	if newTeam, err = mdl.CreateTeam(c, tTeam.Name, tTeam.Description, tTeam.AdminID, tTeam.Private); err != nil {
		t.Errorf("Error: %v", err)
	}

	tests := []struct {
		title      string
		id         int64
		updateTeam gonawintest.TestTeam
		overrideID bool
		newID      int64
		err        string
	}{
		{
			title:      "can update team",
			id:         newTeam.Id,
			updateTeam: gonawintest.TestTeam{Name: "updated team 1", Description: "updated description 1"},
		},
		{
			title:      "cannot update, team not found",
			id:         newTeam.Id,
			updateTeam: gonawintest.TestTeam{Name: "updated team 2", Description: "updated description 2"},
			overrideID: true,
			newID:      -1,
			err:        "no such entity",
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		var team *mdl.Team
		if team, err = mdl.TeamById(c, test.id); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}

		team.Name = test.updateTeam.Name
		team.Description = test.updateTeam.Description
		team.AdminIds[0] = test.updateTeam.AdminID
		team.Private = test.updateTeam.Private

		if test.overrideID {
			team.Id = test.newID
		}

		if err = team.Update(c); err != nil {
			if len(test.err) == 0 {
				t.Errorf("test %v - Error: %v", i, err)
			} else if !strings.Contains(gonawintest.ErrorString(err), test.err) {
				t.Errorf("test %v - Error: %v expected %v", i, err, test.err)
			}
			continue
		}

		var got *mdl.Team
		if got, err = mdl.TeamById(c, team.Id); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = gonawintest.CheckTeam(got, test.updateTeam); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = gonawintest.CheckTeamInvertedIndex(t, c, got, test.updateTeam); err != nil {
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

	testTeams := gonawintest.CreateTestTeams(3)
	teamIDs := gonawintest.CreateTeamsFromTestTeams(t, c, testTeams)

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
		teams   []gonawintest.TestTeam
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
			gonawintest.CreateTestTeams(3)[1:],
			"",
		},
		{
			"non existing teams for given IDs",
			teamIDsWithBadIDs,
			gonawintest.CreateTestTeams(0),
			"",
		},
	}

	for i, test := range tests {
		t.Log(test.title)

		var teams []*mdl.Team
		teams, err = mdl.TeamsByIds(c, test.teamIDs)

		if gonawintest.ErrorString(err) != test.err {
			t.Errorf("test %v error: want err: %s, got: %q", i, test.err, err)
		} else if test.err == "" && teams != nil {
			for i, team := range test.teams {
				if err = gonawintest.CheckTeam(teams[i], team); err != nil {
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
		keys := mdl.TeamsKeysByIds(c, test.ids)
		if len(keys) != len(test.ids) {
			t.Errorf("test %v: keys lenght does not match, expected: %v, got: %v", i, len(test.ids), len(keys))
		}
	}
}
