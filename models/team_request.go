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
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers/log"
)

// TeamRequest represents a request to join a team.
//
type TeamRequest struct {
	ID       int64
	TeamID   int64
	TeamName string
	UserID   int64
	UserName string
	Created  time.Time
}

// TeamRequestJSON is JSON representation of the TeamRequest entity.
//
type TeamRequestJSON struct {
	ID       *int64     `json:"Id,omitempty"`
	TeamID   *int64     `json:"TeamId,omitempty"`
	TeamName *string    `json:",omitempty"`
	UserID   *int64     `json:"UserId,omitempty"`
	UserName *string    `json:",omitempty"`
	Created  *time.Time `json:",omitempty"`
}

// CreateTeamRequest creates a teamrequest with params teamid and userid.
//
func CreateTeamRequest(c appengine.Context, teamID int64, teamName string, userID int64, userName string) (*TeamRequest, error) {
	// create new team request
	teamRequestID, _, err := datastore.AllocateIDs(c, "TeamRequest", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "TeamRequest", "", teamRequestID, nil)

	teamRequest := &TeamRequest{teamRequestID, teamID, teamName, userID, userName, time.Now()}

	_, err = datastore.Put(c, key, teamRequest)
	if err != nil {
		return nil, err
	}

	return teamRequest, nil
}

// Destroy a team request given a teamrequestid.
//
func (tr *TeamRequest) Destroy(c appengine.Context) error {

	var teamRequest *TeamRequest
	var err error

	if teamRequest, err = TeamRequestByID(c, tr.ID); err != nil {
		return fmt.Errorf("Cannot find team request with teamRequestId=%d", tr.ID)
	}

	key := datastore.NewKey(c, "TeamRequest", "", teamRequest.ID, nil)

	return datastore.Delete(c, key)
}

// FindTeamRequest searches for all TeamRequest entities with respect of a filter and a value.
//
func FindTeamRequest(c appengine.Context, filter string, value interface{}) []*TeamRequest {

	q := datastore.NewQuery("TeamRequest").Filter(filter+" =", value)

	var teamRequests []*TeamRequest

	if _, err := q.GetAll(c, &teamRequests); err != nil {
		log.Errorf(c, " teamrequest.Find, error occurred during GetAll: %v", err)
	}

	return teamRequests
}

// findByTeamIDAndUserID searches a request by team id and user id pair.
//
func findByTeamIDAndUserID(c appengine.Context, teamID int64, userID int64) *TeamRequest {

	q := datastore.NewQuery("TeamRequest").Filter("TeamId =", teamID).Filter("UserId =", userID).Limit(1)

	var teamRequests []*TeamRequest

	if _, err := q.GetAll(c, &teamRequests); err == nil && len(teamRequests) > 0 {
		return teamRequests[0]
	} else if len(teamRequests) == 0 {
		log.Infof(c, " teamrequest.findByTeamIDAndUserID, no teamRequests found during GetAll")
	} else {
		log.Errorf(c, " teamrequest.findByTeamIDAndUserID, error occurred during GetAll: %v", err)
	}
	return nil
}

// TeamRequestByID returns a teamrequest if it exist given a teamrequestid.
//
func TeamRequestByID(c appengine.Context, id int64) (*TeamRequest, error) {

	var tr TeamRequest
	key := datastore.NewKey(c, "TeamRequest", "", id, nil)

	if err := datastore.Get(c, key, &tr); err != nil {
		log.Errorf(c, " teamrequest.ById, error occurred during Get: %v", err)
		return &tr, err
	}
	return &tr, nil
}

// WasTeamRequestSent checks if for a team id, user id pair, a request was sent.
//
func WasTeamRequestSent(c appengine.Context, teamID int64, userID int64) bool {
	return findByTeamIDAndUserID(c, teamID, userID) != nil
}

// TeamsRequests returns an array of teamRequest entities from an array of teams.
//
func TeamsRequests(c appengine.Context, teams []*Team) []*TeamRequest {
	var teamRequests []*TeamRequest
	for _, team := range teams {
		teamRequests = append(teamRequests, FindTeamRequest(c, "TeamId", team.ID)...)
	}
	return teamRequests
}
