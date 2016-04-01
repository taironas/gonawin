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

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// Search handler returns the result of a team search in a JSON format.
// It uses parameter 'q' to make the query.
//
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

	var ids []int64
	var err error
	if ids, err = mdl.GetTeamInvertedIndexes(c, words); err != nil {
		return unableToPerformSearch(c, w, desc, err)
	}

	result := mdl.TeamScore(c, keywords, ids)
	log.Infof(c, "%s result from TeamScore: %v", desc, result)

	var teams []*mdl.Team
	if teams, err = mdl.TeamsByIDs(c, result); err != nil {
		log.Infof(c, "%v something failed when calling TeamsByIDs: %v", desc, err)
		return notFound(c, w, keywords)
	}

	if len(teams) == 0 {
		return notFound(c, w, keywords)
	}

	log.Infof(c, "%s ByIds result %v", desc, teams)

	svm := buildTeamSearchViewModel(teams)

	return templateshlp.RenderJSON(w, c, svm)
}

type teamSearchViewModel struct {
	Teams []teamSearchTeamViewModel
}

type teamSearchTeamViewModel struct {
	ID           int64 `json:"Id"`
	Name         string
	AdminIds     []int64
	Private      bool
	Accuracy     float64
	MembersCount int64
	ImageURL     string
}

func buildTeamSearchViewModel(teams []*mdl.Team) teamSearchViewModel {
	tvm := make([]teamSearchTeamViewModel, len(teams))
	for i, t := range teams {
		tvm[i].ID = t.ID
		tvm[i].Name = t.Name
		tvm[i].AdminIds = t.AdminIDs
		tvm[i].Private = t.Private
		tvm[i].Accuracy = t.Accuracy
		tvm[i].MembersCount = t.MembersCount
		tvm[i].ImageURL = helpers.TeamImageURL(t.Name, t.ID)
	}
	return teamSearchViewModel{Teams: tvm}
}

func notFound(c appengine.Context, w http.ResponseWriter, keywords string) error {
	msg := fmt.Sprintf("Oops! Your search - %s - did not match any %s.", keywords, "team")
	data := struct {
		MessageInfo string `json:",omitempty"`
	}{
		msg,
	}
	return templateshlp.RenderJSON(w, c, data)
}

func unableToPerformSearch(c appengine.Context, w http.ResponseWriter, desc string, err error) error {
	log.Errorf(c, "%s teams.Index, error occurred when getting indexes of words: %v", desc, err)
	data := struct {
		MessageDanger string `json:",omitempty"`
	}{
		"Oops! something went wrong, we are unable to perform search query.",
	}
	return templateshlp.RenderJSON(w, c, data)
}
