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

package tournamentrels

import (
	"errors"
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	// mdl "github.com/santiaago/purple-wing/models/tournament"
	// mdl "github.com/santiaago/purple-wing/models/user"
	mdl "github.com/santiaago/purple-wing/models"
	activitymdl "github.com/santiaago/purple-wing/models/activity"
)

// json create handler for tournament relationships
func CreateJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get tournament id
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Tournamentrels Create Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var tournament *mdl.Tournament
		if tournament, err = mdl.TournamentById(c, tournamentId); err != nil {
			log.Errorf(c, "Tournamentrels Create Handler: tournament not found: %v", err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		if err := tournament.Join(c, u); err != nil {
			log.Errorf(c, "Tournamentrels Create Handler: error on Join tournament: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TournamentJson
		fieldsToKeep := []string{"Id", "Name"}
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		// publish new activity
		actor := activitymdl.ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
		object := activitymdl.ActivityEntity{ID: tournament.Id, Type: "tournament", DisplayName: tournament.Name}
		target := activitymdl.ActivityEntity{}
		activitymdl.Publish(c, "tournament", "joined tournament", actor, object, target, u.Id)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json destroy handler for tournament relationships
func DestroyJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get tournament id
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Tournamentrels Destroy Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var tournament *mdl.Tournament
		if tournament, err = mdl.TournamentById(c, tournamentId); err != nil {
			log.Errorf(c, "Tournamentrels Destroy Handler: tournament not found: %v", err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		if err := tournament.Leave(c, u); err != nil {
			log.Errorf(c, "Tournamentrels Destroy Handler: error on Leave team: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeInternal)}
		}

		// return the left tournament
		var tJson mdl.TournamentJson
		fieldsToKeep := []string{"Id", "Name"}
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		// publish new activity
		actor := activitymdl.ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
		object := activitymdl.ActivityEntity{ID: tournament.Id, Type: "tournament", DisplayName: tournament.Name}
		target := activitymdl.ActivityEntity{}
		activitymdl.Publish(c, "tournament", "left tournament", actor, object, target, u.Id)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}
