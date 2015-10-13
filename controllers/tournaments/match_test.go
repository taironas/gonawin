package tournaments

import (
	"testing"

	"appengine/aetest"

	"github.com/taironas/gonawin/test"
)

// UpdateMatchResult tests that result is properly updated.
//
func TestUpdateMatchResult(t *testing.T) {
	var c aetest.Context
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Create 5 users
	users := gonawintest.CreateUsers(c, 5)

	// Create a Tournament
	tournament := gonawintest.CreateWorldCup(c, users[0].Id)

	// add 5 users to the Tournament and place a bet on a match
	for i, u := range users {
		if err = tournament.Join(c, u); err != nil {
			t.Errorf("Error: %v", err)
		}

		if err = gonawintest.PredictOnMatch(c, u, tournament.Matches1stStage[0], int64(i), int64(i+1)); err != nil {
			t.Errorf("Error: %v", err)
		}
	}

	// update the match
}
