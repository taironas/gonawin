package tournaments

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/taironas/gonawin/helpers/handlers"
	"github.com/taironas/gonawin/test"
	"github.com/taironas/route"

	"appengine/aetest"
)

// UpdateMatchResult tests that result is properly updated.
//
func TestUpdateMatchResult(t *testing.T) {
	var i aetest.Instance
	var err error
	options := aetest.Options{StronglyConsistentDatastore: true}

	if i, err = aetest.NewInstance(&options); err != nil {
		t.Fatal(err)
	}
	defer i.Close()

	var c aetest.Context

	if c, err = aetest.NewContext(&options); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Create 5 users
	users := gonawintest.CreateUsers(c, 5)

	// Create a Tournament
	tournament := gonawintest.CreateWorldCup(c, users[0].Id)

	// Add 5 users to the Tournament and place a bet on a match
	for i, u := range users {
		if err = tournament.Join(c, u); err != nil {
			t.Errorf("Error: %v", err)
		}

		if err = gonawintest.PredictOnMatch(c, u, tournament.Matches1stStage[0], int64(i), int64(i+1)); err != nil {
			t.Errorf("Error: %v", err)
		}
	}

	r := new(route.Router)
	r.HandleFunc("/j/tournaments/:tournamentId/matches/:matchId/update", handlers.ErrorHandler(gonawintest.TestingAuthorized(UpdateMatchResult)))

	test := gonawintest.GenerateHandlerFunc(t, UpdateMatchResult)

	var parameters map[string]string
	parameters = make(map[string]string)

	parameters["result"] = "1 0"

	path := fmt.Sprintf("/j/tournaments/%d/matches/%d/update", tournament.Id, tournament.Matches1stStage[0])

	recorder := test(path, "POST", parameters, r)
	if recorder.Code != http.StatusOK {
		t.Errorf("returned %v. Expected %v.", recorder.Code, http.StatusOK)
	}
}
