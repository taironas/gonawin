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

// User search handler.
// Use this handler to search for a user.
//	GET	/j/user/search/			Search for all users respecting the query "q"
//
func Search(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "User Search Handler:"
	keywords := r.FormValue("q")
	if r.Method == "GET" && (len(keywords) > 0) {

		words := helpers.SetOfStrings(keywords)
		ids, err := mdl.GetUserInvertedIndexes(c, words)
		if err != nil {
			log.Errorf(c, "%s users.Index, error occurred when getting indexes of words: %v", desc, err)
			data := struct {
				MessageDanger string `json:",omitempty"`
			}{
				"Oops! something went wrong, we are unable to perform search query.",
			}
			return templateshlp.RenderJson(w, c, data)
		}
		result := mdl.UserScore(c, keywords, ids)
		log.Infof(c, "%s result from UserScore: %v", desc, result)
		users := mdl.UsersByIds(c, result)
		log.Infof(c, "%s ByIds result %v", desc, users)
		if len(users) == 0 {
			msg := fmt.Sprintf("Oops! Your search - %s - did not match any %s.", keywords, "user")
			data := struct {
				MessageInfo string `json:",omitempty"`
			}{
				msg,
			}

			return templateshlp.RenderJson(w, c, data)
		}
		// filter team information to return in json api
		fieldsToKeep := []string{"Id", "Name"}
		usersJson := make([]mdl.UserJson, len(users))
		helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)
		data := struct {
			Users []mdl.UserJson `json:",omitempty"`
		}{
			usersJson,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
