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

package teams

import (
	"errors"
	"net/http"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// RequestInvite handler, use it to request an invitation to a team.
//  POST	/j/teams/requestinvite/[0-9]+/     Request an invitation to a private team with the given id.
// Response: a JSON formatted status message.
//
func RequestInvite(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Request Invite Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	team, err = extract.Team()
	if err != nil {
		return err
	}

	if mdl.WasTeamRequestSent(c, team.ID, u.ID) {
		return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeTeamRequestAlreadySent)}
	}

	if _, err := mdl.CreateTeamRequest(c, team.ID, team.Name, u.ID, u.Username); err != nil {
		log.Errorf(c, "%s teams.Invite, error when trying to create a team request: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotInvite)}
	}

	// return status message
	return templateshlp.RenderJson(w, c, "team request was created")
}

// SendInvite handler, use it to send an invitation to gonawin.
//	POST	/j/teams/sendinvite/[0-9]+/			Send an invitation to a user with the given team id and user id.
// An activity is published when the invitation is sent.
// Response: a JSON formatted status message.
//
func SendInvite(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Send User Invitation Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	team, err = extract.Team()
	if err != nil {
		return err
	}

	var user *mdl.User
	user, err = extract.User()
	if err != nil {
		return err
	}

	if _, err := mdl.CreateUserRequest(c, team.ID, user.ID); err != nil {
		log.Errorf(c, "%s teams.SendInvite, error when trying to create a user request: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotInvite)}
	}

	// publish new activity
	user.Publish(c, "invitation", "has been invited to join team ", team.Entity(), mdl.ActivityEntity{})

	return templateshlp.RenderJson(w, c, "user request was created")
}

// Invited handler, use it to get all the users who were invited to a team.
//	GET  /j/teams/invited/[0-9]+/			Get list of users invited to a given team id.
// Response: array of JSON formatted users.
//
func Invited(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	desc := "Team Invited Handler:"
	c := appengine.NewContext(r)
	extract := extract.NewContext(c, desc, r)

	var teamID int64
	var err error
	teamID, err = extract.TeamID()
	if err != nil {
		return err
	}

	urs := mdl.FindUserRequests(c, "TeamId", teamID)
	var ids []int64
	for _, ur := range urs {
		ids = append(ids, ur.UserId)
	}

	var users []*mdl.User
	if users, err = mdl.UsersByIds(c, ids); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	ivm := buildTeamInvitedViewModel(users)

	return templateshlp.RenderJson(w, c, ivm)
}

type teamInvitedViewModel struct {
	Users []teamInvitedUserViewModel
}

type teamInvitedUserViewModel struct {
	ID       int64 `json:"Id,omitempty"`
	Username string
	Alias    string
	Score    int64
	ImageURL string
}

func buildTeamInvitedViewModel(users []*mdl.User) teamInvitedViewModel {
	uvm := make([]teamInvitedUserViewModel, len(users))
	for i, u := range users {
		uvm[i].ID = u.ID
		uvm[i].Username = u.Username
		uvm[i].Alias = u.Alias
		uvm[i].Score = u.Score
		uvm[i].ImageURL = helpers.UserImageURL(u.Name, u.ID)
	}

	return teamInvitedViewModel{Users: uvm}
}

// AllowRequest handler, use it to allow a user to join a team.
//	GET  /j/teams/allow/[0-9]+/			Allow a request send by a user on a team.
// After this, the user that send the request will be part of the team.
// Response: a JSON formatted status message.
//
func AllowRequest(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Allow Request Handler:"
	extract := extract.NewContext(c, desc, r)

	var teamRequest *mdl.TeamRequest
	var err error
	teamRequest, err = extract.TeamRequest()
	if err != nil {
		return err
	}

	// join user to the team
	var team *mdl.Team
	team, err = mdl.TeamByID(c, teamRequest.TeamId)
	if err != nil {
		log.Errorf(c, "%s team not found. id: %v, err: %v", desc, teamRequest.TeamId, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
	}
	user, err := mdl.UserByID(c, teamRequest.UserId)
	if err != nil {
		log.Errorf(c, "%s user not found, err: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	team.Join(c, user)
	// request is no more needed so clear it from datastore
	teamRequest.Destroy(c)

	return templateshlp.RenderJson(w, c, "team request was handled")
}

// DenyRequest handler, use it to not allow a user to join a team.
//	GET  /j/teams/deny/[0-9]+/			Deny a request send by a user on a team.
// After this, the user will not be able to be part of the team.
// Response: a JSON formatted status message.
//
func DenyRequest(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Deny Request Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var teamRequest *mdl.TeamRequest

	teamRequest, err = extract.TeamRequest()
	if err != nil {
		return err
	}

	// request is no more needed so clear it from datastore
	teamRequest.Destroy(c)

	return templateshlp.RenderJson(w, c, "team request was handled")
}
