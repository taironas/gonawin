package models

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"appengine/aetest"

	"github.com/taironas/gonawin/test"
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
		if err = checkUserInvertedIndex(t, c, got, test.user); err != nil {
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

	var u *User
	if u, err = CreateUser(c, "foo@bar.com", "john.snow", "john snow", "crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	tests := []struct {
		title  string
		userID int64
		user   testUser
		err    string
	}{
		{"can get user by ID", u.Id, testUser{"foo@bar.com", "john.snow", "john snow", "crow", false, ""}, ""},
		{"non existing user for given ID", u.Id + 50, testUser{}, "datastore: no such entity"},
	}

	for _, test := range tests {
		t.Log(test.title)

		var got *User

		got, err = UserById(c, test.userID)

		if gonawintest.ErrorString(err) != test.err {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && got == nil {
			t.Errorf("Error: an user should have been found")
		} else if test.err == "" && got != nil {
			if err = checkUser(got, test.user); err != nil {
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
	testUsers := []testUser{
		{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
		{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
		{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
	}

	var gotIDs []int64

	for _, testUser := range testUsers {
		var got *User
		if got, err = CreateUser(c, testUser.email, testUser.username, testUser.name, testUser.alias, testUser.isAdmin, testUser.auth); err != nil {
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
		users   []testUser
		err     string
	}{
		{
			"can get users by IDs",
			gotIDs,
			[]testUser{
				{"foo@bar.com", "john.snow", "john snow", "crow", false, ""},
				{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
				{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
			},
			"",
		},
		{
			"can get all users by IDs except one",
			userIDsWithOneBadID,
			[]testUser{
				{"foo@bar.com", "robb.stark", "robb stark", "king in the north", false, ""},
				{"foo@bar.com", "jamie.lannister", "jamie lannister", "kingslayer", false, ""},
			},
			"",
		},
		{
			"non existing users for given IDs",
			userIDsWithBadIDs,
			[]testUser{},
			"",
		},
	}

	for _, test := range tests {
		t.Log(test.title)

		var users []*User

		users, err = UsersByIds(c, test.userIDs)

		if gonawintest.ErrorString(err) != test.err {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && users != nil {
			for i, user := range test.users {
				if err = checkUser(users[i], user); err != nil {
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

		key := UserKeyById(c, test.userID)

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

		keys := UserKeysByIds(c, test.userIDs)

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

// TestUserDestroy tests that you can destroy a user.
//
func TestUserDestroy(t *testing.T) {
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
	if err = checkUserInvertedIndex(t, c, got, test.user); err == nil {
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
		t.Errorf("Error: want users count == %d, got %d", len(test.users), len(got))
	}

	for i, user := range test.users {
		if err = checkUser(got[i], user); err != nil {
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
	var user *User
	if user, err = CreateUser(c, "foo@bar.com", "john.snow", "john snow", "crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	/*Test data: non saved user*/
	nonSavedUser := createNonSavedUser("foo@bar.com", "john.snow", "john snow", "crow", false)

	tests := []struct {
		title        string
		userToUpdate *User
		updatedUser  testUser
		err          string
	}{
		{"update user successfully", user, testUser{"foo@bar.com", "white.walkers", "white walkers", "dead", false, ""}, ""},
		{"update non saved user", &nonSavedUser, testUser{"foo@bar.com", "white.walkers", "white walkers", "dead", false, ""}, ""},
	}

	for _, test := range tests {
		t.Log(test.title)

		test.userToUpdate.Username = test.updatedUser.username
		test.userToUpdate.Name = test.updatedUser.name
		test.userToUpdate.Alias = test.updatedUser.alias

		err = test.userToUpdate.Update(c)

		updatedUser, _ := UserById(c, test.userToUpdate.Id)

		if gonawintest.ErrorString(err) != test.err {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && err != nil {
			t.Errorf("Error: user should have been properly updated")
		} else if test.err == "" && updatedUser != nil {
			if err = checkUser(updatedUser, test.updatedUser); err != nil {
				t.Errorf("Error: want user: %v, got: %v", test.updatedUser, updatedUser)
			}
		}
	}
}

// TestSigninUser tests that you can signin a user.
//
func TestSigninUser(t *testing.T) {
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
		user      testUser
		err       string
	}{
		{"can signin user with Email", "Email", testUser{"foo@bar.com", "john.snow", "john snow", "", false, ""}, ""},
		{"can signin user with Username", "Username", testUser{"foo@bar.com", "john.snow", "john snow", "", false, ""}, ""},
		{"cannot signin user", "Name", testUser{"foo@bar.com", "john.snow", "john snow", "", false, ""}, "no valid query name"},
	}

	for _, test := range tests {
		t.Log(test.title)

		var got *User

		got, err = SigninUser(c, test.queryName, test.user.email, test.user.username, test.user.name)

		if !strings.Contains(gonawintest.ErrorString(err), test.err) {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && got == nil {
			t.Errorf("Error: an user should have been found")
		} else if test.err == "" && got != nil {
			if err = checkUser(got, test.user); err != nil {
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
		user        testUser
		teams       []testTeam
		missingTeam bool
	}{
		{"can get teams",
			testUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			[]testTeam{
				{"night's watch", "guards of the wall", 10, false},
				{"Unsullied", "former slaves", 10, false},
				{"Wildlings", "we lived beyond the wall", 10, false},
			},
			false,
		},
		{"user with no team",
			testUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			[]testTeam{},
			false,
		},
		{"user with missing team",
			testUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			[]testTeam{
				{"night's watch", "guards of the wall", 10, false},
				{"Unsullied", "former slaves", 10, false},
				{"Wildlings", "we lived beyond the wall", 10, false},
			},
			true,
		},
	}

	for _, test := range tests {
		t.Log(test.title)

		var user *User
		if user, err = CreateUser(c, test.user.email, test.user.username, test.user.name, test.user.alias, test.user.isAdmin, test.user.auth); err != nil {
			t.Errorf("Error: %v", err)
		}

		for _, team := range test.teams {
			var newTeam *Team
			if newTeam, err = CreateTeam(c, team.name, team.description, team.adminID, team.private); err != nil {
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

		var got []*Team
		got = user.Teams(c)

		if len(got) != len(test.teams) {
			t.Errorf("Error: want teams count == %d, got %d", len(test.teams), len(got))
		}

		for i, team := range test.teams {
			if err = checkTeam(got[i], team); err != nil {
				t.Errorf("test %v - Error: %v", i, err)
			}
		}
	}
}

// TestUserTeamsByPage tests that you can get teams by page.
//
func TestUserTeamsByPage(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title          string
		user           testUser
		paginatedTeams [][]testTeam
		count          int64
		page           int64
	}{
		{
			title: "can get teams by page",
			user:  testUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			paginatedTeams: [][]testTeam{
				{
					{
						name:        "night's watch",
						description: "guards of the wall",
						adminID:     10,
						private:     false,
					},
				},
				{
					{
						name:        "Unsullied",
						description: "former slaves",
						adminID:     10,
						private:     false,
					},
					{
						name:        "Wildlings",
						description: "we lived beyond the wall",
						adminID:     10,
						private:     false,
					},
				},
			},
			count: 2,
			page:  2,
		},
	}

	for _, test := range tests {
		t.Log(test.title)

		var user *User
		if user, err = CreateUser(c, test.user.email, test.user.username, test.user.name, test.user.alias, test.user.isAdmin, test.user.auth); err != nil {
			t.Errorf("Error: %v", err)
		}

		for _, teams := range test.paginatedTeams {
			for _, team := range teams {
				var newTeam *Team
				if newTeam, err = CreateTeam(c, team.name, team.description, team.adminID, team.private); err != nil {
					t.Errorf("Error: %v", err)
				}

				if err = newTeam.Join(c, user); err != nil {
					t.Errorf("Error: %v", err)
				}
			}
		}

		for i := int64(1); i <= test.page; i++ {
			t.Log(fmt.Sprintf("test page %v", i))
			var got []*Team
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
				if err = checkTeam(got[gotIndex], team); err != nil {
					t.Errorf("page %v - Error: %v", i, err)
				}
			}
		}
	}
}

// TestUserTournamentsByPage tests that you can get tournaments by page.
//
func TestUserTournamentsByPage(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title                string
		user                 testUser
		paginatedTournaments [][]testTournament
		count                int64
		page                 int64
	}{
		{
			title: "can get tournaments by page",
			user:  testUser{"foo@bar.com", "john.snow", "john snow", "", false, ""},
			paginatedTournaments: [][]testTournament{
				{
					{name: "2014 FIFA World Cup", description: "football world cup in Brazil", start: time.Now(), end: time.Now(), adminID: 1},
				},
				{
					{name: "2018 FIFA World Cup", description: "football world cup in Russia", start: time.Now(), end: time.Now(), adminID: 1},
					{name: "2016 UEFA Euro", description: "football euro in France", start: time.Now(), end: time.Now(), adminID: 1},
				},
			},
			count: 2,
			page:  2,
		},
	}

	for ti, test := range tests {
		t.Log(test.title)

		var user *User
		if user, err = CreateUser(c, test.user.email, test.user.username, test.user.name, test.user.alias, test.user.isAdmin, test.user.auth); err != nil {
			t.Errorf("test %v Error: %v", ti, err)
		}

		for pti, tournaments := range test.paginatedTournaments {
			for tsi, tournament := range tournaments {
				var newTournament *Tournament
				if newTournament, err = CreateTournament(c, tournament.name, tournament.description, tournament.start, tournament.end, tournament.adminID); err != nil {
					t.Errorf("test %v Error: %v", ti, err)
				}

				if err = newTournament.Join(c, user); err != nil {
					t.Errorf("test %v Error: %v", ti, err)
				}
				// need to upate userIds in test structure.
				// cannot go this before as we need to user.Id.
				test.paginatedTournaments[pti][tsi].userIDs = []int64{user.Id}
			}
		}

		for i := int64(1); i <= test.page; i++ {
			t.Log(fmt.Sprintf("test page %v", i))
			var got []*Tournament
			got = user.TournamentsByPage(c, test.count, i)

			// pagination is reversed to creation order
			paginatedIndex := int64(len(test.paginatedTournaments)) - i

			t.Log(fmt.Sprintf("expected tournaments %+v", test.paginatedTournaments[paginatedIndex]))
			gotTournamentsStr := fmt.Sprintf("got tournaments:\n")
			for _, tt := range got {
				gotTournamentsStr += fmt.Sprintf("%+v\n", *tt)
			}
			t.Log(gotTournamentsStr)

			if len(got) != len(test.paginatedTournaments[paginatedIndex]) {
				t.Errorf("test %v page %v Error: want tournaments count == %d, got %d", ti, i, len(test.paginatedTournaments[paginatedIndex]), len(got))
			}

			for j, tournament := range test.paginatedTournaments[paginatedIndex] {
				// pagination is reversed to creation order
				gotIndex := len(got) - j - 1
				if err = checkTournament(got[gotIndex], &tournament); err != nil {
					t.Errorf("test %v - page %v - Error: %v", ti, i, err)
				}
			}
		}
	}
}

// TestUserAddPredictId tests that predict ID is well added to a user entity.
//
func TestUserAddPredictId(t *testing.T) {

	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title     string
		predictID int64
		err       string
	}{
		{
			"can add predict ID to user",
			42,
			"",
		},
	}

	var user *User
	if user, err = CreateUser(c, "john.snow@winterfell.com", "john.snow", "John Snow", "Crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	for _, test := range tests {
		t.Log(test.title)

		err = user.AddPredictId(c, test.predictID)

		if !strings.Contains(gonawintest.ErrorString(err), test.err) {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && user.PredictIds[0] != test.predictID {
			t.Errorf("Error: a predict ID should have been retrieved from the user")
		}
	}
}

// TestUserAddTournamentId tests that tournament ID is well added to a user entity.
//
func TestUserAddTournamentId(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title        string
		tournamentID int64
		err          string
	}{
		{
			"can add tournament ID to user",
			42,
			"",
		},
		{
			"cannot add twice same tournament ID to user",
			42,
			"AddTournamentId, allready a member",
		},
	}

	var user *User
	if user, err = CreateUser(c, "john.snow@winterfell.com", "john.snow", "John Snow", "Crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	for _, test := range tests {
		t.Log(test.title)
		err = user.AddTournamentId(c, test.tournamentID)

		if !strings.Contains(gonawintest.ErrorString(err), test.err) {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && user.TournamentIds[0] != test.tournamentID {
			t.Errorf("Error: a tournament ID should have been retrieved from the user")
		}
	}
}

// TestUserContainsTournamentId tests if a tournament ID exists for a user entity.
//
func TestUserContainsTournamentId(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title        string
		tournamentID int64
		contains     bool
		index        int
	}{
		{
			"contains tournament ID from user",
			42,
			true,
			0,
		},
		{
			"does not contain tournament ID from user",
			54,
			false,
			-1,
		},
	}

	var user *User
	if user, err = CreateUser(c, "john.snow@winterfell.com", "john.snow", "John Snow", "Crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	if err = user.AddTournamentId(c, tests[0].tournamentID); err != nil {
		t.Errorf("Error: %v", err)
	}

	for _, test := range tests {
		t.Log(test.title)

		contains, index := user.ContainsTournamentId(test.tournamentID)

		if contains != test.contains {
			t.Errorf("Error: want contains: %t, got: %t", test.contains, contains)
		} else if index != test.index {
			t.Errorf("Error: want index: %d, got: %d", test.index, index)
		}
	}
}

// TestUserRemoveTournamentId tests that tournament ID is well removed from a user entity.
//
func TestUserRemoveTournamentId(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title        string
		tournamentID int64
		err          string
	}{
		{
			"can remove tournament ID from user",
			42,
			"",
		},
		{
			"cannot remove tournament ID from user",
			54,
			"RemoveTournamentId, not a member",
		},
	}

	var user *User
	if user, err = CreateUser(c, "john.snow@winterfell.com", "john.snow", "John Snow", "Crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	if err = user.AddTournamentId(c, tests[0].tournamentID); err != nil {
		t.Errorf("Error: %v", err)
	}

	for _, test := range tests {
		t.Log(test.title)

		err = user.RemoveTournamentId(c, test.tournamentID)

		contains, _ := user.ContainsTournamentId(test.tournamentID)

		if !strings.Contains(gonawintest.ErrorString(err), test.err) {
			t.Errorf("Error: want err: %s, got: %q", test.err, err)
		} else if test.err == "" && contains {
			t.Errorf("Error: tournament IDs should be empty")
		}
	}
}

// TestUserContainsTeamId tests if a team ID exists for a user entity.
//
func TestUserContainsTeamId(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		title    string
		teamID   int64
		contains bool
		index    int
	}{
		{
			"contains team ID from user",
			42,
			true,
			0,
		},
		{
			"does not contain team ID from user",
			54,
			false,
			-1,
		},
	}

	var user *User
	if user, err = CreateUser(c, "john.snow@winterfell.com", "john.snow", "John Snow", "Crow", false, ""); err != nil {
		t.Errorf("Error: %v", err)
	}

	if err = user.AddTeamId(c, tests[0].teamID); err != nil {
		t.Errorf("Error: %v", err)
	}

	for _, test := range tests {
		t.Log(test.title)

		contains, index := user.ContainsTeamId(test.teamID)

		if contains != test.contains {
			t.Errorf("Error: want contains: %t, got: %t", test.contains, contains)
		} else if index != test.index {
			t.Errorf("Error: want index: %d, got: %d", test.index, index)
		}
	}
}
