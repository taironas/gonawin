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

// Package users provides the JSON handlers to handle users data in gonawin app.
package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"appengine"
	"appengine/taskqueue"

	"github.com/taironas/route"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"
	mdl "github.com/santiaago/gonawin/models"
)

// use this structure to get information of user in order to update it.
type userData struct {
	User struct {
		Username string
		Name     string
		Alias    string
		Email    string
	}
}

// User index  handler.
func Index(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)

	users := mdl.FindAllUsers(c)

	fieldsToKeep := []string{"Id", "Username", "Name", "Alias", "Email", "Created"}
	usersJson := make([]mdl.UserJson, len(users))
	helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

	return templateshlp.RenderJson(w, c, usersJson)
}

// Show User handler, use it to get the user JSON data.
// including parameter: {teams, tournaments, teamrequests}
// count parameter: default 12
// page parameter: default 1
func Show(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User show handler:"

	rc := requestContext{c, desc, r}
	var user *mdl.User
	var err error

	if user, err = rc.user(); err != nil {
		return err
	}

	log.Infof(c, "User: %v", user)
	log.Infof(c, "User: %v", user.TeamIds)

	fieldsToKeep := []string{"Id", "Username", "Name", "Alias", "Email", "Created", "IsAdmin", "Auth", "TeamIds", "TournamentIds", "Score"}

	var uJson mdl.UserJson
	helpers.InitPointerStructure(user, &uJson, fieldsToKeep)
	log.Infof(c, "%s User: %v", desc, uJson)

	// get with param:
	with := r.FormValue("including")
	params := helpers.SetOfStrings(with)

	var teams []*mdl.Team
	var teamRequests []*mdl.TeamRequest
	var tournaments []*mdl.Tournament
	var invitations []*mdl.Team

	for _, param := range params {
		switch param {
		case "teams":
			// get count parameter, if not present count is set to 20
			strcount := r.FormValue("count")
			count := int64(25)
			if len(strcount) > 0 {
				if n, err := strconv.ParseInt(strcount, 0, 64); err != nil {
					log.Errorf(c, "%s: error during conversion of count parameter: %v", desc, err)
				} else {
					count = n
				}
			}
			// get page parameter, if not present set page to the first one.
			strpage := r.FormValue("page")
			page := int64(1)
			if len(strpage) > 0 {
				if p, err := strconv.ParseInt(strpage, 0, 64); err != nil {
					log.Errorf(c, "%s error during conversion of page parameter: %v", desc, err)
					page = 1
				} else {
					page = p
				}
			}
			teams = user.TeamsByPage(c, count, page)
		case "teamrequests":
			teamRequests = mdl.TeamsRequests(c, teams)
		case "tournaments":
			tournaments = user.Tournaments(c)
		case "invitations":
			invitations = user.Invitations(c)
		}
	}

	// teams
	type team struct {
		Id           int64
		Name         string
		Accuracy     float64
		MembersCount int64
		Private      bool
		ImageURL     string
	}

	ts := make([]team, len(teams))
	for i, t := range teams {
		ts[i].Id = t.Id
		ts[i].Name = t.Name
		ts[i].MembersCount = t.MembersCount
		ts[i].Private = t.Private
		ts[i].ImageURL = helpers.TeamImageURL(t.Name, t.Id)
	}

	// build tournaments stats json
	type TournamentStats struct {
		Id                int64
		Name              string
		ParticipantsCount int
		TeamsCount        int
		Progress          float64
		ImageURL          string
	}

	stats := make([]TournamentStats, len(tournaments))
	for i, t := range tournaments {
		stats[i].Id = t.Id
		stats[i].Name = t.Name
		stats[i].ParticipantsCount = len(t.UserIds)
		stats[i].TeamsCount = len(t.TeamIds)
		stats[i].Progress = t.Progress(c)
		stats[i].ImageURL = helpers.TournamentImageURL(t.Name, t.Id)
	}

	// tournaments
	type tournament struct {
		Id       int64
		Name     string
		UserIds  []int64
		TeamIds  []int64
		ImageURL string
	}

	tournaments2 := make([]tournament, len(tournaments))
	for i, t := range tournaments {
		tournaments2[i].Id = t.Id
		tournaments2[i].UserIds = t.UserIds
		tournaments2[i].TeamIds = t.TeamIds
		tournaments2[i].ImageURL = helpers.TournamentImageURL(t.Name, t.Id)
	}

	// team requests
	teamRequestFieldsToKeep := []string{"Id", "TeamId", "TeamName", "UserId", "UserName"}
	trsJson := make([]mdl.TeamRequestJson, len(teamRequests))
	helpers.TransformFromArrayOfPointers(&teamRequests, &trsJson, teamRequestFieldsToKeep)

	// invitations
	invitationsFieldsToKeep := []string{"Id", "Name"}
	invitationsJson := make([]mdl.TeamJson, len(invitations))
	helpers.TransformFromArrayOfPointers(&invitations, &invitationsJson, invitationsFieldsToKeep)

	// imageURL
	imageURL := helpers.UserImageURL(user.Username, user.Id)

	// data
	data := struct {
		User            mdl.UserJson          `json:",omitempty"`
		Teams           []team                `json:",omitempty"`
		TeamRequests    []mdl.TeamRequestJson `json:",omitempty"`
		Tournaments     []tournament          `json:",omitempty"`
		TournamentStats []TournamentStats     `json:",omitempty"`
		Invitations     []mdl.TeamJson        `json:",omitempty"`
		ImageURL        string                `json:",omitempty"`
	}{
		uJson,
		ts,
		trsJson,
		tournaments2,
		stats,
		invitationsJson,
		imageURL,
	}

	return templateshlp.RenderJson(w, c, data)
}

// User update handler.
func Update(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User update handler:"

	rc := requestContext{c, desc, r}
	var userId int64
	var err error

	if userId, err = rc.userId(); err != nil {
		return err
	}

	if userId != u.Id {
		log.Errorf(c, "%s error user ids do not match. url id:%s user id: %s", desc, userId, u.Id)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
	}

	// only work on name other values should not be editable
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%s Error when reading request body err: %v.", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
	}

	var updatedData userData
	err = json.Unmarshal(body, &updatedData)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body err: %v.", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
	}
	update := false
	if helpers.IsEmailValid(updatedData.User.Email) && updatedData.User.Email != u.Email {
		u.Email = updatedData.User.Email
		update = true
	}

	if helpers.IsStringValid(updatedData.User.Alias) && updatedData.User.Alias != u.Alias {
		u.Alias = updatedData.User.Alias
		update = true
	}

	if update {
		u.Update(c)
	} else {
		data := struct {
			MessageInfo string `json:",omitempty"`
		}{
			"Nothing to update.",
		}
		return templateshlp.RenderJson(w, c, data)
	}

	fieldsToKeep := []string{"Id", "Username", "Name", "Alias", "Email"}
	var uJson mdl.UserJson
	helpers.InitPointerStructure(u, &uJson, fieldsToKeep)

	data := struct {
		MessageInfo string `json:",omitempty"`
		User        mdl.UserJson
	}{
		"User was correctly updated.",
		uJson,
	}
	return templateshlp.RenderJson(w, c, data)
}

// Destroy hander, use this to remove a user from the datastore.
//	POST	/j/user/destroy/[0-9]+/		Destroys the user with the given id.
//
func Destroy(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User Destroy Handler:"

	rc := requestContext{c, desc, r}
	var user *mdl.User
	var err error
	if user, err = rc.user(); err != nil {
		return err
	}

	// delete all team-user relationships
	for _, teamId := range user.TeamIds {
		if mdl.IsTeamAdmin(c, teamId, u.Id) {
			var team *mdl.Team
			if team, err = mdl.TeamById(c, teamId); err != nil {
				log.Errorf(c, "%s team %d not found", desc, teamId)
				return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
			}
			if err = team.RemoveAdmin(c, user.Id); err != nil {
				log.Infof(c, "%s error occurred during admin deletion: %v", desc, err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserIsTeamAdminCannotDelete)}
			}
		} else {
			if err := user.RemoveTeamId(c, teamId); err != nil {
				log.Errorf(c, "%s error when trying to destroy team relationship: %v", desc, err)
			}
		}
	}

	// delete all tournament-user relationships
	for _, tournamentId := range user.TournamentIds {
		if mdl.IsTournamentAdmin(c, tournamentId, u.Id) {
			var tournament *mdl.Tournament
			if tournament, err = mdl.TournamentById(c, tournamentId); err != nil {
				log.Errorf(c, "%s tournament %d not found", desc, tournamentId)
				return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
			}
			if err = tournament.RemoveAdmin(c, user.Id); err != nil {
				log.Infof(c, "%s error occurred during admin deletion: %v", desc, err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserIsTournamentAdminCannotDelete)}
			}
		} else {
			if err := user.RemoveTournamentId(c, tournamentId); err != nil {
				log.Errorf(c, "%s error when trying to destroy tournament relationship: %v", desc, err)
			}
		}
	}

	// send task to delete activities of the user.
	log.Infof(c, "%s Sending to taskqueue: delete activities", desc)

	activityIds, err1 := json.Marshal(u.ActivityIds)
	if err1 != nil {
		log.Errorf(c, "%s Error marshaling", desc, err1)
	}

	task := taskqueue.NewPOSTTask("/a/publish/users/deleteactivities/", url.Values{
		"activity_ids": []string{string(activityIds)},
	})

	if _, err := taskqueue.Add(c, task, ""); err != nil {
		log.Errorf(c, "%s unable to add task to taskqueue.", desc)
	} else {
		log.Infof(c, "%s add task to taskqueue successfully", desc)
	}

	// send task to delete predicts of the user.
	log.Infof(c, "%s Sending to taskqueue: delete predicts", desc)

	predictsIds, err1 := json.Marshal(u.PredictIds)
	if err1 != nil {
		log.Errorf(c, "%s Error marshaling", desc, err1)
	}

	task = taskqueue.NewPOSTTask("/a/publish/users/deletepredicts/", url.Values{
		"predict_ids": []string{string(predictsIds)},
	})

	if _, err := taskqueue.Add(c, task, ""); err != nil {
		log.Errorf(c, "%s unable to add task to taskqueue.", desc)
	} else {
		log.Infof(c, "%s add task to taskqueue successfully", desc)
	}

	// delete the user
	user.Destroy(c)

	// return destroyed status
	msg := fmt.Sprintf("The user %s was correctly deleted!", user.Username)
	data := struct {
		MessageInfo string `json:",omitempty"`
	}{
		msg,
	}

	return templateshlp.RenderJson(w, c, data)
}

// Teams handler, use this to retreive the JSON data of the user teams.
// count parameter: default 12
// page parameter: default 1
func Teams(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User joined teams handler:"

	rc := requestContext{c, desc, r}

	var user *mdl.User
	var err error
	if user, err = rc.user(); err != nil {
		return err
	}

	// get with param:
	var teams []*mdl.Team
	// get count parameter, if not present count is set to 20
	strcount := r.FormValue("count")
	count := int64(25)
	if len(strcount) > 0 {
		if n, err := strconv.ParseInt(strcount, 0, 64); err != nil {
			log.Errorf(c, "%s: error during conversion of count parameter: %v", desc, err)
		} else {
			count = n
		}
	}
	// get page parameter, if not present set page to the first one.
	strpage := r.FormValue("page")
	page := int64(1)
	if len(strpage) > 0 {
		if p, err := strconv.ParseInt(strpage, 0, 64); err != nil {
			log.Errorf(c, "%s error during conversion of page parameter: %v", desc, err)
			page = 1
		} else {
			page = p
		}
	}
	teams = user.TeamsByPage(c, count, page)

	// teams
	teamsFieldsToKeep := []string{"Id", "Name"}
	teamsJson := make([]mdl.TeamJson, len(teams))
	helpers.TransformFromArrayOfPointers(&teams, &teamsJson, teamsFieldsToKeep)

	// data
	data := struct {
		Teams []mdl.TeamJson `json:",omitempty"`
	}{
		teamsJson,
	}
	return templateshlp.RenderJson(w, c, data)
}

// Tournaments user handler, use this to retrieve the JSON data of the tournaments of the user.
// count parameter: default 25
// page parameter: default 1
func Tournaments(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User joined teams handler:"

	rc := requestContext{c, desc, r}

	var user *mdl.User
	var err error
	if user, err = rc.user(); err != nil {
		return err
	}

	// get with param:
	var tournaments []*mdl.Tournament
	// get count parameter, if not present count is set to 25
	strcount := r.FormValue("count")
	count := int64(25)
	if len(strcount) > 0 {
		if n, err := strconv.ParseInt(strcount, 0, 64); err != nil {
			log.Errorf(c, "%s: error during conversion of count parameter: %v", desc, err)
		} else {
			count = n
		}
	}

	// get page parameter, if not present set page to the first one.
	strpage := r.FormValue("page")
	page := int64(1)
	if len(strpage) > 0 {
		if p, err := strconv.ParseInt(strpage, 0, 64); err != nil {
			log.Errorf(c, "%s error during conversion of page parameter: %v", desc, err)
			page = 1
		} else {
			page = p
		}
	}
	tournaments = user.TournamentsByPage(c, count, page)

	// tournaments
	tournamentsFieldsToKeep := []string{"Id", "Name"}
	tournamentsJson := make([]mdl.TournamentJson, len(tournaments))
	helpers.TransformFromArrayOfPointers(&tournaments, &tournamentsJson, tournamentsFieldsToKeep)

	// data
	data := struct {
		Tournaments []mdl.TournamentJson `json:",omitempty"`
	}{
		tournamentsJson,
	}
	return templateshlp.RenderJson(w, c, data)
}

// AllowInvitation user handler, use it to allow an invitation to a team.
func AllowInvitation(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User allow invitation handler:"

	rc := requestContext{c, desc, r}

	var team *mdl.Team
	var err error
	if team, err = rc.team(); err != nil {
		return err
	}

	ur := mdl.FindUserRequestByTeamAndUser(c, team.Id, u.Id)
	if ur == nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// add user as member of team.
	if err = team.Join(c, u); err != nil {
		log.Errorf(c, "Team Join Handler: error on Join team: %v", err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// destroy user request.
	if err = ur.Destroy(c); err != nil {
		log.Errorf(c, "%s error when destroying user request. Error: %v", err)
	}

	var tJson mdl.TeamJson
	fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

	// publish new activity
	var updatedUser *mdl.User
	if updatedUser, err = mdl.UserById(c, u.Id); err != nil {
		log.Errorf(c, "User not found %v", u.Id)
	}
	updatedUser.Publish(c, "team", "joined team", team.Entity(), mdl.ActivityEntity{})

	msg := fmt.Sprintf("You accepted invitation to team %s.", team.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Team        mdl.TeamJson
	}{
		msg,
		tJson,
	}
	return templateshlp.RenderJson(w, c, data)
}

func DenyInvitation(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "User deny invitation handler:"

	if r.Method == "POST" {
		strTeamId, err := route.Context.Get(r, "teamId")
		if err != nil {
			log.Errorf(c, "%s error getting team id, err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var teamId int64
		teamId, err = strconv.ParseInt(strTeamId, 0, 64)
		if err != nil {
			log.Errorf(c, "%s error converting team id from string to int64, err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		// check team
		var team *mdl.Team
		team, err = mdl.TeamById(c, teamId)
		if err != nil {
			log.Errorf(c, "%s team not found", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		ur := mdl.FindUserRequestByTeamAndUser(c, teamId, u.Id)
		if ur != nil {
			// destroy user request.
			if errd := ur.Destroy(c); errd != nil {
				log.Errorf(c, "%s error when destroying user request. Error: %v", errd)
			}
			msg := fmt.Sprintf("You denied an invitation to team %s.", team.Name)
			data := struct {
				MessageInfo string `json:",omitempty"`
			}{
				msg,
			}
			return templateshlp.RenderJson(w, c, data)
		}
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
