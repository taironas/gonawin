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

	"appengine"

	"github.com/santiaago/gonawin/extract"
	"github.com/santiaago/gonawin/helpers"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	"github.com/santiaago/gonawin/helpers/log"
	mdl "github.com/santiaago/gonawin/models"
)

// Index activity handler, use it to get the activities of a user.
func Index(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	desc := "Index activity handler:"
	c := appengine.NewContext(r)
	extract := extract.NewContext(c, desc, r)

	count := extract.Count()
	page := extract.Page()

	activities := mdl.FindActivities(c, u, count, page)
	log.Infof(c, "%s activities = %v", desc, activities)

	json := activitiesToJSON(activities)

	return templateshlp.RenderJson(w, c, json)
}

func activitiesToJSON(activities []*mdl.Activity) []mdl.ActivityJson {
	fieldsToKeep := []string{"ID", "Type", "Verb", "Actor", "Object", "Target", "Published", "UserID"}
	json := make([]mdl.ActivityJson, len(activities))
	helpers.TransformFromArrayOfPointers(&activities, &json, fieldsToKeep)
	return json
}
