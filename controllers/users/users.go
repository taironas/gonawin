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
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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

type Form struct {
	Username      string
	Name          string
	Email         string
	ErrorUsername string
	ErrorName     string
	ErrorEmail    string
}

type UserData struct {
	Username string
	Name     string
	Email    string
}

// Show handler
func Show(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	userId, err := handlers.PermalinkID(r, c, 3)
	if err != nil {
		http.Redirect(w, r, "/m/users/", http.StatusFound)
		return
	}

	funcs := template.FuncMap{
		"Profile": func() bool { return true },
	}

	t := template.Must(template.New("tmpl_user_show").
		Funcs(funcs).
		ParseFiles("templates/user/show.html",
		"templates/user/info.html",
		"templates/user/teams.html",
		"templates/user/tournaments.html",
		"templates/user/requests.html"))

	var user *usermdl.User
	user, err = usermdl.ById(c, userId)
	if err != nil {
		helpers.Error404(w)
		return
	}

	teams := usermdl.Teams(c, userId)
	tournaments := tournamentrelshlp.Tournaments(c, userId)
	teamRequests := teamrelshlp.TeamsRequests(c, teams)

	userData := struct {
		User         *usermdl.User
		Teams        []*teammdl.Team
		Tournaments  []*tournamentmdl.Tournament
		TeamRequests []*teamrequestmdl.TeamRequest
	}{
		user,
		teams,
		tournaments,
		teamRequests,
	}

	templateshlp.RenderWithData(w, r, c, t, userData, funcs, "renderUserShow")
}

// json index user handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		users := usermdl.FindAll(c)

		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created"}
		usersJson := make([]usermdl.UserJson, len(users))
		helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, usersJson)
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// Json show user handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)
	log.Infof(c, "User Show Json Handler")

	if r.Method == "GET" {
		var userId int64
		strUrl := r.URL.String()
		if strings.Contains(strUrl, "?") {
			path := strings.Split(r.URL.String(), "/")
			strPath := path[4]
			strID := strPath[0:strings.Index(strPath, "?")]
			intID, err := strconv.ParseInt(strID, 0, 64)
			if err != nil {
				log.Errorf(c, " error when calling PermalinkID with %v.Error: %v", strID, err)
			}
			userId = intID

		} else {
			intID, err := handlers.PermalinkID(r, c, 4)
			if err != nil {
				return helpers.BadRequest{err}
			}
			userId = intID
		}

		// user
		var user *usermdl.User
		user, err := usermdl.ById(c, userId)
		if err != nil {
			return helpers.BadRequest{err}
		}
		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created", "IsAdmin", "Auth"}
		var uJson usermdl.UserJson
		helpers.InitPointerStructure(user, &uJson, fieldsToKeep)

		// with param:
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
			User         usermdl.UserJson
			Teams        []teammdl.TeamJson
			TeamRequests []teamrequestmdl.TeamRequestJson
			Tournaments  []tournamentmdl.TournamentJson
		}{
			uJson,
			teamsJson,
			trsJson,
			tournamentsJson,
		}

		return templateshlp.RenderJson(w, c, data)
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// json update user handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		userId, err := handlers.PermalinkID(r, c, 4)

		if err != nil {
			return helpers.BadRequest{err}
		}
		if userId != u.Id {
			return helpers.BadRequest{errors.New("User cannot be updated")}
		}

		// only work on name other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when reading request body")}
		}

		var updatedData UserData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when decoding request body")}
		}
		if helpers.IsEmailValid(updatedData.Email) && updatedData.Email != u.Email {
			u.Email = updatedData.Email
			usermdl.Update(c, u)
		}
		fieldsToKeep := []string{"Id", "Username", "Name", "Email", "Created"}
		var uJson usermdl.UserJson
		helpers.InitPointerStructure(u, &uJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, uJson)
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}
