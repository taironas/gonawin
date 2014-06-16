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

package tournaments

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"appengine"
	"appengine/taskqueue"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/handlers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Tournament add admin handler:
//
// Use this handler to add a user as admin of current tournament.
//	GET	/j/tournaments/[0-9]+/admin/add/
//
func AddAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament add admin Handler:"
	if r.Method == "POST" {
		// get tournament id and user id
		tournamentId, err1 := handlers.PermalinkID(r, c, 3)
		userId, err2 := handlers.PermalinkID(r, c, 6)
		if err1 != nil || err2 != nil {
			log.Errorf(c, "%s string value could not be parsed: %v, %v", desc, err1, err2)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tournament *mdl.Tournament
		if tournament, err1 = mdl.TournamentById(c, tournamentId); err1 != nil {
			log.Errorf(c, "%s tournament not found: %v", desc, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var newAdmin *mdl.User
		newAdmin, err := mdl.UserById(c, userId)
		log.Infof(c, "%s User: %v", desc, newAdmin)
		if err != nil {
			log.Errorf(c, "%s user not found", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		if err = tournament.AddAdmin(c, newAdmin.Id); err != nil {
			log.Errorf(c, "%s error on AddAdmin to tournament: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TournamentJson
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		msg := fmt.Sprintf("You added %s as admin of tournament %s.", newAdmin.Name, tournament.Name)
		data := struct {
			MessageInfo string `json:",omitempty"`
			Tournament  mdl.TournamentJson
		}{
			msg,
			tJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Tournament remove admin handler:
//
// Use this handler to remove a user as admin of the current tournament.
//	GET	/j/tournaments/[0-9]+/admin/remove/
//
func RemoveAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament remove admin Handler:"

	if r.Method == "POST" {
		// get tournament id and user id
		tournamentId, err1 := handlers.PermalinkID(r, c, 3)
		userId, err2 := handlers.PermalinkID(r, c, 6)
		if err1 != nil || err2 != nil {
			log.Errorf(c, "%s string value could not be parsed: %v, %v.", desc, err1, err2)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tournament *mdl.Tournament
		if tournament, err1 = mdl.TournamentById(c, tournamentId); err1 != nil {
			log.Errorf(c, "%s tournament not found: %v.", desc, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var oldAdmin *mdl.User
		oldAdmin, err := mdl.UserById(c, userId)
		log.Infof(c, "%s User: %v.", desc, oldAdmin)
		if err != nil {
			log.Errorf(c, "%s user not found.", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		if err = tournament.RemoveAdmin(c, oldAdmin.Id); err != nil {
			log.Errorf(c, "%s error on RemoveAdmin to tournament: %v.", desc, err)
			return &helpers.InternalServerError{Err: err}
		}

		var tJson mdl.TournamentJson
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		msg := fmt.Sprintf("You removed %s as admin of tournament %s.", oldAdmin.Name, tournament.Name)
		data := struct {
			MessageInfo string `json:",omitempty"`
			Tournament  mdl.TournamentJson
		}{
			msg,
			tJson,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return nil
}

// Tournament sync scores handler:
//
// Use this handler to run taks to sync scores of all users in tournament.
//	GET	/j/tournaments/[0-9]+/admin/syncscores/
//
func SyncScores(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	c := appengine.NewContext(r)
	desc := "Tournament sync scores Handler:"
	log.Infof(c, "%v", desc)
	if r.Method == "POST" {
		// get tournament id and user id
		tournamentId, err1 := handlers.PermalinkID(r, c, 3)
		if err1 != nil {
			log.Errorf(c, "%s string value could not be parsed: %v", desc, err1)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tournament *mdl.Tournament
		if tournament, err1 = mdl.TournamentById(c, tournamentId); err1 != nil {
			log.Errorf(c, "%s tournament not found: %v", desc, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		// prepare data to add task to queue.
		b1, errm := json.Marshal(tournament)
		if errm != nil {
			log.Errorf(c, "%s Error marshaling", desc, errm)
		}

		task := taskqueue.NewPOSTTask("/a/sync/scores/", url.Values{
			"tournament": []string{string(b1)},
		})

		if _, err := taskqueue.Add(c, task, "gw-queue"); err != nil {
			log.Errorf(c, "%s unable to add task to taskqueue.", desc)
			return err
		} else {
			log.Infof(c, "%s add task to taskqueue successfully", desc)
		}

		msg := fmt.Sprintf("You send task to synch scores for all users.")
		data := struct {
			MessageInfo string `json:",omitempty"`
		}{
			msg,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
