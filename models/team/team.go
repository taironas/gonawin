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
	"net/http"
	"time"
	
	"appengine"
	"appengine/datastore"
	
	"github.com/santiaago/purple-wing/helpers"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
	usermdl "github.com/santiaago/purple-wing/models/user"
	searchmdl "github.com/santiaago/purple-wing/models/search"
)

type Team struct {
	Id int64
	KeyName string
	Name string
	AdminId int64
	Private bool
	Created time.Time
}

func Create(r *http.Request, name string, adminId int64, private bool) *Team {
	c := appengine.NewContext(r)
	// create new team
	teamId, _, _ := datastore.AllocateIDs(c, "Team", nil, 1)
	key := datastore.NewKey(c, "Team", "", teamId, nil)

	team := &Team{ teamId, helpers.TrimLower(name), name, adminId, private, time.Now() }

	_, err := datastore.Put(c, key, team)
	if err != nil {
		c.Errorf("Create: %v", err)
	}
	searchmdl.AddToTeamInvertedIndex(r, helpers.TrimLower(name), teamId)

	return team
}

func Find(r *http.Request, filter string, value interface{}) *Team {
	q := datastore.NewQuery("Team").Filter(filter + " =", value).Limit(1)
	
	var teams []*Team
	
	if _, err := q.GetAll(appengine.NewContext(r), &teams); err == nil && len(teams) > 0 {
		return teams[0]
	}
	
	return nil
}

func ById(r *http.Request, id int64) (*Team, error) {
	c := appengine.NewContext(r)

	var t Team
	key := datastore.NewKey(c, "Team", "", id, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		c.Errorf("pw: team not found : %v", err)
		return &t, err
	}
	return &t, nil
}


func KeyById(r *http.Request, id int64) (*datastore.Key) {
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "Team", "", id, nil)

	return key
}


func Update(r *http.Request, id int64, t *Team) error {
	c := appengine.NewContext(r)
	// TODO: get old name before updating team
	k := KeyById(r, id)
	if _, err := datastore.Put(c, k, t); err != nil {
		return err
	}
	//searchmdl.UpdateToTeamInvertedIndex(r, oldname, newname, id)
	return nil
}

func FindAll(r *http.Request) []*Team {
	q := datastore.NewQuery("Team")
	
	var teams []*Team
	
	q.GetAll(appengine.NewContext(r), &teams)
	
	return teams
}

func Joined(r *http.Request, teamId int64, userId int64) bool {
	teamRel := teamrelmdl.FindByTeamIdAndUserId(r, teamId, userId)
	return teamRel != nil
}

func Join(r *http.Request, teamId int64, userId int64) error {
	if teamRel := teamrelmdl.Create(r, teamId, userId); teamRel == nil {
		return errors.New("error during team relationship creation")
	}

	return nil
}

func Leave(r *http.Request, teamId int64, userId int64) error {
	return teamrelmdl.Destroy(r, teamId, userId)
}

func IsTeamAdmin(r *http.Request, teamId int64, userId int64) bool {
	if team, err := ById(r, teamId); err == nil {
		return team.AdminId == userId
	}
	
	return false
}

func Players(r *http.Request, teamId int64) []*usermdl.User {
	
	var users []*usermdl.User
	
	teamRels := teamrelmdl.FindByTeamId(r, teamId)
	
	for _, teamRel := range teamRels {
		user, _ := usermdl.ById(r, teamRel.UserId)
		
		users = append(users, user)
	}

	return users
}
