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

package tournaments

import (
	"errors"
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	mdl "github.com/santiaago/purple-wing/models"
)

// A GroupJson is a variable to hold a the name of a group and an array of Teams.
// We use it to group tournament teams information by group to meet world cup organization.
type GroupJson struct {
	Name  string
	Teams []TeamJson
}

// A TeamJson is a variable to hold the basic information of a Team:
// The name of the team, the number of points recorded in the group phase, the goals for and against.
type TeamJson struct {
	Name   string
	Points int64
	GoalsF int64
	GoalsA int64
}

// json tournament groups handler
// use this handler to get groups of a tournament.
func GroupsJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Group Handler:"

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var tournament *mdl.Tournament
		tournament, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "%s tournament with id:%v was not found %v", desc, tournamentId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		groups := mdl.Groups(c, tournament.GroupIds)
		groupsJson := formatGroupsJson(groups)

		data := struct {
			Groups []GroupJson
		}{
			groupsJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Format a TGroup array into a GroupJson array.
func formatGroupsJson(groups []*mdl.Tgroup) []GroupJson {

	groupsJson := make([]GroupJson, len(groups))
	for i, g := range groups {
		groupsJson[i].Name = g.Name
		teams := make([]TeamJson, len(g.Teams))
		for j, t := range g.Teams {
			teams[j].Name = t.Name
			teams[j].Points = g.Points[j]
			teams[j].GoalsF = g.GoalsF[j]
			teams[j].GoalsA = g.GoalsA[j]

		}
		groupsJson[i].Teams = teams
	}
	return groupsJson
}
