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

package activities

import (
	"errors"
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	activitymdl "github.com/santiaago/purple-wing/models/activity"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

// json index activity handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		activities := activitymdl.FindByUser(c, u.Id)

		fieldsToKeep := []string{"Id", "Title", "Verb", "Actor", "Object", "Target"}
		activitiesJson := make([]activitymdl.ActivityJson, len(activities))
		helpers.TransformFromArrayOfPointers(&activities, &activitiesJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, activitiesJson)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}
