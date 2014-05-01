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

package gonawin

import (
	"fmt"
	"net/http"

	"github.com/santiaago/purple-wing/helpers/handlers"

	activitiesctrl "github.com/santiaago/purple-wing/controllers/activities"
	invitectrl "github.com/santiaago/purple-wing/controllers/invite"
	sessionsctrl "github.com/santiaago/purple-wing/controllers/sessions"
	tasksctrl "github.com/santiaago/purple-wing/controllers/tasks"
	teamsctrl "github.com/santiaago/purple-wing/controllers/teams"
	tournamentsctrl "github.com/santiaago/purple-wing/controllers/tournaments"
	usersctrl "github.com/santiaago/purple-wing/controllers/users"
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
	h.HandleFunc("/j/auth/?", handlers.ErrorHandler(sessionsctrl.Authenticate))
	h.HandleFunc("/j/auth/twitter/?", handlers.ErrorHandler(sessionsctrl.TwitterAuth))
	h.HandleFunc("/j/auth/twitter/callback/?", handlers.ErrorHandler(sessionsctrl.TwitterAuthCallback))
	h.HandleFunc("/j/auth/twitter/user/?", handlers.ErrorHandler(sessionsctrl.TwitterUser))
	h.HandleFunc("/j/auth/googleloginurl/?", handlers.ErrorHandler(sessionsctrl.GoogleAccountsLoginURL))
	h.HandleFunc("/j/auth/google/callback/?", handlers.ErrorHandler(sessionsctrl.GoogleAuthCallback))
	h.HandleFunc("/j/auth/google/user/?", handlers.ErrorHandler(sessionsctrl.GoogleUser))
	h.HandleFunc("/j/auth/serviceids/?", handlers.ErrorHandler(sessionsctrl.AuthServiceIds))

	// user
	h.HandleFunc("/j/users/?", handlers.ErrorHandler(handlers.AdminAuthorized(usersctrl.Index)))
	h.HandleFunc("/j/users/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Show)))
	h.HandleFunc("/j/users/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Update)))
	h.HandleFunc("/j/users/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Destroy)))
	h.HandleFunc("/j/users/[0-9]+/scores/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Score)))
	h.HandleFunc("/j/users/search/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Search)))
	h.HandleFunc("/j/users/[0-9]+/teams/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Teams)))
	h.HandleFunc("/j/users/[0-9]+/tournaments/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Tournaments)))

	// team
	h.HandleFunc("/j/teams/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Index)))
	h.HandleFunc("/j/teams/new/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.New)))
	h.HandleFunc("/j/teams/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Show)))
	h.HandleFunc("/j/teams/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Update)))
	h.HandleFunc("/j/teams/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Destroy)))
	h.HandleFunc("/j/teams/requestinvite/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.RequestInvite)))
	h.HandleFunc("/j/teams/sendinvite/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.SendInvite)))
	h.HandleFunc("/j/teams/invited/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Invited)))
	h.HandleFunc("/j/teams/allow/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AllowRequest)))
	h.HandleFunc("/j/teams/deny/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.DenyRequest)))
	h.HandleFunc("/j/teams/search/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Search)))
	h.HandleFunc("/j/teams/[0-9]+/members/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Members)))
	h.HandleFunc("/j/teams/[0-9]+/ranking/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Ranking)))
	h.HandleFunc("/j/teams/[0-9]+/accuracies/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AccuracyByTournament)))
	h.HandleFunc("/j/teams/[0-9]+/accuracies/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Accuracies)))
	h.HandleFunc("/j/teams/[0-9]+/prices/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Prices)))
	h.HandleFunc("/j/teams/[0-9]+/prices/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.PriceByTournament)))
	h.HandleFunc("/j/teams/[0-9]+/prices/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.UpdatePrice)))
	h.HandleFunc("/j/teams/[0-9]+/admin/add/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AddAdmin)))
	h.HandleFunc("/j/teams/[0-9]+/admin/remove/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.RemoveAdmin)))

	// tournament
	h.HandleFunc("/j/tournaments/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Index)))
	h.HandleFunc("/j/tournaments/new/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.New)))
	h.HandleFunc("/j/tournaments/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Show)))
	h.HandleFunc("/j/tournaments/update/[0-9]+/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Update)))
	h.HandleFunc("/j/tournaments/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Destroy)))
	h.HandleFunc("/j/tournaments/search/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Search)))
	h.HandleFunc("/j/tournaments/candidates/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.CandidateTeams)))
	h.HandleFunc("/j/tournaments/[0-9]+/participants/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Participants)))

	// relationships
	h.HandleFunc("/j/teams/join/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Join)))
	h.HandleFunc("/j/teams/leave/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Leave)))
	h.HandleFunc("/j/tournaments/join/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Join)))
	h.HandleFunc("/j/tournaments/leave/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Leave)))
	h.HandleFunc("/j/tournaments/joinasteam/[0-9]+/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.JoinAsTeam)))
	h.HandleFunc("/j/tournaments/leaveasteam/[0-9]+/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.LeaveAsTeam)))

	// invite
	h.HandleFunc("/j/invite/?", handlers.ErrorHandler(handlers.Authorized(invitectrl.Invite)))

	// tournament world cup
	h.HandleFunc("/j/tournaments/newwc/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.NewWorldCup)))
	h.HandleFunc("/j/tournaments/[0-9]+/groups/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Groups)))
	h.HandleFunc("/j/tournaments/[0-9]+/calendar/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Calendar)))
	h.HandleFunc("/j/tournaments/[0-9]+/matches/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Matches)))
	h.HandleFunc("/j/tournaments/[0-9]+/matches/[0-9]+/update/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.UpdateMatchResult)))
	h.HandleFunc("/j/tournaments/[0-9]+/matches/[0-9]+/predict/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Predict)))
	h.HandleFunc("/j/tournaments/[0-9]+/ranking/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Ranking)))
	h.HandleFunc("/j/tournaments/[0-9]+/teams?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Teams)))
	h.HandleFunc("/j/tournaments/[0-9]+/admin/reset/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Reset)))
	h.HandleFunc("/j/tournaments/[0-9]+/matches/simulate/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.SimulateMatches)))
	h.HandleFunc("/j/tournaments/[0-9]+/admin/updateteam/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.UpdateTeam)))
	h.HandleFunc("/j/tournaments/[0-9]+/admin/add/[0-9]+/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.AddAdmin)))
	h.HandleFunc("/j/tournaments/[0-9]+/admin/remove/[0-9]+/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.RemoveAdmin)))

	// activities
	h.HandleFunc("/j/activities/?", handlers.ErrorHandler(handlers.Authorized(activitiesctrl.Index)))

	// admin handlers
	h.HandleFunc("/a/update/scores/", handlers.ErrorHandler(tasksctrl.UpdateScores))
	h.HandleFunc("/a/update/users/scores/", handlers.ErrorHandler(tasksctrl.UpdateUsersScores))
	h.HandleFunc("/a/publish/users/scoreactivities/", handlers.ErrorHandler(tasksctrl.PublishUsersScoreActivities))
	h.HandleFunc("/a/publish/users/deleteactivities/", handlers.ErrorHandler(tasksctrl.DeleteUserActivities))
	h.HandleFunc("/a/create/scoreentities/", handlers.ErrorHandler(tasksctrl.CreateScoreEntities))
	h.HandleFunc("/a/add/scoreentities/score/", handlers.ErrorHandler(tasksctrl.AddScoreToScoreEntities))
	h.HandleFunc("/a/invite/", handlers.ErrorHandler(tasksctrl.Invite))

	http.Handle("/", h)
}
