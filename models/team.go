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
)

type Team struct {
	Id            int64
	KeyName       string
	Name          string
	AdminId       int64
	Private       bool
	Created       time.Time
	UserIds       []int64
	TournamentIds []int64
}

type TeamJson struct {
	Id            *int64     `json:",omitempty"`
	KeyName       *string    `json:",omitempty"`
	Name          *string    `json:",omitempty"`
	AdminId       *int64     `json:",omitempty"`
	Private       *bool      `json:",omitempty"`
	Created       *time.Time `json:",omitempty"`
	UserIds       *[]int64   `json:",omitempty"`
	TournamentIds *[]int64   `json:",omitempty"`
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

	team := &Team{teamId, helpers.TrimLower(name), name, adminId, private, time.Now(), emtpyArray, emtpyArray}

	_, err = datastore.Put(c, key, team)
	if err != nil {
		return nil, err
	}
	// udpate inverted index
	AddToTeamInvertedIndex(c, helpers.TrimLower(name), teamId)

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
		UpdateTeamInvertedIndex(c, oldTeam.KeyName, t.KeyName, t.Id)
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
	if err := t.AddUserId(c, u.Id); err != nil {
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
	if err := t.RemoveUserId(c, u.Id); err != nil {
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

func (t *Team) ContainsTournamentId(id int64) (bool, int) {

	for i, tId := range t.TournamentIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}

// Adds a tournament Id in the TournamentId array.
func (t *Team) AddTournamentId(c appengine.Context, tId int64) error {

	if hasTournament, _ := t.ContainsTournamentId(tId); hasTournament {
		return errors.New(fmt.Sprintf("AddTournamentId, allready a member."))
	}
	t.TournamentIds = append(t.TournamentIds, tId)
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// Remove a tournament Id in the TournamentId array.
func (t *Team) RemoveTournamentId(c appengine.Context, tId int64) error {

	if hasTournament, i := t.ContainsTournamentId(tId); !hasTournament {
		return errors.New(fmt.Sprintf("RemoveTournamentId, not a member."))
	} else {
		// as the order of index in tournamentsId is not important,
		// replace elem at index i with last element and resize slice.
		t.TournamentIds[i] = t.TournamentIds[len(t.TournamentIds)-1]
		t.TournamentIds = t.TournamentIds[0 : len(t.TournamentIds)-1]
	}
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// from a team return an array of tournament the user is involved in.
func (t *Team) Tournaments(c appengine.Context) []*Tournament {

	var tournaments []*Tournament

	for _, tId := range t.TournamentIds {
		tournament, err := TournamentById(c, tId)
		if err != nil {
			log.Errorf(c, " Tournaments, cannot find team with ID=%", tId)
		} else {
			tournaments = append(tournaments, tournament)
		}
	}

	return tournaments
}

// Remove a user Id in the UserId array.
func (t *Team) RemoveUserId(c appengine.Context, uId int64) error {

	if hasUser, i := t.ContainsUserId(uId); !hasUser {
		return errors.New(fmt.Sprintf("RemoveUserId, not a member."))
	} else {
		// as the order of index in usersId is not important,
		// replace elem at index i with last element and resize slice.
		t.UserIds[i] = t.UserIds[len(t.UserIds)-1]
		t.UserIds = t.UserIds[0 : len(t.UserIds)-1]
	}
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// Adds a team Id in the UserId array.
func (t *Team) AddUserId(c appengine.Context, uId int64) error {

	if hasUser, _ := t.ContainsUserId(uId); hasUser {
		return errors.New(fmt.Sprintf("AddUserId, allready a member."))
	}

	t.UserIds = append(t.UserIds, uId)
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

func (t *Team) ContainsUserId(id int64) (bool, int) {

	for i, tId := range t.UserIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}
