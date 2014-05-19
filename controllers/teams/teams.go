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

// Package teams provides the JSON handlers to handle teams data in gonawin app.
//
// It provides the following methods
//
//	GET	/j/teams/				Retreives all teams.
//	POST	/j/teams/new/				Creates a new team.
//	GET	/j/teams/show/[0-9]+/			Retreives the team with the given id.
//	POST	/j/teams/update/[0-9]+/			Updates the team with the given id.
//	POST	/j/teams/destroy/[0-9]+/		Destroys the team with the given id.
//	POST	/j/teams/invite/[0-9]+/			Request an invitation to a private team with the given id.
//	POST	/j/teams/allow/[0-9]+/			Allow a user to be a member of a team with the given id.
//	POST	/j/teams/deny/[0-9]+/			Deny entrance of user to be a member of a team with the given id.
//	POST	/j/teams/join/[0-9]+/			Make a user join a team with the given id.
//	GET	/j/teams/search/			Search for all teams respecting the query "q"
//	GET	/j/teams/[0-9]+/members/		Retreives all members of a team with the given id.
//	GET	/j/teams/[0-9]+/ranking/		Retreives the ranking of a team with the given id.
//	GET	/j/teams/[0-9]+/accuracies/		Retreives all the tournament accuracies of a team with the given id.
//	GET	/j/teams/[0-9]+/accuracies/[0-9]+/	Retreives accuracies of a team with the given id for the specified tournament.
//
//
// Every method below gives more information about every API call, its parameters, and its resutls.
package teams

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"appengine"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/handlers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

type TeamData struct {
	Name        string
	Description string
	Visibility  string
}

type PriceData struct {
	Description string
}

// team Index handler.
//      GET     /j/teams/?			List users not joined by user.
// Parameters:
//   'page' a int indicating the page number.
//   'count' a int indicating the number of teams per page number. default value is 25
// Reponse: array of JSON formatted teams.
func Index(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "teams index handler: "

	if r.Method == "GET" {
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
		// fetch teams
		teams := mdl.GetNotJoinedTeams(c, u, count, page)

		if len(teams) == 0 {
			return templateshlp.RenderEmptyJsonArray(w, c)
		}
		teamsJson := make([]mdl.TeamJson, len(teams))
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private", "Accuracy", "MembersCount"}
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, teamsJson)

	} else {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
}

// team new handler.
//	POST	/j/teams/new/				Creates a new team.
//
func New(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	desc := "Team New Handler:"

	c := appengine.NewContext(r)

	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "%s Error when decoding request body: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
		}

		var data TeamData
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Errorf(c, "%s Error when decoding request body: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
		}

		if len(data.Name) <= 0 {
			log.Errorf(c, "%s 'Name' field cannot be empty", desc)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeNameCannotBeEmpty)}
		} else if t := mdl.FindTeams(c, "KeyName", helpers.TrimLower(data.Name)); t != nil {
			log.Errorf(c, "%s That team name already exists.", desc)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamAlreadyExists)}
		} else {
			team, err := mdl.CreateTeam(c, data.Name, data.Description, u.Id, data.Visibility == "Private")
			if err != nil {
				log.Errorf(c, "%s error when trying to create a team: %v", desc, err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
			}
			// join the team
			if err = team.Join(c, u); err != nil {
				log.Errorf(c, "%s error when trying to create a team relationship: %v", desc, err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
			}
			// publish new activity
			u.Publish(c, "team", "created a new team", team.Entity(), mdl.ActivityEntity{})

			// return the newly created team
			var tJson mdl.TeamJson
			fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
			helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

			msg := fmt.Sprintf("The team %s was correctly created!", team.Name)
			data := struct {
				MessageInfo string `json:",omitempty"`
				Team        mdl.TeamJson
			}{
				msg,
				tJson,
			}

			return templateshlp.RenderJson(w, c, data)
		}
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// team show handler
//	GET	/j/teams/show/[0-9]+/			Retreives the team with the given id.
//
func Show(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	desc := "Team Show Handler:"
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var team *mdl.Team
		if team, err = mdl.TeamById(c, intID); err != nil {
			log.Errorf(c, "%s team not found: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		// get data for json team

		// build team json
		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "Description", "AdminIds", "Private", "TournamentIds", "Accuracy"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		// build players json
		players := team.Players(c)
		fieldsToKeepForPlayer := []string{"Id", "Username", "Alias", "Score"}
		playersJson := make([]mdl.UserJson, len(players))
		helpers.TransformFromArrayOfPointers(&players, &playersJson, fieldsToKeepForPlayer)

		// build tournaments json
		tournaments := team.Tournaments(c)
		type tournament struct {
			Id                int64  `json:",omitempty"`
			Name              string `json:",omitempty"`
			ParticipantsCount int
			TeamsCount        int
			Progress          float64
		}
		ts := make([]tournament, len(tournaments))
		for i, t := range tournaments {
			ts[i].Id = t.Id
			ts[i].Name = t.Name
			ts[i].ParticipantsCount = len(t.UserIds)
			ts[i].TeamsCount = len(t.TeamIds)
			ts[i].Progress = t.Progress(c)
		}

		teamData := struct {
			Team        mdl.TeamJson
			Joined      bool
			RequestSent bool
			Players     []mdl.UserJson
			Tournaments []tournament
		}{
			tJson,
			team.Joined(c, u),
			mdl.WasTeamRequestSent(c, intID, u.Id),
			playersJson,
			ts,
		}
		return templateshlp.RenderJson(w, c, teamData)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// team update handler
//	POST	/j/teams/update/[0-9]+/			Updates the team with the given id.
//
func Update(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	desc := "Team Update Handler:"
	c := appengine.NewContext(r)

	if r.Method == "POST" {

		teamID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotUpdate)}
		}

		if !mdl.IsTeamAdmin(c, teamID, u.Id) {
			log.Errorf(c, "%s user is not admin", desc)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamUpdateForbiden)}
		}

		var team *mdl.Team
		team, err = mdl.TeamById(c, teamID)
		if err != nil {
			log.Errorf(c, "%s team not found. id: %v, err: %v", desc, teamID, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotUpdate)}
		}
		// only work on name and private. Other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "%s Error when reading request body err: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}

		var updatedData TeamData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			log.Errorf(c, "%s Error when decoding request body err: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}

		updatedPrivate := updatedData.Visibility == "private"
		log.Errorf(c, "%s %v. %v", desc, team.Name, team.Private)
		log.Errorf(c, "%s %v. %v", desc, updatedData.Name, updatedPrivate)
		log.Errorf(c, "%s visibility %v.", desc, updatedData.Visibility)
		log.Errorf(c, "%s updateddata %v.", desc, updatedData)

		if helpers.IsStringValid(updatedData.Name) &&
			helpers.IsStringValid(updatedData.Description) &&
			(updatedData.Name != team.Name || updatedData.Description != team.Description || updatedPrivate != team.Private) {
			if updatedData.Name != team.Name {
				// be sure that a team with that name does not exist in datastore.
				if t := mdl.FindTeams(c, "KeyName", helpers.TrimLower(updatedData.Name)); t != nil {
					log.Errorf(c, "%s That team name already exists.", desc)
					return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamAlreadyExists)}
				}
				// update data
				team.Name = updatedData.Name
			}
			team.Description = updatedData.Description
			team.Private = updatedPrivate
			team.Update(c)
		} else {
			log.Errorf(c, "%s Cannot update because updated is not valid.", desc)
			log.Errorf(c, "%s Update name = %s", desc, updatedData.Name)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}

		// publish new activity
		u.Publish(c, "team", "updated team", team.Entity(), mdl.ActivityEntity{})

		// keep only needed fields for json api
		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		msg := fmt.Sprintf("The team %s was correctly updated!", team.Name)
		data := struct {
			MessageInfo string `json:",omitempty"`
			Team        mdl.TeamJson
		}{
			msg,
			tJson,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Team destroy handler
//	POST	/j/teams/destroy/[0-9]+/		Destroys the team with the given id.
//
func Destroy(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Destroy Handler:"
	if r.Method == "POST" {

		teamID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotDelete)}
		}

		if !mdl.IsTeamAdmin(c, teamID, u.Id) {
			log.Errorf(c, "%s user is not admin", desc)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamDeleteForbiden)}
		}
		var team *mdl.Team
		team, err = mdl.TeamById(c, teamID)
		if err != nil {
			log.Errorf(c, "%s team not found. id: %v, err: %v", desc, teamID, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotUpdate)}
		}

		// delete all team-user relationships
		for _, player := range team.Players(c) {
			if err := player.RemoveTeamId(c, team.Id); err != nil {
				log.Errorf(c, "%s error when trying to destroy team relationship: %v", desc, err)
			}
		}
		// delete all tournament-team relationships
		for _, tournament := range team.Tournaments(c) {
			if err := tournament.RemoveTeamId(c, team.Id); err != nil {
				log.Errorf(c, "%serror when trying to destroy tournament relationship: %v", desc, err)
			}
		}
		// delete the team
		team.Destroy(c)

		// publish new activity
		u.Publish(c, "team", "deleted team", team.Entity(), mdl.ActivityEntity{})

		msg := fmt.Sprintf("The team %s was correctly deleted!", team.Name)
		data := struct {
			MessageInfo string `json:",omitempty"`
		}{
			msg,
		}

		// return destroyed status
		return templateshlp.RenderJson(w, c, data)
	} else {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
}

// Request Invite handler.
// A user sends a team a request to join a team if team is private.
// use this handler when you wish to request an invitation to a team.
// this is done when the team in set as 'private' and the user wishes to join it.
func RequestInvite(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	desc := "Team Request Invite Handler:"
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotInvite)}
		}
		// check if team id exist.

		if _, err := mdl.CreateTeamRequest(c, intID, u.Id); err != nil {
			log.Errorf(c, "%s teams.Invite, error when trying to create a team request: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotInvite)}
		}
		// return destroyed status
		return templateshlp.RenderJson(w, c, "team request was created")
	}
	return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Send Invite handler.
// A team sends a user an invitation with tuple (user, team) to join a team.
func SendInvite(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	desc := "Team Send User Invitation Handler:"
	c := appengine.NewContext(r)

	if r.Method == "POST" {

		teamId, err1 := handlers.PermalinkID(r, c, 4)
		if err1 != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotInvite)}
		}
		userId, err2 := handlers.PermalinkID(r, c, 5)
		if err2 != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err2)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotInvite)}
		}

		// check that ids exist in datastore.
		team, err := mdl.TeamById(c, teamId)
		if err != nil {
			log.Errorf(c, "%s team not found. id: %v, err: %v", desc, teamId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotUpdate)}
		}

		user, err := mdl.UserById(c, userId)
		if err != nil {
			log.Errorf(c, "%s team not found. id: %v, err: %v", desc, userId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotUpdate)}
		}

		if _, err := mdl.CreateUserRequest(c, teamId, userId); err != nil {
			log.Errorf(c, "%s teams.SendInvite, error when trying to create a user request: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotInvite)}
		}

		// publish new activity
		user.Publish(c, "invitation", "has been invited to join team ", team.Entity(), mdl.ActivityEntity{})

		return templateshlp.RenderJson(w, c, "user request was created")
	}
	return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Invited handler.
// List all sent invitations of team.
func Invited(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	desc := "Team Invited Handler:"
	c := appengine.NewContext(r)

	if r.Method == "GET" {

		teamId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotInvite)}
		}

		urs := mdl.FindUserRequests(c, "TeamId", teamId)
		var ids []int64
		for _, ur := range urs {
			ids = append(ids, ur.UserId)
		}

		users := mdl.UsersByIds(c, ids)

		// filter team information to return in json api
		fieldsToKeep := []string{"Id", "Username", "Alias"}
		usersJson := make([]mdl.UserJson, len(users))
		helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

		data := struct {
			Users []mdl.UserJson
		}{
			usersJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Allow handler.
// use this handler to allow a request send by a user on a team.
// after this, the user that send the request will be part of the team
func AllowRequest(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Allow Request Handler:"
	if r.Method == "POST" {
		requestId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s teams.AllowRequest, id could not be extracter from url: %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
		}

		if teamRequest, err := mdl.TeamRequestById(c, requestId); err == nil {
			// join user to the team
			var team *mdl.Team
			team, err = mdl.TeamById(c, teamRequest.TeamId)
			if err != nil {
				log.Errorf(c, "%s team not found. id: %v, err: %v", desc, teamRequest.TeamId, err)
				return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
			}
			user, err := mdl.UserById(c, teamRequest.UserId)
			if err != nil {
				log.Errorf(c, "%s user not found, err: %v", desc, err)
				return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
			}

			team.Join(c, user)
			// request is no more needed so clear it from datastore
			teamRequest.Destroy(c)

		} else {
			log.Errorf(c, "%s cannot find team request with id=%d", desc, requestId)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
		}

		return templateshlp.RenderJson(w, c, "team request was handled")

	} else {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
}

// Deny handler.
// use this handler to deny a request send by a user on a team.
// the user will not be able to be part of the team
func DenyRequest(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Deny Request Handler:"
	if r.Method == "POST" {
		requestId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s teams.DenyRequest, id could not be extracter from url: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
		}
		if teamRequest, err := mdl.TeamRequestById(c, requestId); err != nil {
			log.Errorf(c, "%s teams.DenyRequest, team request not found: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
		} else {
			// request is no more needed so clear it from datastore
			teamRequest.Destroy(c)
		}

		return templateshlp.RenderJson(w, c, "team request was handled")

	} else {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
}

// Team search handler.
// Use this handler to search for a team.
//	GET	/j/teams/search/			Search for all teams respecting the query "q"
//
func Search(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Search Handler:"
	keywords := r.FormValue("q")
	if r.Method == "GET" && (len(keywords) > 0) {

		words := helpers.SetOfStrings(keywords)
		ids, err := mdl.GetTeamInvertedIndexes(c, words)
		if err != nil {
			log.Errorf(c, "%s teams.Index, error occurred when getting indexes of words: %v", desc, err)
			data := struct {
				MessageDanger string `json:",omitempty"`
			}{
				"Oops! something went wrong, we are unable to perform search query.",
			}
			return templateshlp.RenderJson(w, c, data)
		}
		result := mdl.TeamScore(c, keywords, ids)
		log.Infof(c, "%s result from TeamScore: %v", desc, result)
		teams := mdl.TeamsByIds(c, result)
		log.Infof(c, "%s ByIds result %v", desc, teams)
		if len(teams) == 0 {
			msg := fmt.Sprintf("Oops! Your search - %s - did not match any %s.", keywords, "team")
			data := struct {
				MessageInfo string `json:",omitempty"`
			}{
				msg,
			}

			return templateshlp.RenderJson(w, c, data)
		}
		// filter team information to return in json api
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private", "Accuracy", "MembersCount"}
		teamsJson := make([]mdl.TeamJson, len(teams))
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)
		// we should not directly return an array. so we add an extra layer.
		data := struct {
			Teams []mdl.TeamJson `json:",omitempty"`
		}{
			teamsJson,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Team members handler.
//	/j/teams/[0-9]+/members/	GET			use this handler to get members of a team.
//
func Members(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Members Handler:"
	if r.Method == "GET" {
		teamId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamMemberNotFound)}
		}
		team, err1 := mdl.TeamById(c, teamId)
		if err1 != nil {
			log.Errorf(c, "%s team not found. id: %v, err: %v", desc, teamId, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamMemberNotFound)}
		}

		// build members json
		members := team.Players(c)
		fieldsToKeepForMember := []string{"Id", "Username", "Alias", "Score"}
		membersJson := make([]mdl.UserJson, len(members))
		helpers.TransformFromArrayOfPointers(&members, &membersJson, fieldsToKeepForMember)

		data := struct {
			Members []mdl.UserJson
		}{
			membersJson,
		}
		return templateshlp.RenderJson(w, c, data)
	} else {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
}

// team Prices handler.
// use this handler to get the prices of a team.
func Prices(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Prices Handler:"

	if r.Method == "GET" {
		teamId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var t *mdl.Team
		t, err = mdl.TeamById(c, teamId)
		if err != nil {
			log.Errorf(c, "%s team with id:%v was not found %v", desc, teamId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		log.Infof(c, "%s ready to build a price array", desc)
		prices := t.Prices(c)

		data := struct {
			Prices []*mdl.Price
		}{
			prices,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Team prices by tournament handler:
//
// Use this handler to get the price of a team for a specific tournament.
//	GET	/j/teams/[0-9]+/prices/[0-9]+/	Retreives price of a team with the given id for the specified tournament.
//
func PriceByTournament(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Prices by tournament Handler:"

	if r.Method == "GET" {
		teamId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var t *mdl.Team
		t, err = mdl.TeamById(c, teamId)
		if err != nil {
			log.Errorf(c, "%s team with id:%v was not found %v", desc, teamId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		tournamentId, err := handlers.PermalinkID(r, c, 5)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		log.Infof(c, "%s ready to get price", desc)
		p := t.PriceByTournament(c, tournamentId)

		data := struct {
			Price *mdl.Price
		}{
			p,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}

}

// Team update price by tournament handler:
//
// Use this handler to get the price of a team for a specific tournament.
//	GET	/j/teams/[0-9]+/prices/update/[0-9]+/	Update Retreives price of a team with the given id for the specified tournament.
//
func UpdatePrice(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team update price Handler:"

	if r.Method == "POST" {
		teamId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var t *mdl.Team
		t, err = mdl.TeamById(c, teamId)
		if err != nil {
			log.Errorf(c, "%s team with id:%v was not found %v", desc, teamId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		tournamentId, err := handlers.PermalinkID(r, c, 6)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		log.Infof(c, "%s ready to get price", desc)
		p := t.PriceByTournament(c, tournamentId)

		// only work on name and private. Other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "%s Error when reading request body err: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}

		var priceData PriceData
		err = json.Unmarshal(body, &priceData)
		if err != nil {
			log.Errorf(c, "%s Error when decoding request body err: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}

		if helpers.IsStringValid(priceData.Description) && (p.Description != priceData.Description) {
			// update data
			p.Description = priceData.Description
			p.Update(c)
		} else {
			log.Errorf(c, "%s Cannot update because updated is not valid.", desc)
			log.Errorf(c, "%s Update description = %s", desc, priceData.Description)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}

		data := struct {
			Price *mdl.Price
		}{
			p,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}

}
