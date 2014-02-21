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
	"fmt"
	"net/http"

	"github.com/santiaago/purple-wing/helpers/handlers"

	invitectrl "github.com/santiaago/purple-wing/controllers/invite"
	sessionsctrl "github.com/santiaago/purple-wing/controllers/sessions"
	teamrelsctrl "github.com/santiaago/purple-wing/controllers/teamrels"
	teamsctrl "github.com/santiaago/purple-wing/controllers/teams"
	tournamentrelsctrl "github.com/santiaago/purple-wing/controllers/tournamentrels"
	tournamentsctrl "github.com/santiaago/purple-wing/controllers/tournaments"
	tournamentteamrelsctrl "github.com/santiaago/purple-wing/controllers/tournamentteamrels"
	usersctrl "github.com/santiaago/purple-wing/controllers/users"
  activitiesctrl "github.com/santiaago/purple-wing/controllers/activities"
)

// temporary main handler: for landing page
func tempHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, gonawin!")
}

// entry point of application
func init() {
	h := new(handlers.RegexpHandler)

	// temporal home page
	h.HandleFunc("/", tempHome)

	// ------------- Json Server -----------------

	// session
	h.HandleFunc("/j/auth/?", handlers.ErrorHandler(sessionsctrl.JsonAuthenticate))
	h.HandleFunc("/j/auth/twitter/?", handlers.ErrorHandler(sessionsctrl.JsonTwitterAuth))
	h.HandleFunc("/j/auth/twitter/callback/?", handlers.ErrorHandler(sessionsctrl.JsonTwitterAuthCallback))
	h.HandleFunc("/j/auth/twitter/user/?", handlers.ErrorHandler(sessionsctrl.JsonTwitterUser))
	// user
	h.HandleFunc("/j/users/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.IndexJson)))
	h.HandleFunc("/j/users/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.ShowJson)))
	h.HandleFunc("/j/users/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.UpdateJson)))
	// team
	h.HandleFunc("/j/teams/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.IndexJson)))
	h.HandleFunc("/j/teams/new/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.NewJson)))
	h.HandleFunc("/j/teams/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.ShowJson)))
	h.HandleFunc("/j/teams/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.UpdateJson)))
	h.HandleFunc("/j/teams/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.DestroyJson)))
	h.HandleFunc("/j/teams/invite/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.InviteJson)))
	h.HandleFunc("/j/teams/allow/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AllowRequestJson)))
	h.HandleFunc("/j/teams/deny/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.DenyRequestJson)))
	h.HandleFunc("/j/teams/search/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.SearchJson)))
	h.HandleFunc("/j/teams/[0-9]+/members/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.MembersJson)))
	// tournament - json
	h.HandleFunc("/j/tournaments/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.IndexJson)))
	h.HandleFunc("/j/tournaments/new/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.NewJson)))
	h.HandleFunc("/j/tournaments/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.ShowJson)))
	h.HandleFunc("/j/tournaments/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.UpdateJson)))
	h.HandleFunc("/j/tournaments/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.DestroyJson)))
	h.HandleFunc("/j/tournaments/search/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.SearchJson)))
	h.HandleFunc("/j/tournaments/candidates/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.CandidateTeamsJson)))
	h.HandleFunc("/j/tournaments/[0-9]+/participants/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.ParticipantsJson)))
	// relationships - json
	h.HandleFunc("/j/teamrels/create/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamrelsctrl.CreateJson)))
	h.HandleFunc("/j/teamrels/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamrelsctrl.DestroyJson)))
	h.HandleFunc("/j/tournamentrels/create/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentrelsctrl.CreateJson)))
	h.HandleFunc("/j/tournamentrels/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentrelsctrl.DestroyJson)))
	h.HandleFunc("/j/tournamentteamrels/create/[0-9]+/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentteamrelsctrl.CreateJson)))
	h.HandleFunc("/j/tournamentteamrels/destroy/[0-9]+/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentteamrelsctrl.DestroyJson)))
	// invite - json
	h.HandleFunc("/j/invite/?", handlers.ErrorHandler(invitectrl.InviteJson))

	// experimental: sar
	h.HandleFunc("/j/tournaments/newwc/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.NewWorldCupJson)))
	h.HandleFunc("/j/tournaments/[0-9]+/groups/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.GroupsJson)))
	h.HandleFunc("/j/tournaments/[0-9]+/calendar/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.CalendarJson)))
	h.HandleFunc("/j/tournaments/[0-9]+/matches/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.MatchesJson)))
	h.HandleFunc("/j/tournaments/[0-9]+/matches/[0-9]+/update/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.UpdateMatchResultJson)))
  
  // activities - json
  h.HandleFunc("/j/activities/?", handlers.ErrorHandler(handlers.Authorized(activitiesctrl.IndexJson)))

	http.Handle("/", h)
}
