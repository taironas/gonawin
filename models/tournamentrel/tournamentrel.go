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
	"net/http"
	"time"
	
	"appengine"
	"appengine/datastore"
)

type TournamentRelationship struct {
	Id int64
	TournamentId int64
	UserId int64
	Created time.Time
}

func Create(r *http.Request, tournamentId int64, userId int64) *TournamentRelationship {
	c := appengine.NewContext(r)
	// create new tournament relationship
	tournamentRelationshipId, _, err := datastore.AllocateIDs(c, "TournamentRelationship", nil, 1)
	if err != nil {
		c.Errorf("pw: TournamentRelationship.Create: %v", err)
	}
	
	key := datastore.NewKey(c, "TournamentRelationship", "", tournamentRelationshipId, nil)

	tournamentRelationship := &TournamentRelationship{ tournamentRelationshipId, tournamentId, userId, time.Now() }

	_, err := datastore.Put(c, key, tournamentRelationship)
	if err != nil {
		c.Errorf("Create: %v", err)
	}

	return tournamentRelationship
}

func Destroy(r *http.Request, tournamentId int64, userId int64) error {
	c := appengine.NewContext(r)
	
	if tournamentRel := FindByTournamentIdAndUserId(r, tournamentId, userId); tournamentRel == nil {
		return errors.New(fmt.Sprintf("Cannot find tournament relationship for tournamentId=%d and userId=%d", tournamentId, userId))
	} else {
		key := datastore.NewKey(c, "TournamentRelationship", "", tournamentRel.Id, nil)
			
		return datastore.Delete(c, key)	
	}
}

func FindByTournamentIdAndUserId(r *http.Request, tournamentId int64, userId int64) *TournamentRelationship {
	q := datastore.NewQuery("TournamentRelationship").Filter("TournamentId =", tournamentId).Filter("UserId =", userId).Limit(1)
	
	var tournamentRels []*TournamentRelationship
	
	if _, err := q.GetAll(appengine.NewContext(r), &tournamentRels); err == nil && len(tournamentRels) > 0 {
		return tournamentRels[0]
	} 
	
	return nil
}

func Find(r *http.Request, filter string, value interface{}) []*TournamentRelationship{
	c:= appengine.NewContext(r)
	
	q := datastore.NewQuery("TournamentRelationship").Filter(filter + " =", value)
	
	var tournamentRels []*TournamentRelationship
	
	if _, err := q.GetAll(c, &tournamentRels); err != nil {
		c.Errorf("pw: error occured in tournamentrel.Find: %v", err)
	}
	
	return tournamentRels
}