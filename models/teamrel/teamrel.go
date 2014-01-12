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
	"time"
	
	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
)

type TeamRelationship struct {
	Id int64
	TeamId int64
	UserId int64
	Created time.Time
}

// Create a teamrel entity given a team id and a user id
func Create(c appengine.Context, teamId int64, userId int64) (*TeamRelationship, error) {
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

// Destroy a teamrel relationship given a team id and a user id
func Destroy(c appengine.Context, teamId int64, userId int64) error {
	
	if teamRel := FindByTeamIdAndUserId(c, teamId, userId); teamRel == nil {
		return errors.New(fmt.Sprintf("Cannot find team relationship for teamId=%d and userId=%d", teamId, userId))
	} else {
		key := datastore.NewKey(c, "TeamRelationship", "", teamRel.Id, nil)
			
		return datastore.Delete(c, key)	
	}
}

// look for a relationship given a team id and user id pair
func FindByTeamIdAndUserId(c appengine.Context, teamId int64, userId int64) *TeamRelationship {
	
	q := datastore.NewQuery("TeamRelationship").Filter("TeamId =", teamId).Filter("UserId =", userId).Limit(1)
	
	var teamRels []*TeamRelationship
	
	if _, err := q.GetAll(c, &teamRels); err == nil && len(teamRels) > 0 {
		return teamRels[0]
	} else {
		log.Errorf(c, " teamrel.FindByTeamIdAndUserId, error occurred during GetAll: %v", err)
		return nil
	}
}

// search for teamrels with respect to the filter and value 
func Find(c appengine.Context, filter string, value interface{}) []*TeamRelationship{
	
	q := datastore.NewQuery("TeamRelationship").Filter(filter + " =", value)
	
	var teamRels []*TeamRelationship
	
	if _, err := q.GetAll(c, &teamRels); err != nil {
		log.Errorf(c, " teamrel.Find, error occurred during GetAll: %v", err)
	}
	
	return teamRels
}
