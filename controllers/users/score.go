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
	"net/http"

	"appengine"

	"github.com/santiaago/gonawin/extract"
	"github.com/santiaago/gonawin/helpers"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"
	mdl "github.com/santiaago/gonawin/models"
)

// User score handler, returns the JSON data of the requested user.
func Score(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	desc := "User Score Handler:"
	c := appengine.NewContext(r)
	extract := extract.NewContext(c, desc, r)

	var user *mdl.User
	var err error
	if user, err = extract.User(); err != nil {
		return err
	}

	scores := user.TournamentsScores(c)

	data := struct {
		Scores []*mdl.ScoreOverall
	}{
		scores,
	}

	return templateshlp.RenderJson(w, c, data)
}
