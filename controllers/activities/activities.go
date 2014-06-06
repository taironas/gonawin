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

// Package activities provides the JSON handlers to get gonawin activities.
package activities

import (
	"errors"
	"net/http"
	"strconv"

	"appengine"

	"github.com/santiaago/gonawin/helpers"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	"github.com/santiaago/gonawin/helpers/log"
	mdl "github.com/santiaago/gonawin/models"
)

// Index activity handler.
func Index(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Index activity handler:"
	if r.Method == "GET" {
		count, err := strconv.ParseInt(r.FormValue("count"), 0, 64)
		if err != nil {
			log.Errorf(c, "%s: error during conversion of count parameter: %v", desc, err)
			count = 20 // set count to default value
		}
		page, err := strconv.ParseInt(r.FormValue("page"), 0, 64)
		if err != nil {
			log.Errorf(c, "%s error during conversion of page parameter: %v", desc, err)
			page = 1
		}
		// fetch user activities
		activities := mdl.FindActivities(c, u, count, page)
		log.Infof(c, "%s activities = %v", desc, activities)

		fieldsToKeep := []string{"ID", "Type", "Verb", "Actor", "Object", "Target", "Published", "UserID"}
		activitiesJson := make([]mdl.ActivityJson, len(activities))
		helpers.TransformFromArrayOfPointers(&activities, &activitiesJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, activitiesJson)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
