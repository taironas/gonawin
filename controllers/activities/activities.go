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
	"math"
	"net/http"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// Index activity handler, use it to get the activities of a user.
// You can pass a 'count' and a 'page' param to the http.Request to
// filter the activities that you want. default values are 20 and 1
// respectively.
//
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

	lastPage := math.Ceil(float64(int64(len(activities)) / count))

	vm := buildIndexActivitiesViewModel(activities, count, page, int64(lastPage))

	return templateshlp.RenderJson(w, c, vm)
}

type indexActivitiesViewModel struct {
	Results activitiesViewModel
	Status  string
}

func buildIndexActivitiesViewModel(activities []*mdl.Activity, perPage, currentPage, lastPage int64) indexActivitiesViewModel {
	return indexActivitiesViewModel{Results: buildActivitiesViewModel(activities, perPage, currentPage, lastPage), Status: "OK"}
}

type activitiesViewModel struct {
	Total       int64
	PerPage     int64
	CurrentPage int64
	LastPage    int64
	Activities  []mdl.ActivityJson
}

func buildActivitiesViewModel(activities []*mdl.Activity, perPage, currentPage, lastPage int64) activitiesViewModel {
	return activitiesViewModel{Total: int64(len(activities)), PerPage: perPage, CurrentPage: currentPage, LastPage: lastPage, Activities: buildJSONActivities(activities)}
}

func buildJSONActivities(activities []*mdl.Activity) []mdl.ActivityJson {
	fieldsToKeep := []string{"ID", "Type", "Verb", "Actor", "Object", "Target", "Published", "UserID"}
	json := make([]mdl.ActivityJson, len(activities))
	helpers.TransformFromArrayOfPointers(&activities, &json, fieldsToKeep)
	return json
}
