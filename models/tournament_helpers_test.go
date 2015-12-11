package models

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"appengine/aetest"
)

type testTournament struct {
	name        string
	description string
	start       time.Time
	end         time.Time
	adminID     int64
	userIDs     []int64
}

// checkTournament checks that a given tournament is equivalent to a test tournament
func checkTournament(got *Tournament, want *testTournament) error {
	var s string
	if got.Name != want.name {
		s = fmt.Sprintf("want Name == %s, got %s", want.name, got.Name)
	} else if got.Description != want.description {
		s = fmt.Sprintf("want Description == %s, got %s", want.description, got.Description)
	} else if got.Start.Sub(want.start).Hours() > 0 {
		s = fmt.Sprintf("want Start == %v, got %v", want.start, got.Start)
	} else if got.End.Sub(want.end).Hours() > 0 {
		s = fmt.Sprintf("want End == %v, got %v", want.end, got.End)
	} else if got.AdminIds[0] != want.adminID {
		s = fmt.Sprintf("want AdminId == %d, got %d", want.adminID, got.AdminIds[0])
	} else if len(got.UserIds) != len(want.userIDs) {
		s = fmt.Sprintf("want UserIds count == %d, got %d", len(want.userIDs), len(got.UserIds))
	} else {
		return nil
	}

	return errors.New(s)
}

// createTestTournaments creates n test tournaments
func createTestTournaments(n int) (testTournaments []*testTournament) {
	for i := 0; i < n; i++ {
		newTournament := &testTournament{
			name:        fmt.Sprintf("tournament %v", i),
			description: fmt.Sprintf("description %v", i),
			start:       time.Now(),
			end:         time.Now(),
			adminID:     10,
		}
		testTournaments = append(testTournaments, newTournament)
	}
	return
}

// createTournaments stores tournaments from test tournaments into the dastore
func createTournaments(t *testing.T, c aetest.Context, testTournaments []testTournament) (tournamentIDs []int64) {
	var err error
	for i, tournament := range testTournaments {
		var got *Tournament
		if got, err = CreateTournament(c, tournament.name, tournament.description, tournament.start, tournament.end, tournament.adminID); err != nil {
			t.Errorf("tournament %v error: %v", i, err)
		}

		tournamentIDs = append(tournamentIDs, got.Id)
	}
	return
}

// createAndJoinTournaments stores tournaments from test tournaments into the dastore and join a given user
func createAndJoinTournaments(t *testing.T, c aetest.Context, testTournaments []*testTournament, user *User) (tournamentIDs []int64) {
	var err error
	for i, tournament := range testTournaments {
		var got *Tournament
		if got, err = CreateTournament(c, tournament.name, tournament.description, tournament.start, tournament.end, tournament.adminID); err != nil {
			t.Errorf("tournament %v error: %v", i, err)
		}

		if err = got.Join(c, user); err != nil {
			t.Errorf("tournament %v error: %v", i, err)
		}

		tournamentIDs = append(tournamentIDs, got.Id)
	}
	return
}

func addUserIDToTournaments(tournaments *[]*testTournament, userID int64) {
	for _, t := range *tournaments {
		t.userIDs = append(t.userIDs, userID)
	}
}
