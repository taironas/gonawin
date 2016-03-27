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
	"fmt"
	"net/http"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// AddAdmin handler, use it to add an admin to a team.
//
//	GET	/j/teams/:teamId/admin/add/:userId
//
func AddAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team add admin Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	if team, err = extract.Team(); err != nil {
		log.Errorf(c, "%s error on AddAdmin to team: extract.Team() %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	var newAdmin *mdl.User
	if newAdmin, err = extract.User(); err != nil {
		log.Errorf(c, "%s error on AddAdmin to team: extract.User() %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	if err = team.AddAdmin(c, newAdmin.ID); err != nil {
		log.Errorf(c, "%s error on AddAdmin to team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	vm := buildTeamAddAdminViewModel(team, newAdmin)
	return templateshlp.RenderJson(w, c, vm)
}

type teamAddAdminViewModel struct {
	MessageInfo string `json:",omitempty"`
	Team        mdl.TeamJSON
}

func buildTeamAddAdminViewModel(team *mdl.Team, newAdmin *mdl.User) teamAddAdminViewModel {

	var t mdl.TeamJSON
	fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &t, fieldsToKeep)

	msg := fmt.Sprintf("You added %s as admin of team %s.", newAdmin.Name, team.Name)
	return teamAddAdminViewModel{msg, t}
}

// RemoveAdmin handler, use it to remove an admin from a team.
//
// Use this handler to remove a user as admin of the current team.
//	GET	/j/teams/:teamId/admin/remove/:userId
//
func RemoveAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team remove admin Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	if team, err = extract.Team(); err != nil {
		log.Errorf(c, "%s error on RemoveAdmin to team: extract.Team() %v.", desc, err)
		return &helpers.InternalServerError{Err: err}
	}

	var oldAdmin *mdl.User
	if oldAdmin, err = extract.User(); err != nil {
		log.Errorf(c, "%s error on RemoveAdmin to team: extract.User() %v.", desc, err)
		return &helpers.InternalServerError{Err: err}
	}

	if err = team.RemoveAdmin(c, oldAdmin.ID); err != nil {
		log.Errorf(c, "%s error on RemoveAdmin to team: %v.", desc, err)
		return &helpers.InternalServerError{Err: err}
	}

	vm := buildTeamRemoveAdminViewModel(team, oldAdmin)
	return templateshlp.RenderJson(w, c, vm)
}

type teamRemoveAdminViewModel struct {
	MessageInfo string `json:",omitempty"`
	Team        mdl.TeamJSON
}

func buildTeamRemoveAdminViewModel(team *mdl.Team, oldAdmin *mdl.User) teamRemoveAdminViewModel {
	var t mdl.TeamJSON
	fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &t, fieldsToKeep)

	msg := fmt.Sprintf("You removed %s as admin of team %s.", oldAdmin.Name, team.Name)

	return teamRemoveAdminViewModel{msg, t}

}
