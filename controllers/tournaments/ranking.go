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

// Json tournament ranking handler:
// Use this handler to get the ranking of a tournament.
// The ranking is an array of users (members) or teams,
// You can specify the rankby parameter to be "users" or "teams".
func RankingJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Ranking Handler:"

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var t *mdl.Tournament
		t, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "%s tournament with id:%v was not found %v", desc, tournamentId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		rankby := r.FormValue("rankby")
		// if wrong data, we set rankby to "users"
		if rankby != "teams" && rankby != "users" {
			rankby = "users"
		}

		if rankby == "users" {
			log.Infof(c, "%s ready to build a user array", desc)
			users := t.RankingByUser(c)

			fieldsToKeep := []string{"Id", "Name", "Score"}
			usersJson := make([]mdl.UserJson, len(users))
			helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

			data := struct {
				Users []mdl.UserJson
			}{
				usersJson,
			}

			return templateshlp.RenderJson(w, c, data)

		} else if rankby == "teams" {
			log.Infof(c, "%s ready to build team array", desc)
			teams := t.RankingByTeam(c)

			fieldsToKeep := []string{"Id", "Name", "Score"}
			teamsJson := make([]mdl.TeamJson, len(teams))
			helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)

			data := struct {
				Teams []mdl.TeamJson
			}{
				teamsJson,
			}
			return templateshlp.RenderJson(w, c, data)
		}
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
