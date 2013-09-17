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
	tournamentsctrl "github.com/santiaago/purple-wing/controllers/tournaments"
	teamrelsctrl "github.com/santiaago/purple-wing/controllers/teamrels"
	"github.com/santiaago/purple-wing/helpers/handlers"
)

func init(){
	h := new(handlers.RegexpHandler)

	// usual pages
	h.HandleFunc("/", pagesctrl.TempHome)
	h.HandleFunc("/m/?", pagesctrl.Home)
	h.HandleFunc("/m/about/?", pagesctrl.About)
	h.HandleFunc("/m/contact/?", pagesctrl.Contact)

	// session
	h.HandleFunc("/m/auth/?", sessionsctrl.Authenticate)
	h.HandleFunc("/m/auth/google/?", sessionsctrl.AuthenticateWithGoogle)
	h.HandleFunc("/m/auth/google/callback/?", sessionsctrl.GoogleAuthCallback)
	h.HandleFunc("/m/auth/twitter/?", sessionsctrl.AuthenticateWithTwitter)
	h.HandleFunc("/m/auth/twitter/callback/?", sessionsctrl.TwitterAuthCallback)
	h.HandleFunc("/m/logout/?", handlers.User(sessionsctrl.SessionLogout))	

	// user
	h.HandleFunc("/m/users/[0-9]+/?", handlers.User(usersctrl.Show))
	h.HandleFunc("/m/users/[0-9]+/edit/?", handlers.User(usersctrl.Edit))

	// admin
	h.HandleFunc("/m/a/?", handlers.Admin(usersctrl.AdminShow))
	h.HandleFunc("/m/a/users/?", handlers.Admin(usersctrl.AdminUsers))
	h.HandleFunc("/m/a/search/?", handlers.Admin(usersctrl.AdminSearch))

	// team
	h.HandleFunc("/m/teams/?", handlers.User(teamsctrl.Index))
	h.HandleFunc("/m/teams/new/?", handlers.User(teamsctrl.New))
	h.HandleFunc("/m/teams/[0-9]+/?", handlers.User(teamsctrl.Show))
	h.HandleFunc("/m/teams/[0-9]+/edit/?", handlers.User(teamsctrl.Edit))

	// tournament
	h.HandleFunc("/m/tournaments/?", handlers.User(tournamentsctrl.Index))
	h.HandleFunc("/m/tournaments/new/?", handlers.User(tournamentsctrl.New))
	h.HandleFunc("/m/tournaments/[0-9]+/?", handlers.User(tournamentsctrl.Show))
	h.HandleFunc("/m/tournaments/[0-9]+/edit/?", handlers.User(tournamentsctrl.Edit))

	// team relationships
	h.HandleFunc("/m/teamrels/[0-9]+/?", handlers.User(teamrelsctrl.Show))

	http.Handle("/", h)
}
