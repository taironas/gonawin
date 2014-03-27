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

// JSON Join handler for tournament relationships
func Join(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Join Handler:"

	if r.Method == "POST" {
		// get tournament id
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var tournament *mdl.Tournament
		if tournament, err = mdl.TournamentById(c, tournamentId); err != nil {
			log.Errorf(c, "%s tournament not found: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		if err := tournament.Join(c, u); err != nil {
			log.Errorf(c, "%s error on Join tournament: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TournamentJson
		fieldsToKeep := []string{"Id", "Name"}
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		// publish new activity
		object := mdl.ActivityEntity{ID: tournament.Id, Type: "tournament", DisplayName: tournament.Name}
		target := mdl.ActivityEntity{}
		u.Publish(c, "tournament", "joined tournament", object, target)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// JSON Leave handler for tournament relationships
func Leave(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Leave Handler:"

	if r.Method == "POST" {
		// get tournament id
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "%s error when extracting permalink id: %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var tournament *mdl.Tournament
		if tournament, err = mdl.TournamentById(c, tournamentId); err != nil {
			log.Errorf(c, "%s tournament not found: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		if err := tournament.Leave(c, u); err != nil {
			log.Errorf(c, "%s error on Leave team: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		// return the left tournament
		var tJson mdl.TournamentJson
		fieldsToKeep := []string{"Id", "Name"}
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		// publish new activity
		object := mdl.ActivityEntity{ID: tournament.Id, Type: "tournament", DisplayName: tournament.Name}
		target := mdl.ActivityEntity{}
		u.Publish(c, "tournament", "left tournament", object, target)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// JSON Join as Team handler for tournament teams realtionship.
func JoinAsTeam(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Join as a Team Handler:"

	if r.Method == "POST" {
		// get tournament and team id
		tournamentId, err1 := handlers.PermalinkID(r, c, 4)
		teamId, err2 := handlers.PermalinkID(r, c, 5)
		if err1 != nil || err2 != nil {
			log.Errorf(c, "%s string value could not be parsed: %v, %v", desc, err1, err2)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tournament *mdl.Tournament
		if tournament, err1 = mdl.TournamentById(c, tournamentId); err1 != nil {
			log.Errorf(c, "%stournament with id: %v was not found %v", desc, tournamentId, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var team *mdl.Team
		if team, err1 = mdl.TeamById(c, teamId); err1 != nil {
			log.Errorf(c, "%s team not found: %v", desc, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if err := tournament.TeamJoin(c, team); err != nil {
			log.Errorf(c, "%s error when trying to join team: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TournamentJson
		fieldsToKeep := []string{"Id", "Name"}
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		// publish new activity
		object := mdl.ActivityEntity{ID: tournament.Id, Type: "tournament", DisplayName: tournament.Name}
		target := mdl.ActivityEntity{}
		team.Publish(c, "tournament", "joined tournament", object, target)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// JSON Leave as Team handler for tournament teams realtionship.
func LeaveAsTeam(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Leave as a Team Handler:"

	if r.Method == "POST" {

		// get tournament and team id
		tournamentId, err1 := handlers.PermalinkID(r, c, 4)
		teamId, err2 := handlers.PermalinkID(r, c, 5)
		if err1 != nil || err2 != nil {
			log.Errorf(c, "%s string value could not be parsed: %v, %v", desc, err1, err2)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
		}
		var tournament *mdl.Tournament
		if tournament, err1 = mdl.TournamentById(c, tournamentId); err1 != nil {
			log.Errorf(c, "%s tournament with id: %v was not found %v", desc, tournamentId, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var team *mdl.Team
		if team, err1 = mdl.TeamById(c, teamId); err1 != nil {
			log.Errorf(c, "team not found: %v", desc, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		// leave team
		if err := tournament.TeamLeave(c, team); err != nil {
			log.Errorf(c, "%s error when trying to leave team: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}
		// return the left tournament

		var tJson mdl.TournamentJson
		fieldsToKeep := []string{"Id", "Name"}
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		// publish new activity
		object := mdl.ActivityEntity{ID: tournament.Id, Type: "tournament", DisplayName: tournament.Name}
		target := mdl.ActivityEntity{}
		u.Publish(c, "tournament", "left tournament", object, target)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
