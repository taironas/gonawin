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
	"errors"
	"fmt"
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

func Create(r *http.Request, teamId int64, userId int64) (*TeamRelationship, error) {
	c := appengine.NewContext(r)
	// create new team relationship
	teamRelationshipId, _, err := datastore.AllocateIDs(c, "TeamRelationship", nil, 1)
	if err != nil {
		return nil, err
	}
	
	key := datastore.NewKey(c, "TeamRelationship", "", teamRelationshipId, nil)

	teamRelationship := &TeamRelationship{ teamRelationshipId, teamId, userId, time.Now() }

	_, err = datastore.Put(c, key, teamRelationship)
	if err != nil {
		return nil, err
	}

	return teamRelationship, nil
}

func Destroy(r *http.Request, teamId int64, userId int64) error {
	c := appengine.NewContext(r)
	
	if teamRel := FindByTeamIdAndUserId(r, teamId, userId); teamRel == nil {
		return errors.New(fmt.Sprintf("Cannot find team relationship for teamId=%d and userId=%d", teamId, userId))
	} else {
		key := datastore.NewKey(c, "TeamRelationship", "", teamRel.Id, nil)
			
		return datastore.Delete(c, key)	
	}
}

func FindByTeamIdAndUserId(r *http.Request, teamId int64, userId int64) *TeamRelationship {
	c:= appengine.NewContext(r)
	
	q := datastore.NewQuery("TeamRelationship").Filter("TeamId =", teamId).Filter("UserId =", userId).Limit(1)
	
	var teamRels []*TeamRelationship
	
	if _, err := q.GetAll(appengine.NewContext(r), &teamRels); err == nil && len(teamRels) > 0 {
		return teamRels[0]
	} else {
		c.Errorf("pw: teamrel.FindByTeamIdAndUserId, error occurred during GetAll: %v", err)
	}
	
	return nil
}

func Find(r *http.Request, filter string, value interface{}) []*TeamRelationship{
	c:= appengine.NewContext(r)
	
	q := datastore.NewQuery("TeamRelationship").Filter(filter + " =", value)
	
	var teamRels []*TeamRelationship
	
	if _, err := q.GetAll(c, &teamRels); err != nil {
		c.Errorf("pw: teamrel.Find, error occurred during GetAll: %v", err)
	}
	
	return teamRels
}