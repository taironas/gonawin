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

package teamrequest

import (
	"errors"
	"fmt"
	"time"
	
	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
)

type TeamRequest struct {
	Id int64
	TeamId int64
	UserId int64
	Created time.Time
}

func Create(c appengine.Context, teamId int64, userId int64) (*TeamRequest, error) {
	// create new team request
	teamRequestId, _, err := datastore.AllocateIDs(c, "TeamRequest", nil, 1)
	if err != nil {
		return nil, err
	}
	
	key := datastore.NewKey(c, "TeamRequest", "", teamRequestId, nil)

	teamRequest := &TeamRequest{ teamRequestId, teamId, userId, time.Now() }

	_, err = datastore.Put(c, key, teamRequest)
	if err != nil {
		return nil, err
	}

	return teamRequest, nil
}

func Destroy(c appengine.Context, teamRequestId int64) error {
	
	if teamRequest, err := ById(c, teamRequestId); err != nil {
		return errors.New(fmt.Sprintf("Cannot find team request with teamRequestId=%d", teamRequestId))
	} else {
		key := datastore.NewKey(c, "TeamRequest", "", teamRequest.Id, nil)
			
		return datastore.Delete(c, key)	
	}
}

func Find(c appengine.Context, filter string, value interface{}) []*TeamRequest {
	
	q := datastore.NewQuery("TeamRequest").Filter(filter + " =", value)
	
	var teamRequests []*TeamRequest
	
	if _, err := q.GetAll(c, &teamRequests); err != nil {
		log.Errorf(c, " teamrequest.Find, error occurred during GetAll: %v", err)
	}
	
	return teamRequests
}

func findByTeamIdAndUserId(c appengine.Context, teamId int64, userId int64) *TeamRequest {
	
	q := datastore.NewQuery("TeamRequest").Filter("TeamId =", teamId).Filter("UserId =", userId).Limit(1)
	
	var teamRequests []*TeamRequest
	
	if _, err := q.GetAll(c, &teamRequests); err == nil && len(teamRequests) > 0 {
		return teamRequests[0]
	} else {
		log.Errorf(c, " teamrequest.findByTeamIdAndUserId, error occurred during GetAll: %v", err)
		return nil
	}
}

func ById(c appengine.Context, id int64) (*TeamRequest, error) {
	
	var tr TeamRequest
	key := datastore.NewKey(c, "TeamRequest", "", id, nil)

	if err := datastore.Get(c, key, &tr); err != nil {
		log.Errorf(c, " teamrequest.ById, error occurred during Get: %v", err)
		return &tr, err
	}
	return &tr, nil
}

func Sent(c appengine.Context, teamId int64, userId int64) bool {
	return findByTeamIdAndUserId(c, teamId, userId) != nil
}
