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
	"strconv"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

// create handler for tournament relationships
func Create(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get tournament id
	tournamentId, err := strconv.ParseInt(r.FormValue("TournamentId"), 10, 64)
	if err != nil {
		log.Errorf(c, " tournaments.Create, string value could not be parsed: %v", err)
	}

	if r.Method == "POST" {
		if err := tournamentmdl.Join(c, tournamentId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " tournamentrels.Create: %v", err)
		}
	}

	http.Redirect(w, r, "/m/tournaments/"+r.FormValue("TournamentId"), http.StatusFound)
}

// destroy handler for tournament relationships
func Destroy(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get tournament id
	tournamentId, err := strconv.ParseInt(r.FormValue("TournamentId"), 10, 64)
	if err != nil {
		log.Errorf(c, " tournaments.Destroy, string value could not be parsed: %v", err)
	}

	if r.Method == "POST" {
		if err := tournamentmdl.Leave(c, tournamentId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " tournamentrels.Destroy: %v", err)
		}
	}

	http.Redirect(w, r, "/m/tournaments/"+r.FormValue("TournamentId"), http.StatusFound)
}

// json create handler for tournament relationships
func CreateJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get tournament id
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, " tournaments.Create, string value could not be parsed: %v", err)
			return helpers.NotFound{err}
		}

		if err := tournamentmdl.Join(c, tournamentId, u.Id); err != nil {
			log.Errorf(c, " tournamentrels.Create: %v", err)
			return helpers.InternalServerError{err}
		}
		// return the joined tournament
		var tournament *tournamentmdl.Tournament
		if tournament, err = tournamentmdl.ById(c, tournamentId); err != nil {
			return helpers.NotFound{err}
		}

		var tJson tournamentmdl.TournamentJson
		helpers.CopyToPtrBasedStructGeneric(tournament, &tJson)
		fieldsToKeep := []string{"Id", "Name"}
		helpers.KeepFields(&tJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tJson)
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// json destroy handler for tournament relationships
func DestroyJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get tournament id
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, " tournaments.Destroy, string value could not be parsed: %v", err)
			return helpers.NotFound{err}
		}

		if err := tournamentmdl.Leave(c, tournamentId, u.Id); err != nil {
			log.Errorf(c, " tournamentrels.Destroy: %v", err)
			return helpers.InternalServerError{err}
		}
		// return the left tournament
		var tournament *tournamentmdl.Tournament
		if tournament, err = tournamentmdl.ById(c, tournamentId); err != nil {
			return helpers.NotFound{err}
		}

		var tJson tournamentmdl.TournamentJson
		helpers.CopyToPtrBasedStructGeneric(tournament, &tJson)
		fieldsToKeep := []string{"Id", "Name"}
		helpers.KeepFields(&tJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tJson)
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}
