/*
 * Copyright (c) 2016 Santiago Arias | Remy Jourde
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
	"time"

	"appengine"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// NewEuro creates a euro tournament.
//
func NewEuro(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "New Euro Handler:"

	tournament, err := mdl.CreateEuro2016(c, u.Id)
	if err != nil {
		log.Errorf(c, "%s error when trying to create a tournament: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentCannotCreate)}
	}

	return templateshlp.RenderJSON(w, c, tournament)
}

// GetEuro returns the json data of the Euro tournament.
//
func GetEuro(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Get Euro Handler:"

	tournaments := mdl.FindTournaments(c, "Name", "2016 UEFA Euro")
	if tournaments == nil {
		log.Errorf(c, "%s Euro tournament was not found.", desc)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}

	tournament := tournaments[0]

	// tournament
	fieldsToKeep := []string{"Id", "Name", "Description"}
	var TournamentJSON mdl.TournamentJSON
	helpers.InitPointerStructure(tournament, &TournamentJSON, fieldsToKeep)
	// formatted start and end
	const layout = "2 January 2006"
	start := tournament.Start.Format(layout)
	end := tournament.End.Format(layout)
	// remaining days
	remainingDays := int64(tournament.Start.Sub(time.Now()).Hours() / 24)
	// data
	data := struct {
		Tournament    mdl.TournamentJSON
		Start         string
		End           string
		RemainingDays int64
	}{
		TournamentJSON,
		start,
		end,
		remainingDays,
	}

	return templateshlp.RenderJSON(w, c, data)

}
