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

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Search handler, use it to get all teams that match the search.
//	GET	/j/teams/search/			Search for all teams respecting the query "q"
//
func Search(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	keywords := r.FormValue("q")
	if r.Method != "GET" || len(keywords) == 0 {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Search Handler:"

	words := helpers.SetOfStrings(keywords)
	ids, err := mdl.GetTeamInvertedIndexes(c, words)
	if err != nil {
		log.Errorf(c, "%s teams.Index, error occurred when getting indexes of words: %v", desc, err)
		data := struct {
			MessageDanger string `json:",omitempty"`
		}{
			"Oops! something went wrong, we are unable to perform search query.",
		}
		return templateshlp.RenderJson(w, c, data)
	}
	result := mdl.TeamScore(c, keywords, ids)
	log.Infof(c, "%s result from TeamScore: %v", desc, result)
	teams := mdl.TeamsByIds(c, result)
	log.Infof(c, "%s ByIds result %v", desc, teams)
	if len(teams) == 0 {
		msg := fmt.Sprintf("Oops! Your search - %s - did not match any %s.", keywords, "team")
		data := struct {
			MessageInfo string `json:",omitempty"`
		}{
			msg,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	// filter team information to return in json api
	type team struct {
		Id           int64
		Name         string
		AdminIds     []int64
		Private      bool
		Accuracy     float64
		MembersCount int64
		ImageURL     string
	}

	ts := make([]team, len(teams))
	for i, t := range teams {
		ts[i].Id = t.Id
		ts[i].Name = t.Name
		ts[i].AdminIds = t.AdminIds
		ts[i].Private = t.Private
		ts[i].Accuracy = t.Accuracy
		ts[i].MembersCount = t.MembersCount
		ts[i].ImageURL = helpers.TeamImageURL(t.Name, t.Id)
	}

	// we should not directly return an array. so we add an extra layer.
	data := struct {
		Teams []team `json:",omitempty"`
	}{
		ts,
	}
	return templateshlp.RenderJson(w, c, data)
}
