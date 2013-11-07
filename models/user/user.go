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
	"net/http"
	"errors"
	"time"
	
	"appengine"
	"appengine/datastore"

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

func Create(r *http.Request, email string, username string, name string, auth string) (*User, error) {
	c := appengine.NewContext(r)
	// create new user
	userId, _, err := datastore.AllocateIDs(c, "User", nil, 1)
	if err != nil {
		c.Errorf("pw: User.Create: %v", err)
	}
	
	key := datastore.NewKey(c, "User", "", userId, nil)
	
	user := &User{ userId, email, username, name, auth, time.Now() }

	_, err = datastore.Put(c, key, user)
	if err != nil {
		c.Errorf("User.Create: %v", err)
		return nil, errors.New("model/user: Unable to put user in Datastore")
	}

	return user, nil
}

func Find(r *http.Request, filter string, value interface{}) *User{
	c := appengine.NewContext(r)
	
	q := datastore.NewQuery("User").Filter(filter + " =", value)
	
	var users []*User
	
	if _, err := q.GetAll(appengine.NewContext(r), &users); err == nil && len(users) > 0 {
		return users[0]
	} else {
		c.Errorf("pw: User.Find: %v", err)
	}

	return nil
}

func ById(r *http.Request, id int64)(*User, error){
	c := appengine.NewContext(r)

	var u User
	key := datastore.NewKey(c, "User", "", id, nil)

	if err := datastore.Get(c, key, &u); err != nil {
		c.Errorf("pw: user not found : %v", err)
		return &u, err
	}
	return &u, nil
}

func KeyById(r *http.Request, id int64)(*datastore.Key){
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "User", "", id, nil)

	return key
}

func Update(r *http.Request, u *User) error{
	c := appengine.NewContext(r)
	k := KeyById(r, u.Id)
	if _, err := datastore.Put(c, k, u); err != nil {
		return err
	}
	return nil
}

func Teams(r *http.Request, userId int64) []*teammdl.Team {
	c := appengine.NewContext(r)
	
	var teams []*teammdl.Team
	
	teamRels := teamrelmdl.Find(r, "UserId", userId)
	
	for _, teamRel := range teamRels {
		team, err := teammdl.ById(r, teamRel.TeamId)
		
		if err != nil {
			c.Errorf("pw: User.Teams, cannot find team with ID=%", teamRel.TeamId)
		} else {
			teams = append(teams, team)
		}
	}

	return teams
}

func AdminTeams(r *http.Request, adminId int64) []*teammdl.Team {
	
	return teammdl.Find(r, "AdminId", adminId)
}
