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

package teams

import (
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers/log"

	mdl "github.com/santiaago/purple-wing/models"
)

// Team add admin handler:
//
// Use this handler to add a user as admin of current team.
//	GET	/j/teams/[0-9]+/admin/add/
//
func AddAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team add admin Handler:"
	log.Infof(c, "%s start", desc)
	return nil
}

// Team remove admin handler:
//
// Use this handler to remove a user as admin of the current team.
//	GET	/j/teams/[0-9]+/admin/remove/
//
func RemoveAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team remove admin Handler:"
	log.Infof(c, "%s start", desc)
	return nil
}
