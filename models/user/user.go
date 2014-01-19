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

package user

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

  teammdl "github.com/santiaago/purple-wing/models/team"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
)

type User struct {
	Id int64
	Email string
	Username string
	Name string
	Auth string
	Created time.Time
}

// creates a user entity, 
func Create(c appengine.Context, email string, username string, name string, auth string) (*User, error) {
	// create new user
	userId, _, err := datastore.AllocateIDs(c, "User", nil, 1)
	if err != nil {
		log.Errorf(c, " User.Create: %v", err)
	}
	
	key := datastore.NewKey(c, "User", "", userId, nil)
	
	user := &User{ userId, email, username, name, auth, time.Now() }

	_, err = datastore.Put(c, key, user)
	if err != nil {
		log.Errorf(c, "User.Create: %v", err)
		return nil, errors.New("model/user: Unable to put user in Datastore")
	}

	return user, nil
}

// search for a user entity given a filter and value
func Find(c appengine.Context, filter string, value interface{}) *User {
	
	q := datastore.NewQuery("User").Filter(filter + " =", value)
	
	var users []*User
	
	if _, err := q.GetAll(c, &users); err == nil && len(users) > 0 {
		return users[0]
	} else {
		log.Errorf(c, " User.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

// find all users present in datastore
func FindAll(c appengine.Context) []*User {
	q := datastore.NewQuery("User")

	var users []*User

	if _, err := q.GetAll(c, &users); err != nil {
		log.Errorf(c, " User.FindAll, error occurred during GetAll call: %v", err)
	}

	return users
}

// find user by id
func ById(c appengine.Context, id int64) (*User, error) {

	var u User
	key := datastore.NewKey(c, "User", "", id, nil)

	if err := datastore.Get(c, key, &u); err != nil {
		log.Errorf(c, " user not found : %v", err)
		return nil, err
	}
	return &u, nil
}

// get key pointer given a user id
func KeyById(c appengine.Context, id int64)(*datastore.Key) {

	key := datastore.NewKey(c, "User", "", id, nil)

	return key
}

// update user given a user pointer
func Update(c appengine.Context, u *User) error {
	k := KeyById(c, u.Id)
	if _, err := datastore.Put(c, k, u); err != nil {
		return err
	}
	return nil
}

// create user from params in datastore and return a pointer to it.
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
	if user = Find(c, queryName, queryValue); user == nil {
		// create user if it does not exist
		if userCreate, err := Create(c, email, username, name, GenerateAuthKey()); err != nil {
			log.Errorf(c, "Signup: %v", err)
			return nil, errors.New("models/user: Unable to create user.")
		} else {
			user = userCreate
		}
	}
	
	return user, nil
}

// generate authentication string key
func GenerateAuthKey() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return ""
	}
	return fmt.Sprintf("%x", b)
}

// from a user id returns an array of teams the user iq involved participates.
func Teams(c appengine.Context, userId int64) []*teammdl.Team {
	
	var teams []*teammdl.Team
	
	teamRels := teamrelmdl.Find(c, "UserId", userId)

	for _, teamRel := range teamRels {
		team, err := teammdl.ById(c, teamRel.TeamId)
		if err != nil {
			log.Errorf(c, " Teams, cannot find team with ID=%", teamRel.TeamId)
		} else {
			teams = append(teams, team)
		}
	}

	return teams
}

// gets an array of teams whos admin id is passed as param
func AdminTeams(c appengine.Context, adminId int64) []*teammdl.Team {
	
	return teammdl.Find(c, "AdminId", adminId)
}
