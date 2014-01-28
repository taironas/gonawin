/*
 * Copyright (c) 2013 Santiago Arias | Remy Jourde | Carlos Bernal
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package team

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
	teaminvidmdl "github.com/santiaago/purple-wing/models/teamInvertedIndex"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
)

type Team struct {
	Id      int64  `json:",omitempty"`
	KeyName string `json:",omitempty"`
	Name    string `json:",omitempty"`
	AdminId int64  `json:",omitempty"`
	Private bool
	Created time.Time `json:"-"`
}

type TeamCounter struct {
	Count int64
}

// Create a team given a name, an admin id and a private mode
func Create(c appengine.Context, name string, adminId int64, private bool) (*Team, error) {
	// create new team
	teamId, _, err := datastore.AllocateIDs(c, "Team", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "Team", "", teamId, nil)

	team := &Team{teamId, helpers.TrimLower(name), name, adminId, private, time.Now()}

	_, err = datastore.Put(c, key, team)
	if err != nil {
		return nil, err
	}
	// udpate inverted index
	teaminvidmdl.Add(c, helpers.TrimLower(name), teamId)
	// update team counter
	errIncrement := datastore.RunInTransaction(c, func(c appengine.Context) error {
		var err1 error
		_, err1 = incrementTeamCounter(c, datastore.NewKey(c, "TeamCounter", "singleton", 0, nil))
		return err1
	}, nil)
	if errIncrement != nil {
		log.Errorf(c, " Error incrementing TeamCounter")
	}

	return team, err
}

// Destroy a team given a team id
func Destroy(c appengine.Context, teamId int64) error {

	if team, err := ById(c, teamId); err != nil {
		return errors.New(fmt.Sprintf("Cannot find team with teamId=%d", teamId))
	} else {
		key := datastore.NewKey(c, "Team", "", team.Id, nil)

		return datastore.Delete(c, key)
	}
}

// given a filter and a value look query the datastore for teams and returns an array of team pointers.
func Find(c appengine.Context, filter string, value interface{}) []*Team {

	q := datastore.NewQuery("Team").Filter(filter+" =", value)

	var teams []*Team

	if _, err := q.GetAll(c, &teams); err == nil {
		return teams
	} else {
		log.Errorf(c, " Team.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

// get a team given an id
func ById(c appengine.Context, id int64) (*Team, error) {

	var t Team
	key := datastore.NewKey(c, "Team", "", id, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, " team not found : %v", err)
		return &t, err
	}
	return &t, nil
}

// get a team key given an id
func KeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Team", "", id, nil)

	return key
}

// update a team given an id and a team pointer
func Update(c appengine.Context, id int64, t *Team) error {
	// update key name
	t.KeyName = helpers.TrimLower(t.Name)
	k := KeyById(c, id)
	oldTeam := new(Team)
	if err := datastore.Get(c, k, oldTeam); err == nil {
		if _, err = datastore.Put(c, k, t); err != nil {
			return err
		}
		teaminvidmdl.Update(c, oldTeam.Name, t.Name, id)
	}

	return nil
}

// get all teams in datastore
func FindAll(c appengine.Context) []*Team {
	q := datastore.NewQuery("Team")

	var teams []*Team

	if _, err := q.GetAll(c, &teams); err != nil {
		log.Errorf(c, " Team.FindAll, error occurred during GetAll call: %v", err)
	}

	return teams
}

// find with respect to array of ids
func ByIds(c appengine.Context, ids []int64) []*Team {

	var teams []*Team
	for _, id := range ids {
		if team, err := ById(c, id); err == nil {
			teams = append(teams, team)
		} else {
			log.Errorf(c, " Team.ByIds, error occurred during ByIds call: %v", err)
		}
	}
	return teams
}

// checkes if a user has joined a team or not
func Joined(c appengine.Context, teamId int64, userId int64) bool {
	teamRel := teamrelmdl.FindByTeamIdAndUserId(c, teamId, userId)
	return teamRel != nil
}

// make a user join a team
func Join(c appengine.Context, teamId int64, userId int64) error {
	if _, err := teamrelmdl.Create(c, teamId, userId); err != nil {
		return errors.New(fmt.Sprintf(" Team.Join, error during team relationship creation: %v", err))
	}

	return nil
}

// make a user leave a team
// Todo: Should we check that the user is indeed a memeber of the team?
func Leave(c appengine.Context, teamId int64, userId int64) error {
	return teamrelmdl.Destroy(c, teamId, userId)
}

// check if user is admin of that team
func IsTeamAdmin(c appengine.Context, teamId int64, userId int64) bool {

	if team, err := ById(c, teamId); err == nil {
		return team.AdminId == userId
	} else {
		log.Errorf(c, " Team.IsTeamAdmin, error occurred during ById call: %v", err)
		return false
	}
}

// increment team counter to keep track of the number of teams present in datastore
func incrementTeamCounter(c appengine.Context, key *datastore.Key) (int64, error) {
	var x TeamCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count++
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

// decrement team counter to keep track of the number of teams present in datastore
func decrementTeamCounter(c appengine.Context, key *datastore.Key) (int64, error) {
	var x TeamCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count--
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

// get the current number of teams present in datastore
func GetTeamCounter(c appengine.Context) (int64, error) {
	key := datastore.NewKey(c, "TeamCounter", "singleton", 0, nil)
	var x TeamCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
}

// given a id, and a word, get the frequency of that word in the team terms
func GetWordFrequencyForTeam(c appengine.Context, id int64, word string) int64 {

	if teams := Find(c, "Id", id); teams != nil {
		return helpers.CountTerm(strings.Split(teams[0].KeyName, " "), word)
	}
	return 0
}
