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

package users

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	mdl "github.com/santiaago/purple-wing/models"
)

type UserData struct {
	Username string
	Name     string
	Email    string
}

// json index user handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		if !u.IsAdmin {
			log.Errorf(c, "Tournament Index Handler: user is not admin, User list can only be shown for admin users.")
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotFound)}
		}

		users := mdl.FindAllUsers(c)

		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created"}
		usersJson := make([]mdl.UserJson, len(users))
		helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, usersJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Json show user handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		var userId int64
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "User Show Handler: error when extracting permalink for url: %v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}
		userId = intID

		// user
		var user *mdl.User
		user, err = mdl.UserById(c, userId)
		if err != nil {
			log.Errorf(c, "User Show Handler: user not found")
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created", "IsAdmin", "Auth"}
		var uJson mdl.UserJson
		helpers.InitPointerStructure(user, &uJson, fieldsToKeep)

		// get with param:
		with := r.FormValue("including")
		params := helpers.SetOfStrings(with)
		var teams []*mdl.Team
		var teamRequests []*mdl.TeamRequest
		var tournaments []*mdl.Tournament
		for _, param := range params {
			switch param {
			case "teams":
				teams = user.Teams(c)
			case "teamrequests":
				teamRequests = mdl.TeamsRequests(c, teams)
			case "tournaments":
				tournaments = user.Tournaments(c)
			}
		}

		// teams
		teamsFieldsToKeep := []string{"Id", "Name"}
		teamsJson := make([]mdl.TeamJson, len(teams))
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, teamsFieldsToKeep)
		// tournaments
		tournamentfieldsToKeep := []string{"Id", "Name"}
		tournamentsJson := make([]mdl.TournamentJson, len(tournaments))
		helpers.TransformFromArrayOfPointers(&tournaments, &tournamentsJson, tournamentfieldsToKeep)
		// team requests
		teamRequestFieldsToKeep := []string{"Id", "TeamId", "UserId"}
		trsJson := make([]mdl.TeamRequestJson, len(teamRequests))
		helpers.TransformFromArrayOfPointers(&teamRequests, &trsJson, teamRequestFieldsToKeep)

		// data
		data := struct {
			User         mdl.UserJson          `json:",omitempty"`
			Teams        []mdl.TeamJson        `json:",omitempty"`
			TeamRequests []mdl.TeamRequestJson `json:",omitempty"`
			Tournaments  []mdl.TournamentJson  `json:",omitempty"`
		}{
			uJson,
			teamsJson,
			trsJson,
			tournamentsJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// json update user handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		userId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "User Update Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFoundCannotUpdate)}
		}
		if userId != u.Id {
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
		}

		// only work on name other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "User Show handler: Error when reading request body err: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
		}

		var updatedData UserData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			log.Errorf(c, "User Show handler: Error when decoding request body err: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
		}
		if helpers.IsEmailValid(updatedData.Email) && updatedData.Email != u.Email {
			u.Email = updatedData.Email
			u.Update(c)
		}
		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created"}
		var uJson mdl.UserJson
		helpers.InitPointerStructure(u, &uJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, uJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
