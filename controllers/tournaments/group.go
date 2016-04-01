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
	"errors"
	"net/http"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// A GroupJSON is a variable to hold a the name of a group and an array of Teams.
// We use it to group tournament teams information by group to meet world cup organization.
type GroupJSON struct {
	Name  string
	Teams []TeamJSON
}

// A TeamJSON is a variable to hold the basic information of a Team:
// The name of the team, the number of points recorded in the group phase, the goals for and against.
type TeamJSON struct {
	Name   string
	Points int64
	GoalsF int64
	GoalsA int64
	Iso    string
}

// Groups handelr sends the JSON tournament groups data.
// use this handler to get groups of a tournament.
func Groups(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Group Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var tournament *mdl.Tournament
	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	groups := mdl.Groups(c, tournament.GroupIds)
	groupsJSON := formatGroupsJSON(groups)

	data := struct {
		Groups []GroupJSON
	}{
		groupsJSON,
	}

	return templateshlp.RenderJSON(w, c, data)
}

// Format a TGroup array into a GroupJSON array.
//
func formatGroupsJSON(groups []*mdl.Tgroup) []GroupJSON {

	groupsJSON := make([]GroupJSON, len(groups))
	for i, g := range groups {
		groupsJSON[i].Name = g.Name
		teams := make([]TeamJSON, len(g.Teams))
		for j, t := range g.Teams {
			teams[j].Name = t.Name
			teams[j].Points = g.Points[j]
			teams[j].GoalsF = g.GoalsF[j]
			teams[j].GoalsA = g.GoalsA[j]
			teams[j].Iso = t.Iso
		}
		groupsJSON[i].Teams = teams
	}
	return groupsJSON
}
