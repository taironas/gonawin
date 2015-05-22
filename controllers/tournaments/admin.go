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

	"github.com/santiaago/gonawin/extract"
	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// AddAdmin let you add an admin to a tournament.
//	GET	/j/tournaments/[0-9]+/admin/add/
//
func AddAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament add admin Handler:"
	extract := extract.NewContext(c, desc, r)

	var tournament *mdl.Tournament
	var err error

	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	var userId int64
	if userId, err = extract.UserId(); err != nil {
		return err
	}

	var newAdmin *mdl.User
	if newAdmin, err = extract.Admin(userId); err != nil {
		return err
	}

	// add admin to tournament
	if err = tournament.AddAdmin(c, newAdmin.Id); err != nil {
		log.Errorf(c, "%s error on AddAdmin to tournament: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// send response
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

// RemoveAdmin handler lets you remove an admin from a tournament.
//
// Use this handler to remove a user as admin of the current tournament.
//	GET	/j/tournaments/[0-9]+/admin/remove/
//
func RemoveAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament remove admin Handler:"
	extract := extract.NewContext(c, desc, r)

	var tournament *mdl.Tournament
	var err error

	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	var userId int64
	if userId, err = extract.UserId(); err != nil {
		return err
	}

	var oldAdmin *mdl.User
	if oldAdmin, err = extract.Admin(userId); err != nil {
		return err
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

// SyncScores handler lets you syncronize scores in tournament.
//
// Use this handler to run taks to sync scores of all users in tournament.
//	GET	/j/tournaments/[0-9]+/admin/syncscores/
//
func SyncScores(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament sync scores Handler:"
	extract := extract.NewContext(c, desc, r)

	log.Infof(c, "%v", desc)

	var err error
	var tournament *mdl.Tournament

	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	// prepare data to add task to queue.
	b1, errm := json.Marshal(tournament)
	if errm != nil {
		log.Errorf(c, "%s Error marshaling", desc, errm)
	}

	task := taskqueue.NewPOSTTask("/a/sync/scores/", url.Values{
		"tournament": []string{string(b1)},
	})

	if _, err = taskqueue.Add(c, task, "gw-queue"); err != nil {
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

// ActivatePhase handler let you  activate phase of tournament.
//
// Use this handler to activate all the matches of given phase in tournament.
//	GET	/j/tournaments/[0-9]+/admin/activatephase/
//
func ActivatePhase(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament activate phase handler:"
	extract := extract.NewContext(c, desc, r)

	log.Infof(c, "%v", desc)

	var err error
	var tournament *mdl.Tournament

	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	phaseName := r.FormValue("phase")

	matches := mdl.GetMatchesByPhase(c, tournament, phaseName)

	for _, match := range matches {
		match.Ready = true
	}

	return mdl.UpdateMatches(c, matches)
}
