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

package tournamentrel

import (
	"errors"
	"fmt"
	"time"
	
	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
)

type TournamentRelationship struct {
	Id int64
	TournamentId int64
	UserId int64
	Created time.Time
}

// create a tournamentrel entity given a tournament and a user id pair
func Create(c appengine.Context, tournamentId int64, userId int64) (*TournamentRelationship, error) {
	// create new tournament relationship
	tournamentRelationshipId, _, err := datastore.AllocateIDs(c, "TournamentRelationship", nil, 1)
	if err != nil {
		return nil, err
	}
	
	key := datastore.NewKey(c, "TournamentRelationship", "", tournamentRelationshipId, nil)

	tournamentRelationship := &TournamentRelationship{ tournamentRelationshipId, tournamentId, userId, time.Now() }

	_, err = datastore.Put(c, key, tournamentRelationship)
	if err != nil {
		return nil, err
	}

	return tournamentRelationship, nil
}

// destroy a tournamentrel given a pair tournamentid and userid
func Destroy(c appengine.Context, tournamentId int64, userId int64) error {
	
	if tournamentRel := FindByTournamentIdAndUserId(c, tournamentId, userId); tournamentRel == nil {
		return errors.New(fmt.Sprintf("Cannot find tournament relationship for tournamentId=%d and userId=%d", tournamentId, userId))
	} else {
		key := datastore.NewKey(c, "TournamentRelationship", "", tournamentRel.Id, nil)
			
		return datastore.Delete(c, key)	
	}
}

//  return a pointer to a tournament relationship entity given a pair tournamentid, userid
func FindByTournamentIdAndUserId(c appengine.Context, tournamentId int64, userId int64) *TournamentRelationship {
	
	q := datastore.NewQuery("TournamentRelationship").Filter("TournamentId =", tournamentId).Filter("UserId =", userId).Limit(1)
	
	var tournamentRels []*TournamentRelationship
	
	if _, err := q.GetAll(c, &tournamentRels); err == nil && len(tournamentRels) > 0 {
		return tournamentRels[0]
	} else {
		log.Errorf(c, " tournamentrel.FindByTournamentIdAndUserId, error occurred during GetAll: %v", err)
		return nil
	}
}

// return an array of tournament rels given a filter, value pair
func Find(c appengine.Context, filter string, value interface{}) []*TournamentRelationship {
	
	q := datastore.NewQuery("TournamentRelationship").Filter(filter + " =", value)
	
	var tournamentRels []*TournamentRelationship
	
	if _, err := q.GetAll(c, &tournamentRels); err != nil {
		log.Errorf(c, " tournamentrel.Find, error occurred during GetAll: %v", err)
	}
	
	return tournamentRels
}
