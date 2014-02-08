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
	teamrelshlp "github.com/santiaago/purple-wing/helpers/teamrels"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	tournamentrelshlp "github.com/santiaago/purple-wing/helpers/tournamentrels"

	teammdl "github.com/santiaago/purple-wing/models/team"
	teamrequestmdl "github.com/santiaago/purple-wing/models/teamrequest"
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

type UserData struct {
	Username string
	Name     string
	Email    string
}

// json index user handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		if !u.IsAdmin {
			log.Errorf(c, "Tournament Index Handler: user is not admin, User list can only be shown for admin users.")
			return helpers.BadRequest{errors.New(helpers.ErrorCodeNotFound)}
		}

		users := usermdl.FindAll(c)

		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created"}
		usersJson := make([]usermdl.UserJson, len(users))
		helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, usersJson)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// Json show user handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		var userId int64
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "User Show Handler: error when extracting permalink for url: %v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeUserNotFound)}
		}
		userId = intID

		// user
		var user *usermdl.User
		user, err = usermdl.ById(c, userId)
		if err != nil {
			log.Errorf(c, "User Show Handler: user not found")
			return helpers.NotFound{errors.New(helpers.ErrorCodeUserNotFound)}
		}

		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created", "IsAdmin", "Auth"}
		var uJson usermdl.UserJson
		helpers.InitPointerStructure(user, &uJson, fieldsToKeep)

		// get with param:
		with := r.FormValue("including")
		params := helpers.SetOfStrings(with)
		var teams []*teammdl.Team
		var teamRequests []*teamrequestmdl.TeamRequest
		var tournaments []*tournamentmdl.Tournament
		for _, param := range params {
			switch param {
			case "teams":
				teams = usermdl.Teams(c, userId)
			case "teamrequests":
				teamRequests = teamrelshlp.TeamsRequests(c, teams)
			case "tournaments":
				tournaments = tournamentrelshlp.Tournaments(c, userId)
			}
		}

		// teams
		teamsFieldsToKeep := []string{"Id", "Name"}
		teamsJson := make([]teammdl.TeamJson, len(teams))
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, teamsFieldsToKeep)
		// tournaments
		tournamentfieldsToKeep := []string{"Id", "Name"}
		tournamentsJson := make([]tournamentmdl.TournamentJson, len(tournaments))
		helpers.TransformFromArrayOfPointers(&tournaments, &tournamentsJson, tournamentfieldsToKeep)
		// team requests
		teamRequestFieldsToKeep := []string{"Id", "TeamId", "UserId"}
		trsJson := make([]teamrequestmdl.TeamRequestJson, len(teamRequests))
		helpers.TransformFromArrayOfPointers(&teamRequests, &trsJson, teamRequestFieldsToKeep)

		// data
		data := struct {
			User         usermdl.UserJson                 `json:",omitempty"`
			Teams        []teammdl.TeamJson               `json:",omitempty"`
			TeamRequests []teamrequestmdl.TeamRequestJson `json:",omitempty"`
			Tournaments  []tournamentmdl.TournamentJson   `json:",omitempty"`
		}{
			uJson,
			teamsJson,
			trsJson,
			tournamentsJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json update user handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		userId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "User Update Handler: error when extracting permalink id: %v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeUserNotFoundCannotUpdate)}
		}
		if userId != u.Id {
			return helpers.BadRequest{errors.New(helpers.ErrorCodeUserCannotUpdate)}
			return helpers.BadRequest{errors.New("User cannot be updated")}
		}

		// only work on name other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "User Show handler: Error when reading request body err: %v", err)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeUserCannotUpdate)}
		}

		var updatedData UserData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			log.Errorf(c, "User Show handler: Error when decoding request body err: %v", err)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeUserCannotUpdate)}
		}
		if helpers.IsEmailValid(updatedData.Email) && updatedData.Email != u.Email {
			u.Email = updatedData.Email
			usermdl.Update(c, u)
		}
		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created"}
		var uJson usermdl.UserJson
		helpers.InitPointerStructure(u, &uJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, uJson)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}
