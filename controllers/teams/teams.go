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

package teams

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	mdl "github.com/santiaago/purple-wing/models"
)

type TeamData struct {
	Name       string
	Visibility string
}

// json index handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		teams := mdl.FindAllTeams(c)
		if len(teams) == 0 {
			return templateshlp.RenderEmptyJsonArray(w, c)
		}
		teamsJson := make([]mdl.TeamJson, len(teams))
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)
		return templateshlp.RenderJson(w, c, teamsJson)
	} else {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
}

// json new handler
func NewJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "Team New Handler: Error when decoding request body: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
		}

		var data TeamData
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Errorf(c, "Team New Handler: Error when decoding request body: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
		}

		if len(data.Name) <= 0 {
			log.Errorf(c, "Team New Handler: 'Name' field cannot be empty")
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeNameCannotBeEmpty)}
		} else if t := mdl.FindTeams(c, "KeyName", helpers.TrimLower(data.Name)); t != nil {
			log.Errorf(c, "Team New Handler: That team name already exists.")
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamAlreadyExists)}
		} else {
			team, err := mdl.CreateTeam(c, data.Name, u.Id, data.Visibility == "Private")
			if err != nil {
				log.Errorf(c, "Team New Handler: error when trying to create a team: %v", err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
			}
			// join the team
			if err = u.AddTeamId(c, team.Id); err != nil {
				log.Errorf(c, "Team New Handler: error when trying to create a team relationship: %v", err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
			}
			// publish new activity
			actor := mdl.ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
			object := mdl.ActivityEntity{ID: team.Id, Type: "team", DisplayName: team.Name}
			target := mdl.ActivityEntity{}
			mdl.Publish(c, "team", "created a new team", actor, object, target, u.Id)

			// return the newly created team
			var tJson mdl.TeamJson
			fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
			helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

			return templateshlp.RenderJson(w, c, tJson)
		}
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// json show handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Show Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var team *mdl.Team
		if team, err = mdl.TeamById(c, intID); err != nil {
			log.Errorf(c, "Team Show Handler: team not found: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		// get data for json team

		// build team json
		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		// build players json
		players := team.Players(c)
		fieldsToKeepForPlayer := []string{"Id", "Username"}
		playersJson := make([]mdl.UserJson, len(players))
		helpers.TransformFromArrayOfPointers(&players, &playersJson, fieldsToKeepForPlayer)

		teamData := struct {
			Team        mdl.TeamJson
			Joined      bool
			RequestSent bool
			Players     []mdl.UserJson
		}{
			tJson,
			team.Joined(c, u),
			mdl.WasTeamRequestSent(c, intID, u.Id),
			playersJson,
		}
		return templateshlp.RenderJson(w, c, teamData)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// json update handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {

		teamID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Update Handler: error when extracting permalink id: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotUpdate)}
		}

		if !mdl.IsTeamAdmin(c, teamID, u.Id) {
			log.Errorf(c, "Team Update Handler: user is not admin")
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamUpdateForbiden)}
		}

		var team *mdl.Team
		team, err = mdl.TeamById(c, teamID)
		if err != nil {
			log.Errorf(c, "Team Update handler: team not found. id: %v, err: %v", teamID, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotUpdate)}
		}
		// only work on name and private. Other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "Team Update handler: Error when reading request body err: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}

		var updatedData TeamData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			log.Errorf(c, "Team Update handler: Error when decoding request body err: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}

		updatedPrivate := updatedData.Visibility == "Private"

		if helpers.IsStringValid(updatedData.Name) && (updatedData.Name != team.Name || updatedPrivate != team.Private) {

			// be sure that team with that name does not exist in datastore
			if t := mdl.FindTeams(c, "KeyName", helpers.TrimLower(updatedData.Name)); t != nil {
				log.Errorf(c, "Team Update Handler: That team name already exists.")
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamAlreadyExists)}
			}
			// update data
			team.Name = updatedData.Name
			team.Private = updatedPrivate
			team.Update(c)
		} else {
			log.Errorf(c, "Cannot update because updated data are not valid")
			log.Errorf(c, "Update name = %s", updatedData.Name)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
		}
		// keep only needed fields for json api
		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// json destroy handler
func DestroyJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {

		teamID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Destroy Handler: error when extracting permalink id: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotDelete)}
		}

		if !mdl.IsTeamAdmin(c, teamID, u.Id) {
			log.Errorf(c, "Team Destroy Handler: user is not admin")
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamDeleteForbiden)}
		}
		var team *mdl.Team
		team, err = mdl.TeamById(c, teamID)
		if err != nil {
			log.Errorf(c, "Team Destroy handler: team not found. id: %v, err: %v", teamID, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotUpdate)}
		}

		// delete all team-user relationships
		for _, player := range team.Players(c) {
			if err := player.RemoveTeamId(c, team.Id); err != nil {
				log.Errorf(c, "Team Destroy Handler: error when trying to destroy team relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, tournament := range team.Tournaments(c) {
			if err := tournament.RemoveTeamId(c, team.Id); err != nil {
				log.Errorf(c, "Team Destroy Handler: error when trying to destroy tournament relationship: %v", err)
			}
		}
		// delete the team
		team.Destroy(c)

		// return destroyed status
		return templateshlp.RenderJson(w, c, "team has been destroyed")
	} else {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
}

// json invite handler
// use this handler when you wish to request an invitation to a team.
// this is done when the team in set as 'private' and the user wishes to join it.
func InviteJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Invite Handler: error when extracting permalink id: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFoundCannotInvite)}
		}

		if _, err := mdl.CreateTeamRequest(c, intID, u.Id); err != nil {
			log.Errorf(c, "Team Invite Handler: teams.Invite, error when trying to create a team request: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotInvite)}
		}
		// return destroyed status
		return templateshlp.RenderJson(w, c, "team request was created")
	}
	return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Json Allow handler
// use this handler to allow a request send by a user on a team.
// after this, the user that send the request will be part of the team
func AllowRequestJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		requestId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Allow Request Handler: teams.AllowRequest, id could not be extracter from url: %v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
		}

		if teamRequest, err := mdl.TeamRequestById(c, requestId); err == nil {
			// join user to the team
			var team *mdl.Team
			team, err = mdl.TeamById(c, teamRequest.TeamId)
			if err != nil {
				log.Errorf(c, "Team Allow Request handler: team not found. id: %v, err: %v", teamRequest.TeamId, err)
				return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
			}
			user, err := mdl.UserById(c, teamRequest.UserId)
			if err != nil {
				log.Errorf(c, "Team Allow Handler: user not found")
				return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
			}

			team.Join(c, user)
			// request is no more needed so clear it from datastore
			teamRequest.Destroy(c)

		} else {
			log.Errorf(c, "Team Allow Request Handler: cannot find team request with id=%d", requestId)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
		}

		return templateshlp.RenderJson(w, c, "team request was handled")

	} else {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
}

// Json Deny handler
// use this handler to deny a request send by a user on a team.
// the user will not be able to be part of the team
func DenyRequestJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		requestId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Deny Request Handler: teams.AllowRequest, id could not be extracter from url: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
		}
		if teamRequest, err := mdl.TeamRequestById(c, requestId); err != nil {
			log.Errorf(c, "Team Deny Request Handler: teams.AllowRequest, team request not found: %v", err)
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

// json search handler
// use this handler to search for a team.
func SearchJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	keywords := r.FormValue("q")
	if r.Method == "GET" && (len(keywords) > 0) {

		words := helpers.SetOfStrings(keywords)
		ids, err := mdl.GetTeamInvertedIndexes(c, words)
		if err != nil {
			log.Errorf(c, "Team Search Handler: teams.Index, error occurred when getting indexes of words: %v", err)
			data := struct {
				MessageDanger string `json:",omitempty"`
			}{
				"Oops! something went wrong, we are unable to perform search query.",
			}
			return templateshlp.RenderJson(w, c, data)
		}
		result := mdl.TeamScore(c, keywords, ids)
		log.Infof(c, "result from TeamScore: %v", result)
		teams := mdl.TeamsByIds(c, result)
		log.Infof(c, "ByIds result %v", teams)
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
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
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

// json team members handler
// use this handler to get members of a team.
func MembersJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	log.Infof(c, "json team members handler.")

	if r.Method == "GET" {
		teamId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "Team Members Handler: error extracting permalink err:%v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamMemberNotFound)}
		}
		team, err1 := mdl.TeamById(c, teamId)
		if err1 != nil {
			log.Errorf(c, "Team Allow Request handler: team not found. id: %v, err: %v", teamId, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamMemberNotFound)}
		}

		// build members json
		members := team.Players(c)
		fieldsToKeepForMember := []string{"Id", "Username"}
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

// json Join handler for team relations
func JoinJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get team id
		teamId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Join Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		var team *mdl.Team
		if team, err = mdl.TeamById(c, teamId); err != nil {
			log.Errorf(c, "Team Join Handler: team not found: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if err := team.Join(c, u); err != nil {
			log.Errorf(c, "Team Join Handler: error on Join team: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		// publish new activity
		actor := mdl.ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
		object := mdl.ActivityEntity{ID: team.Id, Type: "team", DisplayName: team.Name}
		target := mdl.ActivityEntity{}
		mdl.Publish(c, "team", "joined team", actor, object, target, u.Id)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// json destroy handler for team relations
func LeaveJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get team id
		teamId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Leave Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if mdl.IsTeamAdmin(c, teamId, u.Id) {
			log.Errorf(c, "Team Leave Handler: Team administrator cannot leave the team")
			return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeTeamAdminCannotLeave)}
		}

		var team *mdl.Team
		if team, err = mdl.TeamById(c, teamId); err != nil {
			log.Errorf(c, "Team Leave Handler: team not found: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		if err := team.Leave(c, u); err != nil {
			log.Errorf(c, "Team Leave Handler: error on Leave team: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		helpers.CopyToPointerStructure(team, &tJson)
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.KeepFields(&tJson, fieldsToKeep)

		// publish new activity
		actor := mdl.ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
		object := mdl.ActivityEntity{ID: team.Id, Type: "team", DisplayName: team.Name}
		target := mdl.ActivityEntity{}
		mdl.Publish(c, "team", "left team", actor, object, target, u.Id)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
