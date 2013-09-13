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

package teamrel

import (
	"net/http"
	"time"
	
	"appengine"
	"appengine/datastore"
)

type TeamRelationship struct {
	Id int64
	TeamId int64
	UserId int64
	Created time.Time
}

func Create(r *http.Request, teamId int64, userId int64) *TeamRelationship {
	c := appengine.NewContext(r)
	// create new team relationship
	teamRelationshipId, _, _ := datastore.AllocateIDs(c, "TeamRelationship", nil, 1)
	key := datastore.NewKey(c, "TeamRelationship", "", teamRelationshipId, nil)

	teamRelationship := &TeamRelationship{ teamRelationshipId, teamId, userId, time.Now() }

	_, err := datastore.Put(c, key, teamRelationship)
	if err != nil {
		c.Errorf("Create: %v", err)
	}

	return teamRelationship
}

func Destroy(r *http.Request, teamId int64, userId int64) error {
	return nil
}

func Find(r *http.Request, filter string, value interface{}) *TeamRelationship {
	q := datastore.NewQuery("TeamRelationship").Filter(filter + " =", value).Limit(1)
	
	var teamRels []*TeamRelationship
	
	if _, err := q.GetAll(appengine.NewContext(r), &teamRels); err == nil && len(teamRels) > 0 {
		return teamRels[0]
	}
	
	return nil
}