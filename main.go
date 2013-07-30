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

package pw

import (
	"net/http"

	pagesctrl "github.com/santiaago/purple-wing/controllers/pages"
	sessionsctrl "github.com/santiaago/purple-wing/controllers/sessions"
	usersctrl "github.com/santiaago/purple-wing/controllers/users"
	teamsctrl "github.com/santiaago/purple-wing/controllers/teams"
	"github.com/santiaago/purple-wing/helpers"
)

func init(){
	h := new(helpers.RegexpHandler)
	// usual pages
	h.HandleFunc("/", pagesctrl.TempHome)
	h.HandleFunc("/m/?", pagesctrl.Home)
	h.HandleFunc("/m/about/?", pagesctrl.About)
	h.HandleFunc("/m/contact/?", pagesctrl.Contact)
	// session
	h.HandleFunc("/m/auth/?", sessionsctrl.SessionAuth)
	h.HandleFunc("/m/oauth2callback/?", sessionsctrl.SessionAuthCallback)
	h.HandleFunc("/m/logout/?", sessionsctrl.SessionLogout)	
	// user
	h.HandleFunc("/m/users/[0-9]+/?", usersctrl.UserShow)
	h.HandleFunc("/m/users/[0-9]+/edit/?", usersctrl.UserEdit)
	// admin
	h.HandleFunc("/m/a/?", usersctrl.AdminShow)
	h.HandleFunc("/m/a/users/?", usersctrl.AdminUsers)
	// team
	h.HandleFunc("/m/teams/?", teamsctrl.TeamIndex)
	h.HandleFunc("/m/teams/new/?", teamsctrl.TeamNew)
	h.HandleFunc("/m/teams/[0-9]+/?", teamsctrl.TeamShow)
	h.HandleFunc("/m/teams/[0-9]+/edit/?", teamsctrl.TeamEdit)

	http.Handle("/", h)
}
