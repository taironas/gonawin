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

// Package activities provides the JSON handlers to get gonawin activities.
package activities

import (
	"errors"
	"net/http"
	"sort"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	mdl "github.com/santiaago/purple-wing/models"
)

// json index activity handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		// fetch user activities
		activities := u.Activities(c)
		// fetch activities for user tournaments
		for _, tournament := range u.Tournaments(c) {
			activities = append(activities, tournament.Activities(c)...)
			// sort current slice by date of publication
			sort.Sort(activities)
		}
		// fetch activities for user teams
		for _, team := range u.Teams(c) {
			activities = append(activities, team.Activities(c)...)
			// sort current slice by date of publication
			sort.Sort(activities)
		}

		fieldsToKeep := []string{"ID", "Type", "Verb", "Actor", "Object", "Target", "Published", "UserID"}
		activitiesJson := make([]mdl.ActivityJson, len(activities))
		helpers.TransformFromArrayOfPointers(&activities, &activitiesJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, activitiesJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
