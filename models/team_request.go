/*
 * Copyright (c) 2014 Santiago Arias | Remy Jourde
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

package models

import (
	"errors"
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/gonawin/helpers/log"
)

type TeamRequest struct {
	Id       int64
	TeamId   int64
	TeamName string
	UserId   int64
	UserName string
	Created  time.Time
}

type TeamRequestJson struct {
	Id       *int64     `json:",omitempty"`
	TeamId   *int64     `json:",omitempty"`
	TeamName *string    `json:",omitempty"`
	UserId   *int64     `json:",omitempty"`
	UserName *string    `json:",omitempty"`
	Created  *time.Time `json:",omitempty"`
}

// Create a teamrequest with params teamid and userid
func CreateTeamRequest(c appengine.Context, teamId int64, teamName string, userId int64, userName string) (*TeamRequest, error) {
	// create new team request
	teamRequestId, _, err := datastore.AllocateIDs(c, "TeamRequest", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "TeamRequest", "", teamRequestId, nil)

	teamRequest := &TeamRequest{teamRequestId, teamId, teamName, userId, userName, time.Now()}

	_, err = datastore.Put(c, key, teamRequest)
	if err != nil {
		return nil, err
	}

	return teamRequest, nil
}

// destroy a team request given a teamrequestid
func (tr *TeamRequest) Destroy(c appengine.Context) error {

	if teamRequest, err := TeamRequestById(c, tr.Id); err != nil {
		return errors.New(fmt.Sprintf("Cannot find team request with teamRequestId=%d", tr.Id))
	} else {
		key := datastore.NewKey(c, "TeamRequest", "", teamRequest.Id, nil)

		return datastore.Delete(c, key)
	}
}

// Search for all TeamRequest entities with respect of a filter and a value.
func FindTeamRequest(c appengine.Context, filter string, value interface{}) []*TeamRequest {

	q := datastore.NewQuery("TeamRequest").Filter(filter+" =", value)

	var teamRequests []*TeamRequest

	if _, err := q.GetAll(c, &teamRequests); err != nil {
		log.Errorf(c, " teamrequest.Find, error occurred during GetAll: %v", err)
	}

	return teamRequests
}

// search a request by team id and user id pair
func findByTeamIdAndUserId(c appengine.Context, teamId int64, userId int64) *TeamRequest {

	q := datastore.NewQuery("TeamRequest").Filter("TeamId =", teamId).Filter("UserId =", userId).Limit(1)

	var teamRequests []*TeamRequest

	if _, err := q.GetAll(c, &teamRequests); err == nil && len(teamRequests) > 0 {
		return teamRequests[0]
	} else if len(teamRequests) == 0 {
		log.Infof(c, " teamrequest.findByTeamIdAndUserId, no teamRequests found during GetAll")
	} else {
		log.Errorf(c, " teamrequest.findByTeamIdAndUserId, error occurred during GetAll: %v", err)
	}
	return nil
}

// return a teamrequest if it exist given a teamrequestid
func TeamRequestById(c appengine.Context, id int64) (*TeamRequest, error) {

	var tr TeamRequest
	key := datastore.NewKey(c, "TeamRequest", "", id, nil)

	if err := datastore.Get(c, key, &tr); err != nil {
		log.Errorf(c, " teamrequest.ById, error occurred during Get: %v", err)
		return &tr, err
	}
	return &tr, nil
}

// checks if for a team id, user id pair, a request was sent
func WasTeamRequestSent(c appengine.Context, teamId int64, userId int64) bool {
	return findByTeamIdAndUserId(c, teamId, userId) != nil
}

// Return an array of teamRequest entities from an array of teams.
func TeamsRequests(c appengine.Context, teams []*Team) []*TeamRequest {
	var teamRequests []*TeamRequest
	for _, team := range teams {
		teamRequests = append(teamRequests, FindTeamRequest(c, "TeamId", team.Id)...)
	}
	return teamRequests
}
