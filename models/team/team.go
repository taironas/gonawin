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
	"net/http"
	"time"
	
	"appengine"
	"appengine/datastore"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
)

type Team struct {
	Id int64
	Name string
	AdminId int64
	Created time.Time
}

func Create(r *http.Request, name string, adminId int64) *Team {
	c := appengine.NewContext(r)
	// create new team
	teamId, _, _ := datastore.AllocateIDs(c, "Team", nil, 1)
	key := datastore.NewKey(c, "Team", "", teamId, nil)

	team := &Team{ teamId, name, adminId, time.Now() }

	_, err := datastore.Put(c, key, team)
	if err != nil {
		c.Errorf("Create: %v", err)
	}

	return team;
}

func Find(r *http.Request, filter string, value interface{}) *Team {
	q := datastore.NewQuery("Team").Filter(filter + " =", value)
	
	var teams []*Team
	
	if _, err := q.GetAll(appengine.NewContext(r), &teams); err == nil && len(teams) > 0 {
		return teams[0]
	}
	
	return nil
}

func FindAll(r *http.Request) []*Team {
	q := datastore.NewQuery("Team")
	
	var teams []*Team
	
	q.GetAll(appengine.NewContext(r), &teams)
	
	return teams
}

func (t *Team) Members() []*usermdl.User {
	var users []*usermdl.User
	
	return users
}