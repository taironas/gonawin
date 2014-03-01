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

package models

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
)

type Team struct {
	Id      int64
	KeyName string
	Name    string
	AdminId int64
	Private bool
	Created time.Time
	UserIds []int64
}

type TeamJson struct {
	Id      *int64     `json:",omitempty"`
	KeyName *string    `json:",omitempty"`
	Name    *string    `json:",omitempty"`
	AdminId *int64     `json:",omitempty"`
	Private *bool      `json:",omitempty"`
	Created *time.Time `json:",omitempty"`
	UserIds *[]int64   `json:",omitempty"`
}

// Create a team given a name, an admin id and a private mode.
func CreateTeam(c appengine.Context, name string, adminId int64, private bool) (*Team, error) {
	// create new team
	teamId, _, err := datastore.AllocateIDs(c, "Team", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "Team", "", teamId, nil)
	emtpyArray := make([]int64, 0)

	team := &Team{teamId, helpers.TrimLower(name), name, adminId, private, time.Now(), emtpyArray}

	_, err = datastore.Put(c, key, team)
	if err != nil {
		return nil, err
	}
	// udpate inverted index
	teaminvidmdl.Add(c, helpers.TrimLower(name), teamId)

	return team, err
}

// Destroy a team given a team id.
func (t *Team) Destroy(c appengine.Context) error {

	if _, err := TeamById(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf("Cannot find team with Id=%d", t.Id))
	} else {
		key := datastore.NewKey(c, "Team", "", t.Id, nil)

		return datastore.Delete(c, key)
	}
}

// Given a filter and a value look query the datastore for teams and returns an array of team pointers.
func FindTeams(c appengine.Context, filter string, value interface{}) []*Team {

	q := datastore.NewQuery("Team").Filter(filter+" =", value)

	var teams []*Team

	if _, err := q.GetAll(c, &teams); err == nil {
		return teams
	} else {
		log.Errorf(c, " Team.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

// Get a team given an id.
func TeamById(c appengine.Context, id int64) (*Team, error) {

	var t Team
	key := datastore.NewKey(c, "Team", "", id, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, " team not found : %v", err)
		return &t, err
	}
	return &t, nil
}

// get a team key given an id
func TeamKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Team", "", id, nil)

	return key
}

// Update a team given an id and a team pointer.
func (t *Team) Update(c appengine.Context) error {
	// update key name
	t.KeyName = helpers.TrimLower(t.Name)
	k := TeamKeyById(c, t.Id)
	oldTeam := new(Team)
	if err := datastore.Get(c, k, oldTeam); err == nil {
		if _, err = datastore.Put(c, k, t); err != nil {
			return err
		}
		// use lower trim names as team inverted index store them like this.
		teaminvidmdl.Update(c, oldTeam.KeyName, t.KeyName, t.Id)
	}
	return nil
}

// Get all teams in datastore.
func FindAllTeams(c appengine.Context) []*Team {
	q := datastore.NewQuery("Team")

	var teams []*Team

	if _, err := q.GetAll(c, &teams); err != nil {
		log.Errorf(c, " Team.FindAll, error occurred during GetAll call: %v", err)
	}

	return teams
}

// Get an array of pointers to Teams with respect to an array of ids.
func TeamsByIds(c appengine.Context, ids []int64) []*Team {

	var teams []*Team
	for _, id := range ids {
		if team, err := TeamById(c, id); err == nil {
			teams = append(teams, team)
		} else {
			log.Errorf(c, " Team.ByIds, error occurred during ByIds call: %v", err)
		}
	}
	return teams
}

// Checkes if a user has joined a team or not.
func (t *Team) Joined(c appengine.Context, u *User) bool {

	hasTeam, _ := u.ContainsTeamId(t.Id)
	return hasTeam
}

// Make a user join a team.
func (t *Team) Join(c appengine.Context, u *User) error {
	// add
	if err := u.AddTeamId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Team.Join, error joining tournament for user:%v Error: %v", u.Id, err))
	}

	return nil
}

// make a user leave a team
// Todo: Should we check that the user is indeed a memeber of the team?
func (t *Team) Leave(c appengine.Context, u *User) error {
	if err := u.RemoveTeamId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Team.Leave, error leaving team for user:%v Error: %v", u.Id, err))
	}
	return nil

}

// Check if user is admin of that team.
func IsTeamAdmin(c appengine.Context, teamId int64, userId int64) bool {

	if team, err := TeamById(c, teamId); err == nil {
		return team.AdminId == userId
	} else {
		log.Errorf(c, " Team.IsTeamAdmin, error occurred during ById call: %v", err)
		return false
	}
}

// Given a id, and a word, get the frequency of that word in the team terms.
func GetWordFrequencyForTeam(c appengine.Context, id int64, word string) int64 {

	if teams := FindTeams(c, "Id", id); teams != nil {
		return helpers.CountTerm(strings.Split(teams[0].KeyName, " "), word)
	}
	return 0
}

// from a team id return an array of users/ players that participates in it.
func (t *Team) Players(c appengine.Context) []*User {

	var users []*User

	for _, uId := range t.UserIds {
		user, err := UserById(c, uId)
		if err != nil {
			log.Errorf(c, " Players, cannot find user with ID=%", uId)
		} else {
			users = append(users, user)
		}
	}
	return users
}
