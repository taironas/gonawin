package gonawintest

import (
	"fmt"
	"testing"
	"time"

	mdl "github.com/taironas/gonawin/models"
)

// TestTournament represents a testing tournament
type TestTournament struct {
	name        string
	description string
	start       time.Time
	end         time.Time
	adminID     int64
}

// CreateTestTournaments creates n test tournaments
func CreateTestTournaments(n int) (testTournaments []TestTournament) {
	for i := 0; i < n; i++ {
		newTournament := TestTournament{
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

// CreateTournaments stores tournaments from test tournaments into the dastore
func CreateTournaments(t *testing.T, c aetest.Context, testTournaments []TestTournament) (tournamentIDs []int64) {
	var err error
	for i, tournament := range testTournaments {
		var got *mdl.Tournament
		if got, err = mdl.CreateTournament(c, tournament.name, tournament.description, tournament.start, tournament.end, tournament.adminID); err != nil {
			t.Errorf("tournament %v error: %v", i, err)
		}

		tournamentIDs = append(tournamentIDs, got.Id)
	}
	return
}

// CreateAndJoinTournaments stores tournaments from test tournaments into the dastore and join a given user
func CreateAndJoinTournaments(t *testing.T, c aetest.Context, testTournaments []TestTournament, user *mdl.User) (tournamentIDs []int64) {
	var err error
	for i, tournament := range testTournaments {
		var got *mdl.Tournament
		if got, err = mdl.CreateTournament(c, tournament.name, tournament.description, tournament.start, tournament.end, tournament.adminID); err != nil {
			t.Errorf("tournament %v error: %v", i, err)
		}

		if err = got.Join(c, user); err != nil {
			t.Errorf("tournament %v error: %v", i, err)
		}

		tournamentIDs = append(tournamentIDs, got.Id)
	}
	return
}
