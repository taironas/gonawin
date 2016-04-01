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
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers/log"
)

// UserRequest represents the user request entity.
//
type UserRequest struct {
	ID      int64
	TeamID  int64
	UserID  int64
	Created time.Time
}

// UserRequestJSON is JSON representation of the UserRequest structure.
//
type UserRequestJSON struct {
	ID      *int64     `json:"Id,omitempty"`
	TeamID  *int64     `json:"TeamId,omitempty"`
	UserID  *int64     `json:"UserId,omitempty"`
	Created *time.Time `json:",omitempty"`
}

// CreateUserRequest creates a user request with params teamid and userid
//
func CreateUserRequest(c appengine.Context, teamID int64, userID int64) (*UserRequest, error) {
	// create new team request
	id, _, err := datastore.AllocateIDs(c, "UserRequest", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "UserRequest", "", id, nil)

	ur := &UserRequest{id, teamID, userID, time.Now()}

	_, err = datastore.Put(c, key, ur)
	if err != nil {
		return nil, err
	}

	return ur, nil
}

// Destroy a user request given a teamrequestid.
//
func (ur *UserRequest) Destroy(c appengine.Context) error {
	key := datastore.NewKey(c, "UserRequest", "", ur.ID, nil)
	return datastore.Delete(c, key)
}

// FindUserRequests searches for all TeamRequest entities with respect of a filter and a value.
//
func FindUserRequests(c appengine.Context, filter string, value interface{}) []*UserRequest {

	q := datastore.NewQuery("UserRequest").Filter(filter+" =", value)

	var userRequests []*UserRequest

	if _, err := q.GetAll(c, &userRequests); err != nil {
		log.Errorf(c, "userRequest.Find, error occurred during GetAll: %v", err)
	}
	return userRequests
}

// FindUserRequestByTeamAndUser searches a request by team id and user id pair.
//
func FindUserRequestByTeamAndUser(c appengine.Context, teamID int64, userID int64) *UserRequest {

	q := datastore.NewQuery("UserRequest").Filter("TeamId =", teamID).Filter("UserId =", userID).Limit(1)

	var userRequests []*UserRequest

	if _, err := q.GetAll(c, &userRequests); err == nil && len(userRequests) > 0 {
		return userRequests[0]
	} else if len(userRequests) == 0 {
		log.Infof(c, "userrequest.findUserRequestByTeamAndUser, no user request found during GetAll")
	} else {
		log.Errorf(c, "userrequest.findUserRequestByTeamAndUser, error occurred during GetAll: %v", err)
	}
	return nil
}

// UserRequestByID returns a user request if it exist given a user request id.
//
func UserRequestByID(c appengine.Context, id int64) (*UserRequest, error) {

	var ur UserRequest
	key := datastore.NewKey(c, "UserRequest", "", id, nil)

	if err := datastore.Get(c, key, &ur); err != nil {
		log.Errorf(c, " userrequest.ById, error occurred during Get: %v", err)
		return &ur, err
	}
	return &ur, nil
}
