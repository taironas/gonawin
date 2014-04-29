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

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	mdl "github.com/santiaago/purple-wing/models"
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
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		if !u.IsAdmin {
			log.Errorf(c, "User Index Handler: user is not admin, User list can only be shown for admin users.")
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotFound)}
		}

		users := mdl.FindAllUsers(c)

		fieldsToKeep := []string{"Id", "Username", "Name", "Alias", "Email", "Created"}
		usersJson := make([]mdl.UserJson, len(users))
		helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, usersJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// User show handler.
// including parameter: {teams, tournaments, teamrequests}
// count parameter: default 12
// page parameter: default 1
func Show(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "User show handler:"
	if r.Method == "GET" {
		var userId int64
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink for url: %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}
		userId = intID

		// user
		var user *mdl.User
		user, err = mdl.UserById(c, userId)
		log.Infof(c, "User: %v", user)
		log.Infof(c, "User: %v", user.TeamIds)
		if err != nil {
			log.Errorf(c, "%s user not found", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		fieldsToKeep := []string{"Id", "Username", "Name", "Alias", "Email", "Created", "IsAdmin", "Auth", "TeamIds", "TournamentIds", "Score"}
		var uJson mdl.UserJson
		helpers.InitPointerStructure(user, &uJson, fieldsToKeep)
		log.Infof(c, "%s User: %v", desc, uJson.TeamIds)
		log.Infof(c, "%s User: %v", desc, uJson)
		log.Infof(c, "%s User: %v", desc, *uJson.TeamIds)

		// get with param:
		with := r.FormValue("including")
		params := helpers.SetOfStrings(with)
		var teams []*mdl.Team
		var teamRequests []*mdl.TeamRequest
		var tournaments []*mdl.Tournament
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
			}
		}

		// teams
		teamsFieldsToKeep := []string{"Id", "Name", "Accuracy", "MembersCount", "Private"}
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

// User update handler.
func Update(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "User update handler:"
	if r.Method == "POST" {
		userId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v.", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFoundCannotUpdate)}
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
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// User destroy handler
//	POST	/j/user/destroy/[0-9]+/		Destroys the user with the given id.
//
func Destroy(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "User Destroy Handler:"
	if r.Method == "POST" {
		userId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFoundCannotDelete)}
		}
		if userId != u.Id {
			log.Errorf(c, "%s error user ids do not match. url id:%s user id: %s", desc, userId, u.Id)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFoundCannotDelete)}
		}
    
		// send task to delete activities of the user.
		log.Infof(c, "%s Sending to taskqueue: delete activities", desc)

		bactivityIds, err1 := json.Marshal(u.ActivityIds)
		if err1 != nil {
			log.Errorf(c, "%s Error marshaling", desc, err1)
		}

		task := taskqueue.NewPOSTTask("/a/publish/users/deleteactivities/", url.Values{
			"activity_ids": []string{string(bactivityIds)},
		})

		if _, err := taskqueue.Add(c, task, ""); err != nil {
			log.Errorf(c, "%s unable to add task to taskqueue.", desc)
		} else {
			log.Infof(c, "%s add task to taskqueue successfully", desc)
		}

		// user
		var user *mdl.User
		user, err = mdl.UserById(c, userId)
		if err != nil {
			log.Errorf(c, "%s user not found", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		// delete all team-user relationships
		for _, teamId := range user.TeamIds {
			if mdl.IsTeamAdmin(c, teamId, u.Id) {
				var team *mdl.Team
				if team, err = mdl.TeamById(c, teamId); err != nil {
					log.Errorf(c, "%s team %d not found", desc, teamId)
					return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
				}
				if err = team.RemoveAdmin(c, userId); err != nil {
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
				if err = tournament.RemoveAdmin(c, userId); err != nil {
					log.Infof(c, "%s error occurred during admin deletion: %v", desc, err)
					return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserIsTournamentAdminCannotDelete)}
				}
			} else {
				if err := user.RemoveTournamentId(c, tournamentId); err != nil {
					log.Errorf(c, "%s error when trying to destroy tournament relationship: %v", desc, err)
				}
			}
		}
		// delete the user
		user.Destroy(c)

		msg := fmt.Sprintf("The user %s was correctly deleted!", user.Username)
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

// User teams  handler.
// count parameter: default 12
// page parameter: default 1
func Teams(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "User joined teams handler:"

	if r.Method == "GET" {
		var userId int64
		intID, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink for url: %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}
		userId = intID

		// user
		var user *mdl.User
		user, err = mdl.UserById(c, userId)
		log.Infof(c, "User: %v", user)
		log.Infof(c, "User: %v", user.TeamIds)
		if err != nil {
			log.Errorf(c, "%s user not found", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
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
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// User tournaments  handler.
// count parameter: default 25
// page parameter: default 1
func Tournaments(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "User joined teams handler:"

	if r.Method == "GET" {
		var userId int64
		intID, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink for url: %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}
		userId = intID

		// user
		var user *mdl.User
		user, err = mdl.UserById(c, userId)
		log.Infof(c, "User: %v", user)
		log.Infof(c, "User: %v", user.TeamIds)
		if err != nil {
			log.Errorf(c, "%s user not found", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
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
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
