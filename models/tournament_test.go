package models

import (
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
		want       *testTournament
	}{
		{
			name: "Simple create",
			tournament: Tournament{
				Name:        "Foo",
				Description: "Foo description",
				Start:       time.Now(),
				End:         time.Now(),
				AdminIds:    []int64{1},
			},
			want: &testTournament{
				name:        "Foo",
				description: "Foo description",
				start:       time.Now(),
				end:         time.Now(),
				adminID:     1,
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

func TestTournamentDestroy(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	log.Infof(c, "Test Destory Tournament")

	test := struct {
		name       string
		tournament Tournament
		want       *testTournament
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
	key := TournamentKeyByID(c, tournament.Id)
	if err := datastore.Get(c, key, &dummy); err != nil {
		t.Fatal(err)
	}

	// destroy it
	if got := tournament.Destroy(c); got != nil {
		t.Errorf("TestDestroyTournament(%q): got %v wanted %v", test.name, got, test.want)
	}

	// make a query on datastore to be sure it is not there.
	if got, err1 := TournamentByID(c, tournament.Id); err1 == nil {
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
		for i := range test.tournaments {
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
				key := TournamentKeyByID(c, got.Id)
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
