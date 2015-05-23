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

package users

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

type userViewModel struct {
	Id       int64
	Username string
	Alias    string
	Score    int64
	ImageURL string
}

// Search handler returns the result of a user search in a JSON format.
// It uses parameter 'q' to make the query.
//
//	GET	/j/user/search/			Search for all users respecting the query "q"
//
func Search(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	keywords := r.FormValue("q")
	if r.Method != "GET" || len(keywords) == 0 {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User Search Handler:"

	words := helpers.SetOfStrings(keywords)

	var ids []int64
	var err error
	if ids, err = mdl.GetUserInvertedIndexes(c, words); err != nil {
		return unableToPerformSearch(c, w, err)
	}

	result := mdl.UserScore(c, keywords, ids)
	log.Infof(c, "%s result from UserScore: %v", desc, result)

	var users []*mdl.User
	if users = mdl.UsersByIds(c, result); len(users) == 0 {
		return notFound(c, w, keywords)
	}

	log.Infof(c, "%s ByIds result %v", desc, users)

	uvm := buildUserViewModel(users)

	data := struct {
		Users []userViewModel `json:",omitempty"`
	}{
		uvm,
	}
	return templateshlp.RenderJson(w, c, data)
}

func buildUserViewModel(users []*mdl.User) []userViewModel {
	uvm := make([]userViewModel, len(users))
	for i, u := range users {
		uvm[i].Id = u.Id
		uvm[i].Username = u.Username
		uvm[i].Alias = u.Alias
		uvm[i].Score = u.Score
		uvm[i].ImageURL = helpers.UserImageURL(u.Name, u.Id)
	}
	return uvm
}

func notFound(c appengine.Context, w http.ResponseWriter, keywords string) error {
	msg := fmt.Sprintf("Oops! Your search - %s - did not match any %s.", keywords, "user")
	data := struct {
		MessageInfo string `json:",omitempty"`
	}{
		msg,
	}
	return templateshlp.RenderJson(w, c, data)

}

func unableToPerformSearch(c appengine.Context, w http.ResponseWriter, err error) error {
	log.Errorf(c, "%s users.Index, error occurred when getting indexes of words: %v", desc, err)
	data := struct {
		MessageDanger string `json:",omitempty"`
	}{
		"Oops! something went wrong, we are unable to perform search query.",
	}
	return templateshlp.RenderJson(w, c, data)
}
