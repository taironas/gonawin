package gonawintest

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"appengine/aetest"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
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
				AdminIds:    make([]int64, 1),
			},
			want: &Tournament{
				Name:            "Foo",
				Description:     "Foo description",
				Start:           time.Now(),
				End:             time.Now(),
				AdminIds:        make([]int64, 1),
				GroupIds:        make([]int64, 0),
				Matches1stStage: make([]int64, 0),
				Matches2ndStage: make([]int64, 0),
				UserIds:         make([]int64, 0),
				TeamIds:         make([]int64, 0),
			},
		},
	}
	for i, test := range tests {
		got, _ := CreateTournament(c, test.tournament.Name, test.tournament.Description, test.tournament.Start, test.tournament.End, test.tournament.AdminIds[0])
		if got == nil && test.want != nil {
			t.Errorf("TestCreateTournament(%q): got nil wanted %v", test.name, *test.want)
		} else if got != nil && test.want == nil {
			t.Errorf("TestCreateTournament(%q): got %v wanted nil", test.name, *got)
		} else if got == nil && test.want == nil {
			// This is OK
		} else if err = checkTournament(got, test.want); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
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
		want       *Tournament
	}{
		name: "destroy tournament",
		tournament: Tournament{
			Name:        "Foo",
			Description: "Foo description",
			Start:       time.Now(),
			End:         time.Now(),
			AdminIds:    make([]int64, 1),
		},
		want: nil,
	}

	// create tournament
	tournament, _ := CreateTournament(c, test.tournament.Name, test.tournament.Description, test.tournament.Start, test.tournament.End, test.tournament.AdminIds[0])

	// perform a get query so that the results of the unapplied write are visible to subsequent global queries.
	dummy := Tournament{}
	key := TournamentKeyById(c, tournament.Id)
	if err := datastore.Get(c, key, &dummy); err != nil {
		t.Fatal(err)
	}

	// destory it
	if got := tournament.Destroy(c); got != nil {
		t.Errorf("TestDestroyTournament(%q): got %v wanted %v", test.name, got, test.want)
	}

	// make a query on datastore to be sure it is not there.
	if got, err1 := TournamentById(c, tournament.Id); err1 == nil {
		t.Errorf("TestDestroyTournament(%q): got %v wanted %v", test.name, got, test.want)
	}
}

func TestFindTournaments(t *testing.T) {

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	log.Infof(c, "Test Find Tournament")

	tournaments := []Tournament{
		Tournament{
			Name:        "bar",
			Description: "Foo description",
			Start:       time.Now(),
			End:         time.Now(),
			AdminIds:    make([]int64, 1),
		},
		Tournament{
			Name:        "foobar",
			Description: "foo description",
			Start:       time.Now(),
			End:         time.Now(),
			AdminIds:    make([]int64, 1),
		},
		Tournament{
			Name:        "foobarfoo",
			Description: "Foo description",
			Start:       time.Now(),
			End:         time.Now(),
			AdminIds:    make([]int64, 1),
		},
	}

	tests := []struct {
		name        string
		tournaments []Tournament
		queries     []string
		want        struct {
			Len   int
			Names []string
		}
	}{
		{
			name:        "find tournaments in empty datastore",
			tournaments: make([]Tournament, 0),
			queries:     []string{"foo"},
			want: struct {
				Len   int
				Names []string
			}{
				Len:   int(0),
				Names: make([]string, 0),
			},
		},
		{
			name:        "find tournaments in full datastore",
			tournaments: tournaments,
			queries:     []string{"bar", "foobar", "foobarfoo"},
			want: struct {
				Len   int
				Names []string
			}{
				Len:   int(1),
				Names: []string{"bar", "foobar", "foobarfoo"},
			},
		},
	}

	// create tournaments
	for _, test := range tests {
		for i, _ := range test.tournaments {
			if got, err1 := CreateTournament(
				c,
				test.tournaments[i].Name,
				test.tournaments[i].Description,
				test.tournaments[i].Start,
				test.tournaments[i].End,
				test.tournaments[i].AdminIds[0]); err1 != nil {
				t.Errorf("TestFindTournaments(%q): error creating tournaments: %v", test.name, err1)
			} else {
				// perform a get query so that the results of the unapplied write are visible to subsequent global queries.
				dummy := Tournament{}
				key := TournamentKeyById(c, got.Id)
				if err := datastore.Get(c, key, &dummy); err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	// search tournaments
	log.Infof(c, "Test Find Tournament: start search ok tournaments")
	for _, test := range tests {
		for _, query := range test.queries {
			got := FindTournaments(c, "KeyName", query)
			if test.want.Len != len(got) {
				t.Errorf("TestFindTournaments(%q): got array of %v  wanted %v: query:%v t:%v", test.name, len(got), test.want.Len, query, got)
			}
			for _, tour := range got {
				if !helpers.SliceContains(test.want.Names, tour.Name) {
					t.Errorf("TestFindTournaments(%q): name not found. got %v  wanted among %v", test.name, tour.Name, test.want.Names)
				}
			}
		}
	}
}

func checkTournament(got *Tournament, want *Tournament) error {
	var s string
	if got.Name != want.Name {
		s = fmt.Sprintf("want Name == %s, got %s", want.Name, got.Name)
	} else if got.Description != want.Description {
		s = fmt.Sprintf("want Description == %s, got %s", want.Description, got.Description)
	} else if len(got.AdminIds) != len(want.AdminIds) {
		s = fmt.Sprintf("want AdminIds count == %d, got %d", len(want.AdminIds), len(got.AdminIds))
	} else if len(got.GroupIds) != len(want.GroupIds) {
		s = fmt.Sprintf("want GroupIds count == %d, got %d", len(want.GroupIds), len(got.GroupIds))
	} else if len(got.Matches1stStage) != len(want.Matches1stStage) {
		s = fmt.Sprintf("want Matches1stStage count == %d, got %d", len(want.Matches1stStage), len(got.Matches1stStage))
	} else if len(got.Matches2ndStage) != len(want.Matches2ndStage) {
		s = fmt.Sprintf("want Matches2ndStage count == %d, got %d", len(want.Matches2ndStage), len(got.Matches2ndStage))
	} else if len(got.UserIds) != len(want.UserIds) {
		s = fmt.Sprintf("want UserIds count == %d, got %d", len(want.UserIds), len(got.UserIds))
	} else if len(got.TeamIds) != len(want.TeamIds) {
		s = fmt.Sprintf("want TeamIds count == %d, got %d", len(want.TeamIds), len(got.TeamIds))
	} else {
		return nil
	}

	return errors.New(s)
}
