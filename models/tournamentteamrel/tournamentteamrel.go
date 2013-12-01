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

package tournamentteamrel

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	
	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
)

type TournamentTeamRelationship struct {
	Id int64
	TournamentId int64
	TeamId int64
	Created time.Time
}

func Create(r *http.Request, tournamentId int64, teamId int64) (*TournamentTeamRelationship, error) {
	c := appengine.NewContext(r)
	// create new tournament team relationship
	tournamentteamRelationshipId, _, err := datastore.AllocateIDs(c, "TournamentTeamRelationship", nil, 1)
	if err != nil {
		return nil, err
	}
	
	key := datastore.NewKey(c, "TournamentTeamRelationship", "", tournamentteamRelationshipId, nil)

	tournamentTeamRelationship := &TournamentTeamRelationship{ tournamentteamRelationshipId, tournamentId, teamId, time.Now() }

	_, err = datastore.Put(c, key, tournamentTeamRelationship)
	if err != nil {
		return nil, err
	}

	return tournamentTeamRelationship, nil
}

func Destroy(r *http.Request, tournamentId int64, teamId int64) error {
	c := appengine.NewContext(r)
	
	if tournamentteamRel := FindByTournamentIdAndTeamId(r, tournamentId, teamId); tournamentteamRel == nil {
		return errors.New(fmt.Sprintf("Cannot find tournament relationship for tournamentId=%d and teamId=%d", tournamentId, teamId))
	} else {
		key := datastore.NewKey(c, "TournamentTeamRelationship", "", tournamentteamRel.Id, nil)
			
		return datastore.Delete(c, key)	
	}
}

func FindByTournamentIdAndTeamId(r *http.Request, tournamentId int64, teamId int64) *TournamentTeamRelationship {
	c:= appengine.NewContext(r)
	
	q := datastore.NewQuery("TournamentTeamRelationship").Filter("TournamentId =", tournamentId).Filter("TeamId =", teamId).Limit(1)
	
	var tournamentteamRels []*TournamentTeamRelationship
	
	if _, err := q.GetAll(appengine.NewContext(r), &tournamentteamRels); err == nil && len(tournamentteamRels) > 0 {
		return tournamentteamRels[0]
	} else {
		log.Errorf(c, " tournamentteamrel.FindByTournamentIdAndTeamId, error occurred during GetAll: %v", err)
		return nil
	}
}

func Find(r *http.Request, filter string, value interface{}) []*TournamentTeamRelationship {
	c:= appengine.NewContext(r)
	
	q := datastore.NewQuery("TournamentTeamRelationship").Filter(filter + " =", value)
	
	var tournamentteamRels []*TournamentTeamRelationship
	
	if _, err := q.GetAll(c, &tournamentteamRels); err != nil {
		log.Errorf(c, " tournamentteamrel.Find, error occurred during GetAll: %v", err)
	}
	
	return tournamentteamRels
}
