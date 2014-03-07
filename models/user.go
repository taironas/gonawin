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
)

type User struct {
	Id                    int64
	Email                 string
	Username              string
	Name                  string
	IsAdmin               bool    // is user gonawin admin.
	Auth                  string  // authentication auth token
	PredictIds            []int64 // current user predicts.
	ArchivedPredictInds   []int64 // archived user predicts.
	TournamentIds         []int64 // current tournament ids of user <=> tournaments user subscribed.
	ArchivedTournamentIds []int64 // archived tournament ids of user <=> finnished tournametns user subscribed.
	TeamIds               []int64 // current team ids of user <=> teams user belongs to.
	Score                 int64   // overall user score.
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
	Score                 *int64     `json:",omitempty"`
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

	user := &User{userId, email, username, name, isAdmin, auth, emptyArray, emptyArray, emptyArray, emptyArray, emptyArray, int64(0), time.Now()}

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
		// publish new activity
		user.Publish(c, "welcome", "joined gonawin", ActivityEntity{}, ActivityEntity{})
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

// Update an array of users.
func UpdateUsers(c appengine.Context, users []*User) error {
	keys := make([]*datastore.Key, len(users))
	for i, _ := range keys {
		keys[i] = UserKeyById(c, users[i].Id)
	}
	if _, err := datastore.PutMulti(c, keys, users); err != nil {
		return err
	}
	return nil
}

func (u *User) PredictFromMatchId(c appengine.Context, mId int64) (*Predict, error) {
	predicts := PredictsByIds(c, u.PredictIds)
	for i, p := range predicts {
		if p.MatchId == mId {
			return predicts[i], nil
		}
	}
	return nil, nil
}

func (u *User) ScoreForMatch(c appengine.Context, m *Tmatch) (int64, error) {
	desc := "Score for match:"
	var p *Predict
	var err1 error
	if p, err1 = u.PredictFromMatchId(c, m.Id); err1 == nil && p == nil {
		return 0, nil
	} else if err1 != nil {
		log.Errorf(c, "%s unable to get predict for current user %v: %v", desc, u.Id, err1)
		return 0, nil
	}
	return computeScore(m, p), nil
}

// Sort users by score
type UserByScore []*User

func (a UserByScore) Len() int           { return len(a) }
func (a UserByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UserByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }

// Find all activities
func (u *User) Activities(c appengine.Context) []*Activity {
	q := datastore.NewQuery("Activity").Filter("CreatorID=", u.Id).Order("-Published")

	var activities []*Activity

	if _, err := q.GetAll(c, &activities); err != nil {
		log.Errorf(c, "model/activity, FindAll: error occurred during GetAll call: %v", err)
	}

	return activities
}

func (u *User) Publish(c appengine.Context, activityType string, verb string, object ActivityEntity, target ActivityEntity) error {
	var activity Activity
	activity.Type = activityType
	activity.Verb = verb
	activity.Actor = ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
	activity.Object = object
	activity.Target = target
	activity.Published = time.Now()
	activity.CreatorID = u.Id

	return activity.save(c)
}
