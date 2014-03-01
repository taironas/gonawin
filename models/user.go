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
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
	//teammdl "github.com/santiaago/purple-wing/models/team"
)

type User struct {
	Id                    int64
	Email                 string
	Username              string
	Name                  string
	IsAdmin               bool
	Auth                  string
	PredictIds            []int64 // Current predicts of user.
	ArchivedPredictInds   []int64 // Archived predicts of user.
	TournamentIds         []int64 // Current tournament ids of user.
	ArchivedTournamentIds []int64 // Archived tournament ids of user.
	TeamIds               []int64 // Current team ids of user.
	ArchivedTeamIds       []int64 // Archived team ids of user.
	Created               time.Time
}

type UserJson struct {
	Id                    *int64     `json:",omitempty"`
	Email                 *string    `json:",omitempty"`
	Username              *string    `json:",omitempty"`
	Name                  *string    `json:",omitempty"`
	IsAdmin               *bool      `json:",omitempty"`
	Auth                  *string    `json:",omitempty"`
	PredictIds            *[]int64   `json:",omitempty"`
	ArchivedPredictInds   *[]int64   `json:",omitempty"`
	TournamentIds         *[]int64   `json:",omitempty"`
	ArchivedTournamentIds *[]int64   `json:",omitempty"`
	TeamIds               *[]int64   `json:",omitempty"`
	ArchivedTeamIds       *[]int64   `json:",omitempty"`
	Created               *time.Time `json:",omitempty"`
}

// Creates a user entity.
func CreateUser(c appengine.Context, email string, username string, name string, isAdmin bool, auth string) (*User, error) {
	// create new user
	userId, _, err := datastore.AllocateIDs(c, "User", nil, 1)
	if err != nil {
		log.Errorf(c, " User.Create: %v", err)
	}

	key := datastore.NewKey(c, "User", "", userId, nil)

	emptyArray := make([]int64, 0)

	user := &User{userId, email, username, name, isAdmin, auth, emptyArray, emptyArray, emptyArray, emptyArray, emptyArray, emptyArray, time.Now()}

	_, err = datastore.Put(c, key, user)
	if err != nil {
		log.Errorf(c, "User.Create: %v", err)
		return nil, errors.New("model/user: Unable to put user in Datastore")
	}

	return user, nil
}

// Search for a user entity given a filter and value.
func FindUser(c appengine.Context, filter string, value interface{}) *User {

	q := datastore.NewQuery("User").Filter(filter+" =", value)

	var users []*User

	if _, err := q.GetAll(c, &users); err == nil && len(users) > 0 {
		return users[0]
	} else if len(users) == 0 {
		log.Infof(c, " User.Find, error occurred during GetAll")
	} else {
		log.Errorf(c, " User.Find, error occurred during GetAll: %v", err)
	}
	return nil
}

// Find all users present in datastore.
func FindAllUsers(c appengine.Context) []*User {
	q := datastore.NewQuery("User")

	var users []*User

	if _, err := q.GetAll(c, &users); err != nil {
		log.Errorf(c, "FindAllUser, error occurred during GetAll call: %v", err)
	}

	return users
}

// Find a user entity by id.
func UserById(c appengine.Context, id int64) (*User, error) {

	var u User
	key := datastore.NewKey(c, "User", "", id, nil)

	if err := datastore.Get(c, key, &u); err != nil {
		log.Errorf(c, " user not found : %v", err)
		return nil, err
	}
	return &u, nil
}

// Get key pointer given a user id.
func UserKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "User", "", id, nil)

	return key
}

// Update user given a user pointer.
func (u *User) Update(c appengine.Context) error {
	k := UserKeyById(c, u.Id)
	if _, err := datastore.Put(c, k, u); err != nil {
		return err
	}
	return nil
}

// Create user from params in datastore and return a pointer to it.
func SigninUser(w http.ResponseWriter, r *http.Request, queryName string, email string, username string, name string) (*User, error) {

	c := appengine.NewContext(r)
	var user *User

	queryValue := ""
	if queryName == "Email" {
		queryValue = email
	} else if queryName == "Username" {
		queryValue = username
	} else {
		return nil, errors.New("models/user: no valid query name.")
	}

	// find user
	if user = FindUser(c, queryName, queryValue); user == nil {
		// create user if it does not exist
		isAdmin := false
		if userCreate, err := CreateUser(c, email, username, name, isAdmin, GenerateAuthKey()); err != nil {
			log.Errorf(c, "Signup: %v", err)
			return nil, errors.New("models/user: Unable to create user.")
		} else {
			user = userCreate
		}
	}

	return user, nil
}

// Generate authentication string key.
func GenerateAuthKey() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

// From a user id returns an array of teams the user iq involved participates.
func (u *User) Teams(c appengine.Context) []*Team {

	var teams []*Team

	//teamRels := teamrelmdl.Find(c, "UserId", userId)

	for _, tId := range u.TeamIds {
		t, err := TeamById(c, tId)
		if err != nil {
			log.Errorf(c, "Teams, cannot find team with ID=%", tId)
		} else {
			teams = append(teams, t)
		}
	}

	return teams
}

// Adds a predict Id in the PredictId array.
func (u *User) AddPredictId(c appengine.Context, pId int64) error {

	u.PredictIds = append(u.PredictIds, pId)
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

// Adds a tournament Id in the TournamentId array.
func (u *User) AddTournamentId(c appengine.Context, tId int64) error {

	if hasTournament, _ := u.ContainsTournamentId(tId); hasTournament {
		return errors.New(fmt.Sprintf("AddTournamentId, allready a member."))
	}

	u.TournamentIds = append(u.TournamentIds, tId)
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

// Adds a tournament Id in the TournamentId array.
func (u *User) RemoveTournamentId(c appengine.Context, tId int64) error {

	if hasTournament, i := u.ContainsTournamentId(tId); !hasTournament {
		return errors.New(fmt.Sprintf("RemoveTournamentId, not a member."))
	} else {
		// as the order of index in tournamentsId is not important,
		// replace elem at index i with last element and resize slice.
		u.TournamentIds[i] = u.TournamentIds[len(u.TournamentIds)-1]
		u.TournamentIds = u.TournamentIds[0 : len(u.TournamentIds)-1]
	}
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

func (u *User) ContainsTournamentId(id int64) (bool, int) {

	for i, tId := range u.TournamentIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}

// from a user return an array of tournament the user is involved in.
func (u *User) Tournaments(c appengine.Context) []*Tournament {

	var tournaments []*Tournament

	//tournamentRels := tournamentrelmdl.Find(c, "UserId", userId)

	for _, tId := range u.TournamentIds {
		t, err := TournamentById(c, tId)
		if err != nil {
			log.Errorf(c, " Tournaments, cannot find team with ID=%", tId)
		} else {
			tournaments = append(tournaments, t)
		}
	}

	return tournaments
}

// Adds a team Id in the TeamId array.
func (u *User) AddTeamId(c appengine.Context, tId int64) error {

	if hasTeam, _ := u.ContainsTeamId(tId); hasTeam {
		return errors.New(fmt.Sprintf("AddTeamId, allready a member."))
	}

	u.TeamIds = append(u.TeamIds, tId)
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

// Adds a team Id in the TeamId array.
func (u *User) RemoveTeamId(c appengine.Context, tId int64) error {

	if hasTeam, i := u.ContainsTeamId(tId); !hasTeam {
		return errors.New(fmt.Sprintf("RemoveTeamId, not a member."))
	} else {
		// as the order of index in teamsId is not important,
		// replace elem at index i with last element and resize slice.
		u.TeamIds[i] = u.TeamIds[len(u.TeamIds)-1]
		u.TeamIds = u.TeamIds[0 : len(u.TeamIds)-1]
	}
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

func (u *User) ContainsTeamId(id int64) (bool, int) {

	for i, tId := range u.TeamIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}
