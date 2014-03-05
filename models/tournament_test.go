package models

import (
	"testing"
	"time"

	"appengine/aetest"

	"github.com/santiaago/purple-wing/helpers/log"
)

func TestCreateTournament(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	log.Infof(c, "Test Create Tournament")

	tests := []struct {
		name       string
		tournament Tournament
		want       *Tournament
	}{
		{
			name: "Simple create",
			tournament: Tournament{
				Name:        "Foo",
				Description: "Foo description",
				Start:       time.Now(),
				End:         time.Now(),
				AdminId:     int64(0),
			},
			want: &Tournament{
				Name:            "Foo",
				Description:     "Foo description",
				Start:           time.Now(),
				End:             time.Now(),
				AdminId:         int64(0),
				GroupIds:        make([]int64, 0),
				Matches1stStage: make([]int64, 0),
				Matches2ndStage: make([]int64, 0),
				UserIds:         make([]int64, 0),
				TeamIds:         make([]int64, 0),
			},
		},
	}
	for _, test := range tests {
		got, _ := CreateTournament(c, test.tournament.Name, test.tournament.Description, test.tournament.Start, test.tournament.End, test.tournament.AdminId)
		if got == nil && test.want != nil {
			t.Errorf("TestCreateTournament(%q): got nil wanted %v", test.name, *test.want)
		} else if got != nil && test.want == nil {
			t.Errorf("TestCreateTournament(%q): got %v wanted nil", test.name, *got)
		} else if got == nil && test.want == nil {
			// This is OK
		} else if got.Name != test.want.Name ||
			got.Description != test.want.Description ||
			got.AdminId != test.want.AdminId ||
			len(got.GroupIds) != len(test.want.GroupIds) ||
			len(got.Matches1stStage) != len(test.want.Matches1stStage) ||
			len(got.Matches2ndStage) != len(test.want.Matches2ndStage) ||
			len(got.UserIds) != len(test.want.UserIds) ||
			len(got.TeamIds) != len(test.want.TeamIds) {
			t.Errorf("TestCreateTournament(%q): got %v wanted %v", test.name, *got, *test.want)
		}
	}
}

func TestDestroyTournament(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	log.Infof(c, "Test Destory Tournament")

	test := struct {
		name       string
		tournament Tournament
		want *Tournament
	}{
		name: "destroy tournament",
		tournament:Tournament{
			Name:        "Foo",
			Description: "Foo description",
			Start:       time.Now(),
			End:         time.Now(),
			AdminId:     int64(0),
		},
		want:nil,
	}
	
	// create tournament
	tournament, _ := CreateTournament(c, test.tournament.Name, test.tournament.Description, test.tournament.Start, test.tournament.End, test.tournament.AdminId)
	
	// destory it
	if got := tournament.Destroy(c); got != nil{
		t.Errorf("TestDestroyTournament(%q): got %v wanted %v", test.name, got, test.want)
	}

	// make a query on datastore to be sure it is not there.
	if got, err1 := TournamentById(c, tournament.Id); err1 == nil{
		t.Errorf("TestDestroyTournament(%q): got %v wanted %v", test.name, got, test.want)
	}
}
