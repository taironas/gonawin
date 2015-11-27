package gonawintest

import (
	"fmt"
	"strings"
	"testing"
	"time"

	mdl "github.com/taironas/gonawin/models"
	"github.com/taironas/gonawin/tests/helpers"

	"appengine/aetest"
)

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
		user  gonawintest.TestUser
	}{
		{"can create user", gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""}},
	}

	for i, test := range tests {
		t.Log(test.title)
		var got *mdl.User
		if got, err = mdl.CreateUser(c, test.user.Email, test.user.Username, test.user.Name, test.user.Alias, test.user.IsAdmin, test.user.Auth); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = gonawintest.CheckUser(got, test.user); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if err = gonawintest.CheckUserInvertedIndex(t, c, got, test.user); err != nil {
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

	var u *mdl.User
	if u, err = mdl.CreateUser(c, "foo@bar.com", "john.snow", "john snow", "crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	tests := []struct {
		title  string
		userID int64
		user   gonawintest.TestUser
		err    string
	}{
		{"can get user by ID", u.Id, gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""}, ""},
		{"non existing user for given ID", u.Id + 50, gonawintest.TestUser{}, "datastore: no such entity"},
	}

	for _, test := range tests {
		t.Log(test.title)

		var got *mdl.User

		got, err = mdl.UserById(c, test.userID)

		if gonawintest.ErrorString(err) != test.err {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && got == nil {
			t.Errorf("Error: an user should have been found")
		} else if test.err == "" && got != nil {
			if err = gonawintest.CheckUser(got, test.user); err != nil {
				t.Errorf("Error: want user: %v, got: %v", test.user, got)
			}
		}
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

	// Test data: good user ID
	testUsers := []gonawintest.TestUser{
		{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
		{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
		{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
	}

	var gotIDs []int64

	for _, testUser := range testUsers {
		var got *mdl.User
		if got, err = mdl.CreateUser(c, testUser.Email, testUser.Username, testUser.Name, testUser.Alias, testUser.IsAdmin, testUser.Auth); err != nil {
			t.Errorf("Error: %v", err)
		}

		gotIDs = append(gotIDs, got.Id)
	}

	// Test data: only one bad user ID
	userIDsWithOneBadID := make([]int64, len(gotIDs))
	copy(userIDsWithOneBadID, gotIDs)
	userIDsWithOneBadID[0] = userIDsWithOneBadID[0] + 50

	// Test data: bad user IDs
	userIDsWithBadIDs := make([]int64, len(gotIDs))
	copy(userIDsWithBadIDs, gotIDs)
	userIDsWithBadIDs[0] = userIDsWithBadIDs[0] + 50
	userIDsWithBadIDs[1] = userIDsWithBadIDs[1] + 50
	userIDsWithBadIDs[2] = userIDsWithBadIDs[2] + 50

	tests := []struct {
		title   string
		userIDs []int64
		users   []gonawintest.TestUser
		err     string
	}{
		{
			"can get users by IDs",
			gotIDs,
			[]gonawintest.TestUser{
				{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
				{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
				{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
			},
			"",
		},
		{
			"can get all users by IDs except one",
			userIDsWithOneBadID,
			[]gonawintest.TestUser{
				{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
				{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
			},
			"",
		},
		{
			"non existing users for given IDs",
			userIDsWithBadIDs,
			[]gonawintest.TestUser{},
			"",
		},
	}

	for _, test := range tests {
		t.Log(test.title)

		var users []*mdl.User

		users, err = mdl.UsersByIds(c, test.userIDs)

		if gonawintest.ErrorString(err) != test.err {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && users != nil {
			for i, user := range test.users {
				if err = gonawintest.CheckUser(users[i], user); err != nil {
					t.Errorf("Error: want user: %v, got: %v", user, users[i])
				}
			}
		}
	}
}

// TestUserKeyById tests that you can get a user key by its ID.
//
func TestUserKeyById(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title  string
		userID int64
	}{
		{"can get user key by ID", 15},
	}

	for _, test := range tests {
		t.Log(test.title)

		key := mdl.UserKeyById(c, test.userID)

		if key.IntID() != test.userID {
			t.Errorf("Error: want key ID: %v, got: %v", test.userID, key.IntID())
		}
	}
}

// TestUserKeysByIds tests that you can get a list of user keys by their IDs.
//
func TestUserKeysByIds(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title   string
		userIDs []int64
	}{
		{
			"can get user keys by IDs",
			[]int64{25, 666, 2042},
		},
	}

	for _, test := range tests {
		t.Log(test.title)

		keys := mdl.UserKeysByIds(c, test.userIDs)

		if len(keys) != len(test.userIDs) {
			t.Errorf("Error: want number of user IDs: %d, got: %d", len(test.userIDs), len(keys))
		}

		for i, userID := range test.userIDs {
			if keys[i].IntID() != userID {
				t.Errorf("Error: want key ID: %d, got: %d", userID, keys[i].IntID())
			}
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
		user  gonawintest.TestUser
	}{
		"can destroy user", gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
	}

	t.Log(test.title)
	var got *mdl.User
	if got, err = mdl.CreateUser(c, test.user.Email, test.user.Username, test.user.Name, test.user.Alias, test.user.IsAdmin, test.user.Auth); err != nil {
		t.Errorf("Error: %v", err)
	}

	if err = got.Destroy(c); err != nil {
		t.Errorf("Error: %v", err)
	}

	var u *mdl.User
	if u, err = mdl.UserById(c, got.Id); u != nil {
		t.Errorf("Error: user found, not properly destroyed")
	}
	if err = gonawintest.CheckUserInvertedIndex(t, c, got, test.user); err == nil {
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
		user  gonawintest.TestUser
	}{
		"can find user", gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
	}

	t.Log(test.title)

	if _, err = mdl.CreateUser(c, test.user.Email, test.user.Username, test.user.Name, test.user.Alias, test.user.IsAdmin, test.user.Auth); err != nil {
		t.Errorf("Error: %v", err)
	}

	var got *mdl.User
	if got = mdl.FindUser(c, "Username", "john.snow"); got == nil {
		t.Errorf("Error: user not found by Username")
	}

	if got = mdl.FindUser(c, "Name", "john snow"); got == nil {
		t.Errorf("Error: user not found by Name")
	}

	if got = mdl.FindUser(c, "Alias", "crow"); got == nil {
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
		users []gonawintest.TestUser
	}{
		"can find users",
		[]gonawintest.TestUser{
			{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
			{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
			{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
		},
	}

	t.Log(test.title)

	for _, user := range test.users {
		if _, err = mdl.CreateUser(c, user.Email, user.Username, user.Name, user.Alias, user.IsAdmin, user.Auth); err != nil {
			t.Errorf("Error: %v", err)
		}
	}

	var got []*mdl.User
	if got = mdl.FindAllUsers(c); got == nil {
		t.Errorf("Error: users not found")
	}

	if len(got) != len(test.users) {
		t.Errorf("Error: want users count == %d, got %d", len(test.users), len(got))
	}

	for i, user := range test.users {
		if err = gonawintest.CheckUser(got[i], user); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
	}
}

// TestUserUpdate tests that you can update a user.
//
func TestUserUpdate(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	/*Test data: saved user*/
	var user *mdl.User
	if user, err = mdl.CreateUser(c, "foo@bar.com", "john.snow", "john snow", "crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	/*Test data: non saved user*/
	nonSavedUser := gonawintest.CreateNonSavedUser("foo@bar.com", "john.snow", "john snow", "crow", false)

	tests := []struct {
		title        string
		userToUpdate *mdl.User
		updatedUser  gonawintest.TestUser
		err          string
	}{
		{"update user successfully", user, gonawintest.TestUser{"foo@bar.com", "white.walkers", "white walkers", "dead", false, ""}, ""},
		{"update non saved user", &nonSavedUser, gonawintest.TestUser{"foo@bar.com", "white.walkers", "white walkers", "dead", false, ""}, ""},
	}

	for _, test := range tests {
		t.Log(test.title)

		test.userToUpdate.Username = test.updatedUser.Username
		test.userToUpdate.Name = test.updatedUser.Name
		test.userToUpdate.Alias = test.updatedUser.Alias

		err = test.userToUpdate.Update(c)

		updatedUser, _ := mdl.UserById(c, test.userToUpdate.Id)

		if gonawintest.ErrorString(err) != test.err {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && err != nil {
			t.Errorf("Error: user should have been properly updated")
		} else if test.err == "" && updatedUser != nil {
			if err = gonawintest.CheckUser(updatedUser, test.updatedUser); err != nil {
				t.Errorf("Error: want user: %v, got: %v", test.updatedUser, updatedUser)
			}
		}
	}
}

// TestUserSigninUser tests that you can signin a user.
//
func TestUserSigninUser(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title     string
		queryName string
		user      gonawintest.TestUser
		err       string
	}{
		{"can signin user with Email", "Email", gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "", false, ""}, ""},
		{"can signin user with Username", "Username", gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "", false, ""}, ""},
		{"cannot signin user", "Name", gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "", false, ""}, "no valid query name"},
	}

	for _, test := range tests {
		t.Log(test.title)

		var got *mdl.User

		got, err = mdl.SigninUser(c, test.queryName, test.user.Email, test.user.Username, test.user.Name)

		if !strings.Contains(gonawintest.ErrorString(err), test.err) {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && got == nil {
			t.Errorf("Error: an user should have been found")
		} else if test.err == "" && got != nil {
			if err = gonawintest.CheckUser(got, test.user); err != nil {
				t.Errorf("Error: want user: %v, got: %v", test.user, got)
			}
		}
	}
}

// TestUserTeams tests that you can get teams of a given user.
//
func TestUserTeams(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title       string
		user        gonawintest.TestUser
		teams       []gonawintest.TestTeam
		missingTeam bool
	}{
		{"can get teams",
			gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			[]gonawintest.TestTeam{
				{"night's watch", "guards of the wall", 10, false},
				{"Unsullied", "former slaves", 10, false},
				{"Wildlings", "we lived beyond the wall", 10, false},
			},
			false,
		},
		{"user with no team",
			gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			[]gonawintest.TestTeam{},
			false,
		},
		{"user with missing team",
			gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			[]gonawintest.TestTeam{
				{"night's watch", "guards of the wall", 10, false},
				{"Unsullied", "former slaves", 10, false},
				{"Wildlings", "we lived beyond the wall", 10, false},
			},
			true,
		},
	}

	for _, test := range tests {
		t.Log(test.title)

		var user *mdl.User
		if user, err = mdl.CreateUser(c, test.user.Email, test.user.Username, test.user.Name, test.user.Alias, test.user.IsAdmin, test.user.Auth); err != nil {
			t.Errorf("Error: %v", err)
		}

		for _, team := range test.teams {
			var newTeam *mdl.Team
			if newTeam, err = mdl.CreateTeam(c, team.Name, team.Description, team.AdminID, team.Private); err != nil {
				t.Errorf("Error: %v", err)
			}

			if err = newTeam.Join(c, user); err != nil {
				t.Errorf("Error: %v", err)
			}
		}

		if test.missingTeam {
			if err = user.AddTeamId(c, 666 /*extra team ID*/); err != nil {
				t.Errorf("Error: %v", err)
			}
		}

		var got []*mdl.Team
		got = user.Teams(c)

		if len(got) != len(test.teams) {
			t.Errorf("Error: want teams count == %d, got %d", len(test.teams), len(got))
		}

		for i, team := range test.teams {
			if err = gonawintest.CheckTeam(got[i], team); err != nil {
				t.Errorf("test %v - Error: %v", i, err)
			}
		}
	}
}

// TestTeamsByPage tests that you can get teams by page.
//
func TestTeamsByPage(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title          string
		user           gonawintest.TestUser
		paginatedTeams [][]gonawintest.TestTeam
		count          int64
		page           int64
	}{
		{
			title: "can get teams by page",
			user:  gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			paginatedTeams: [][]gonawintest.TestTeam{
				{
					{
						Name:        "night's watch",
						Description: "guards of the wall",
						AdminID:     10,
						Private:     false,
					},
				},
				{
					{
						Name:        "Unsullied",
						Description: "former slaves",
						AdminID:     10,
						Private:     false,
					},
					{
						Name:        "Wildlings",
						Description: "we lived beyond the wall",
						AdminID:     10,
						Private:     false,
					},
				},
			},
			count: 2,
			page:  2,
		},
	}

	for _, test := range tests {
		t.Log(test.title)

		var user *mdl.User
		if user, err = mdl.CreateUser(c, test.user.Email, test.user.Username, test.user.Name, test.user.Alias, test.user.IsAdmin, test.user.Auth); err != nil {
			t.Errorf("Error: %v", err)
		}

		for _, teams := range test.paginatedTeams {
			for _, team := range teams {
				var newTeam *mdl.Team
				if newTeam, err = mdl.CreateTeam(c, team.Name, team.Description, team.AdminID, team.Private); err != nil {
					t.Errorf("Error: %v", err)
				}

				if err = newTeam.Join(c, user); err != nil {
					t.Errorf("Error: %v", err)
				}
			}
		}

		for i := int64(1); i <= test.page; i++ {
			t.Log(fmt.Sprintf("test page %v", i))
			var got []*mdl.Team
			got = user.TeamsByPage(c, test.count, i)

			// pagination is reversted to creation order
			paginatedIndex := int64(len(test.paginatedTeams)) - i

			t.Log(fmt.Sprintf("expected teams %+v", test.paginatedTeams[paginatedIndex]))
			gotTeamsStr := fmt.Sprintf("got teams:\n")
			for _, tt := range got {
				gotTeamsStr += fmt.Sprintf("%+v\n", *tt)
			}
			t.Log(gotTeamsStr)

			if len(got) != len(test.paginatedTeams[paginatedIndex]) {
				t.Errorf("page %v Error: want teams count == %d, got %d", i, len(test.paginatedTeams), len(got))
			}

			for j, team := range test.paginatedTeams[paginatedIndex] {
				// pagination is reversted to creation order
				gotIndex := len(got) - j - 1
				if err = gonawintest.CheckTeam(got[gotIndex], team); err != nil {
					t.Errorf("page %v - Error: %v", i, err)
				}
			}
		}
	}
}

// TestTournamentsByPage tests that you can get tournaments by page.
//
func TestTournamentsByPage(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title                string
		user                 gonawintest.TestUser
		paginatedTournaments [][]gonawintest.TestTournament
		count                int64
		page                 int64
	}{
		{
			title: "can get tournaments by page",
			user:  gonawintest.TestUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			paginatedTournaments: [][]gonawintest.TestTournament{
				{
					{Name: "2014 FIFA World Cup", Description: "football world cup in Brazil", Start: time.Now(), End: time.Now(), AdminID: 1},
				},
				{
					{Name: "2018 FIFA World Cup", Description: "football world cup in Russia", Start: time.Now(), End: time.Now(), AdminID: 1},
					{Name: "2016 UEFA Euro", Description: "football euro in France", Start: time.Now(), End: time.Now(), AdminID: 1},
				},
			},
			count: 2,
			page:  2,
		},
	}

	for ti, test := range tests {
		t.Log(test.title)

		var user *mdl.User
		if user, err = mdl.CreateUser(c, test.user.Email, test.user.Username, test.user.Name, test.user.Alias, test.user.IsAdmin, test.user.Auth); err != nil {
			t.Errorf("test %v Error: %v", ti, err)
		}

		for pti, tournaments := range test.paginatedTournaments {
			for tsi, tournament := range tournaments {
				var newTournament *mdl.Tournament
				if newTournament, err = mdl.CreateTournament(c, tournament.Name, tournament.Description, tournament.Start, tournament.End, tournament.AdminID); err != nil {
					t.Errorf("test %v Error: %v", ti, err)
				}

				if err = newTournament.Join(c, user); err != nil {
					t.Errorf("test %v Error: %v", ti, err)
				}
				// need to upate userIds in test structure.
				// cannot go this before as we need to user.Id.
				test.paginatedTournaments[pti][tsi].UserIDs = []int64{user.Id}
			}
		}

		for i := int64(1); i <= test.page; i++ {
			var got []*mdl.Tournament
			got = user.TournamentsByPage(c, test.count, i)

			// pagination is reversed to creation order
			paginatedIndex := int64(len(test.paginatedTournaments)) - i

			if len(got) != len(test.paginatedTournaments[paginatedIndex]) {
				t.Errorf("test %v page %v Error: want tournaments count == %d, got %d", ti, i, len(test.paginatedTournaments[paginatedIndex]), len(got))
			}

			for j, tournament := range test.paginatedTournaments[paginatedIndex] {
				// pagination is reversed to creation order
				gotIndex := len(got) - j - 1
				if err = gonawintest.CheckTournament(got[gotIndex], tournament); err != nil {
					t.Errorf("test %v - page %v - Error: %v", ti, i, err)
				}
			}
		}
	}
}
