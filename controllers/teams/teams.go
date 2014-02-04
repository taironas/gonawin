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
	"io/ioutil"
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
	"github.com/santiaago/purple-wing/helpers/handlers"
	teamrelshlp "github.com/santiaago/purple-wing/helpers/teamrels"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	tournamentrelshlp "github.com/santiaago/purple-wing/helpers/tournamentrels"

	searchmdl "github.com/santiaago/purple-wing/models/search"
	teammdl "github.com/santiaago/purple-wing/models/team"
	teaminvidmdl "github.com/santiaago/purple-wing/models/teamInvertedIndex"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
	teamrequestmdl "github.com/santiaago/purple-wing/models/teamrequest"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

type TeamData struct {
	Name       string
	Visibility string
}

// json index handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		teams := teammdl.FindAll(c)
		if len(teams) == 0 {
			return templateshlp.RenderEmptyJsonArray(w, c)
		}
		teamsJson := make([]teammdl.TeamJson, len(teams))
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)
		return templateshlp.RenderJson(w, c, teamsJson)
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}

// json new handler
func NewJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when reading request body")}
		}

		var data TeamData
		err = json.Unmarshal(body, &data)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when decoding request body")}
		}

		if len(data.Name) <= 0 {
			return helpers.InternalServerError{errors.New("'Name' field cannot be empty")}
		} else if t := teammdl.Find(c, "KeyName", helpers.TrimLower(data.Name)); t != nil {
			return helpers.InternalServerError{errors.New("That team name already exists")}
		} else {
			team, err := teammdl.Create(c, data.Name, u.Id, data.Visibility == "Private")
			if err != nil {
				log.Errorf(c, "error when trying to create a team: %v", err)
				return helpers.InternalServerError{errors.New("error when trying to create a team")}
			}
			// join the team
			_, err = teamrelmdl.Create(c, team.Id, u.Id)
			if err != nil {
				log.Errorf(c, " error when trying to create a team relationship: %v", err)
				return helpers.InternalServerError{errors.New("error when trying to create a team relationship")}
			}
			// return the newly created team
			var tJson teammdl.TeamJson
			fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
			helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

			return templateshlp.RenderJson(w, c, tJson)
		}
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}

// json show handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	intID, err := handlers.PermalinkID(r, c, 4)
	if err != nil {
		return helpers.NotFound{err}
	}

	if r.Method == "GET" {
		var team *teammdl.Team
		if team, err = teammdl.ById(c, intID); err != nil {
			return helpers.NotFound{err}
		}
		// get data for json team

		// build team json
		var tJson teammdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		// build players json
		players := teamrelshlp.Players(c, intID)
		fieldsToKeepForPlayer := []string{"Id", "Username"}
		playersJson := make([]usermdl.UserJson, len(players))
		helpers.TransformFromArrayOfPointers(&players, &playersJson, fieldsToKeepForPlayer)

		teamData := struct {
			Team        teammdl.TeamJson
			Joined      bool
			RequestSent bool
			Players     []usermdl.UserJson
		}{
			tJson,
			teammdl.Joined(c, intID, u.Id),
			teamrequestmdl.Sent(c, intID, u.Id),
			playersJson,
		}
		return templateshlp.RenderJson(w, c, teamData)
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}

// json update handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	teamID, err := handlers.PermalinkID(r, c, 4)
	if err != nil {
		return helpers.NotFound{err}
	}

	if r.Method == "POST" {
		if !teammdl.IsTeamAdmin(c, teamID, u.Id) {
			return helpers.BadRequest{errors.New("Team can only be updated by the team administrator")}
		}

		var team *teammdl.Team
		team, err = teammdl.ById(c, teamID)
		if err != nil {
			log.Errorf(c, " Team Edit handler: team not found. id: %v", teamID)
			return helpers.NotFound{err}
		}
		// only work on name and private. Other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when reading request body")}
		}

		var updatedData TeamData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when decoding request body")}
		}

		updatedPrivate := updatedData.Visibility == "Private"

		if helpers.IsStringValid(updatedData.Name) && (updatedData.Name != team.Name || updatedPrivate != team.Private) {
			team.Name = updatedData.Name
			team.Private = updatedPrivate
			teammdl.Update(c, teamID, team)
		} else {
			log.Errorf(c, "Cannot update because updated data are not valid")
			log.Errorf(c, "Update name = %s", updatedData.Name)
		}
		// keep only needed fields for json api
		var tJson teammdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tJson)
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// json destroy handler
func DestroyJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	teamID, err := handlers.PermalinkID(r, c, 4)
	if err != nil {
		return helpers.NotFound{err}
	}

	if r.Method == "POST" {
		if !teammdl.IsTeamAdmin(c, teamID, u.Id) {
			return helpers.BadRequest{errors.New("Team can only be deleted by the team administrator")}
		}

		// delete all team-user relationships
		for _, player := range teamrelshlp.Players(c, teamID) {
			if err := teamrelmdl.Destroy(c, teamID, player.Id); err != nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, tournament := range tournamentrelshlp.Teams(c, teamID) {
			if err := tournamentteamrelmdl.Destroy(c, tournament.Id, teamID); err != nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the team
		teammdl.Destroy(c, teamID)

		// return destroyed status
		return templateshlp.RenderJson(w, c, "team has been destroyed")
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}

// json invite handler
// use this handler when you wish to request an invitation to a team.
// this is done when the team in set as 'private' and the user wishes to join it.
func InviteJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			return helpers.NotFound{err}
		}

		if _, err := teamrequestmdl.Create(c, intID, u.Id); err != nil {
			log.Errorf(c, " teams.Invite, error when trying to create a team request: %v", err)
			return helpers.InternalServerError{errors.New("Error when sending invite.")}
		}
		// return destroyed status
		return templateshlp.RenderJson(w, c, "team request was created")
	}
	return helpers.NotFound{errors.New("Not supported.")}
}

// Json Allow handler
// use this handler to allow a request send by a user on a team.
// after this, the user that that send the request will be part of the team
func AllowRequestJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		requestId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, " teams.AllowRequest, id could not be extracter from url: %v", err)
			return helpers.NotFound{err}
		}

		if teamRequest, err := teamrequestmdl.ById(c, requestId); err == nil {
			// join user to the team
			teammdl.Join(c, teamRequest.TeamId, teamRequest.UserId)
		} else {
			appengine.NewContext(r).Errorf(" cannot find team request with id=%d", requestId)
		}
		// request is no more needed so clear it from datastore
		teamrequestmdl.Destroy(c, requestId)

		return templateshlp.RenderJson(w, c, "team request was handled")

	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// Json Deny handler
// use this handler to deny a request send by a user on a team.
// the user will not be able to be part of the team
func DenyRequestJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		requestId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, " teams.AllowRequest, id could not be extracter from url: %v", err)
			return helpers.NotFound{err}
		}

		// request is no more needed so clear it from datastore
		teamrequestmdl.Destroy(c, requestId)

		return templateshlp.RenderJson(w, c, "team request was handled")

	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// json search handler
// use this handler to search for a team.
func SearchJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)
	log.Infof(c, "json search handler.")
	keywords := r.FormValue("q")
	if r.Method == "GET" && (len(keywords) > 0) {
		words := helpers.SetOfStrings(keywords)
		ids, err := teaminvidmdl.GetIndexes(c, words)
		if err != nil {
			log.Errorf(c, " teams.Index, error occurred when getting indexes of words: %v", err)
		}
		result := searchmdl.TeamScore(c, keywords, ids)
		log.Infof(c, "result from TeamScore: %v", result)
		teams := teammdl.ByIds(c, result)
		log.Infof(c, "ByIds result %v", teams)
		if len(teams) == 0 {
			return templateshlp.RenderEmptyJson(w, c)
		}
		// filter team information to return in json api
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		teamsJson := make([]teammdl.TeamJson, len(teams))
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)
		// we should not directly return an array. so we add an extra layer.
		data := struct{
			Teams []teammdl.TeamJson `json:",omitempty"`
		}{
			teamsJson,
		}
		return templateshlp.RenderJson(w, c, data)
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}

// json team members handler
// use this handler to get members of a team.
func MembersJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)
	log.Infof(c, "json team members handler.")
	
	teamId, err := handlers.PermalinkID(r, c, 3)
	if err != nil {
		return helpers.NotFound{err}
	}
	
	if r.Method == "GET" {
		// build members json
		members := teamrelshlp.Players(c, teamId)
		fieldsToKeepForMember := []string{"Id", "Username"}
		membersJson := make([]usermdl.UserJson, len(members))
		helpers.TransformFromArrayOfPointers(&members, &membersJson, fieldsToKeepForMember)
	
		data := struct {
			Members     []usermdl.UserJson
		}{
			membersJson,
		}
		return templateshlp.RenderJson(w, c, data)
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}
