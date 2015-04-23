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

package teams

import (
	"errors"
	"net/http"
	"strconv"

	"appengine"

	"github.com/taironas/route"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Accuracies handler, use it to get the accuracies of a team.
//
// Use this handler to get the accuracies of a team.
//	GET	/j/teams/:teamId/accuracies	retrieves all the tournament accuracies of a team with the given id.
//
// The response is an array of accurracies for the specified team group by tournament with the last 5 progressions.
func Accuracies(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Accuracies Handler:"

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

	var t *mdl.Team
	t, err = mdl.TeamById(c, teamId)
	if err != nil {
		log.Errorf(c, "%s team with id:%v was not found %v", desc, teamId, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}

	log.Infof(c, "%s ready to build a acc array", desc)
	accs := t.AccuraciesGroupByTournament(c, 5)

	data := struct {
		Accuracies *[]mdl.AccuracyOverall
	}{
		accs,
	}

	return templateshlp.RenderJson(w, c, data)
}

// Team accuracies by tournament handler:
//
// Use this handler to get the accuracies of a team for a specific tournament.
//	GET	/j/teams/:teamId/accuracies/:tournamentId	retrieves accuracies of a team with the given id for the specified tournament.
//
// The response is an array of accurracies for the specified team team group by tournament with all it's progressions.
func AccuracyByTournament(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Accuracies by tournament Handler:"

	if r.Method == "GET" {
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

		var t *mdl.Team
		t, err = mdl.TeamById(c, teamId)
		if err != nil {
			log.Errorf(c, "%s team with id:%v was not found %v", desc, teamId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		strTournamentId, err := route.Context.Get(r, "tournamentId")
		if err != nil {
			log.Errorf(c, "%s error getting tournament id, err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var tournamentId int64
		tournamentId, err = strconv.ParseInt(strTournamentId, 0, 64)
		if err != nil {
			log.Errorf(c, "%s error converting tournament id from string to int64, err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var tour *mdl.Tournament
		tour, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "%s tournament with id:%v was not found %v", desc, tournamentId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		log.Infof(c, "%s ready to build a acc array", desc)
		acc := t.AccuracyByTournament(c, tour)

		data := struct {
			Accuracy *mdl.AccuracyOverall
		}{
			acc,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
