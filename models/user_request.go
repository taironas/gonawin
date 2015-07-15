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

type UserRequest struct {
	Id      int64
	TeamId  int64
	UserId  int64
	Created time.Time
}

type UserRequestJson struct {
	Id      *int64     `json:",omitempty"`
	TeamId  *int64     `json:",omitempty"`
	UserId  *int64     `json:",omitempty"`
	Created *time.Time `json:",omitempty"`
}

// Create a user request with params teamid and userid
func CreateUserRequest(c appengine.Context, teamId int64, userId int64) (*UserRequest, error) {
	// create new team request
	id, _, err := datastore.AllocateIDs(c, "UserRequest", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "UserRequest", "", id, nil)

	ur := &UserRequest{id, teamId, userId, time.Now()}

	_, err = datastore.Put(c, key, ur)
	if err != nil {
		return nil, err
	}

	return ur, nil
}

// destroy a user request given a teamrequestid
func (ur *UserRequest) Destroy(c appengine.Context) error {
	key := datastore.NewKey(c, "UserRequest", "", ur.Id, nil)
	return datastore.Delete(c, key)
}

// Search for all TeamRequest entities with respect of a filter and a value.
func FindUserRequests(c appengine.Context, filter string, value interface{}) []*UserRequest {

	q := datastore.NewQuery("UserRequest").Filter(filter+" =", value)

	var userRequests []*UserRequest

	if _, err := q.GetAll(c, &userRequests); err != nil {
		log.Errorf(c, "userRequest.Find, error occurred during GetAll: %v", err)
	}
	return userRequests
}

// search a request by team id and user id pair
func FindUserRequestByTeamAndUser(c appengine.Context, teamId int64, userId int64) *UserRequest {

	q := datastore.NewQuery("UserRequest").Filter("TeamId =", teamId).Filter("UserId =", userId).Limit(1)

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

// return a user request if it exist given a user request id.
func UserRequestById(c appengine.Context, id int64) (*UserRequest, error) {

	var ur UserRequest
	key := datastore.NewKey(c, "UserRequest", "", id, nil)

	if err := datastore.Get(c, key, &ur); err != nil {
		log.Errorf(c, " userrequest.ById, error occurred during Get: %v", err)
		return &ur, err
	}
	return &ur, nil
}
