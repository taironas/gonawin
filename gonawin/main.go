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

package gonawin

import (
	"net/http"

	"github.com/santiaago/gonawin/helpers/handlers"
	"github.com/taironas/route"

	activitiesctrl "github.com/santiaago/gonawin/controllers/activities"
	invitectrl "github.com/santiaago/gonawin/controllers/invite"
	sessionsctrl "github.com/santiaago/gonawin/controllers/sessions"
	tasksctrl "github.com/santiaago/gonawin/controllers/tasks"
	teamsctrl "github.com/santiaago/gonawin/controllers/teams"
	tournamentsctrl "github.com/santiaago/gonawin/controllers/tournaments"
	usersctrl "github.com/santiaago/gonawin/controllers/users"
)

// entry point of application
func init() {

	r := new(route.Router)

	// ------------- Json Server -----------------

	// session
	r.HandleFunc("/j/auth/?", handlers.ErrorHandler(sessionsctrl.Authenticate))
	r.HandleFunc("/j/auth/twitter/?", handlers.ErrorHandler(sessionsctrl.TwitterAuth))
	r.HandleFunc("/j/auth/twitter/callback/?", handlers.ErrorHandler(sessionsctrl.TwitterAuthCallback))
	r.HandleFunc("/j/auth/twitter/user/?", handlers.ErrorHandler(sessionsctrl.TwitterUser))
	r.HandleFunc("/j/auth/googleloginurl/?", handlers.ErrorHandler(sessionsctrl.GoogleAccountsLoginURL))
	r.HandleFunc("/j/auth/google/callback/?", handlers.ErrorHandler(sessionsctrl.GoogleAuthCallback))
	r.HandleFunc("/j/auth/google/user/?", handlers.ErrorHandler(sessionsctrl.GoogleUser))
	r.HandleFunc("/j/auth/google/deletecookie/?", handlers.ErrorHandler(sessionsctrl.GoogleDeleteCookie))
	r.HandleFunc("/j/auth/serviceids/?", handlers.ErrorHandler(sessionsctrl.AuthServiceIds))

	// user
	r.HandleFunc("/j/users/?", handlers.ErrorHandler(handlers.AdminAuthorized(usersctrl.Index)))
	r.HandleFunc("/j/users/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Show)))
	r.HandleFunc("/j/users/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Update)))
	r.HandleFunc("/j/users/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Destroy)))
	r.HandleFunc("/j/users/[0-9]+/scores/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Score)))
	r.HandleFunc("/j/users/search/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Search)))
	r.HandleFunc("/j/users/[0-9]+/teams/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Teams)))
	r.HandleFunc("/j/users/[0-9]+/tournaments/?", handlers.ErrorHandler(handlers.Authorized(usersctrl.Tournaments)))
	r.HandleFunc("/j/users/allow/[0-9]+?", handlers.ErrorHandler(handlers.Authorized(usersctrl.AllowInvitation)))
	r.HandleFunc("/j/users/deny/[0-9]+?", handlers.ErrorHandler(handlers.Authorized(usersctrl.DenyInvitation)))

	// team
	r.HandleFunc("/j/teams/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Index)))
	r.HandleFunc("/j/teams/new/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.New)))
	r.HandleFunc("/j/teams/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Show)))
	r.HandleFunc("/j/teams/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Update)))
	r.HandleFunc("/j/teams/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Destroy)))
	r.HandleFunc("/j/teams/requestinvite/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.RequestInvite)))
	r.HandleFunc("/j/teams/sendinvite/[0-9]+/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.SendInvite)))
	r.HandleFunc("/j/teams/invited/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Invited)))
	r.HandleFunc("/j/teams/allow/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AllowRequest)))
	r.HandleFunc("/j/teams/deny/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.DenyRequest)))
	r.HandleFunc("/j/teams/search/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Search)))
	r.HandleFunc("/j/teams/[0-9]+/members/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Members)))
	r.HandleFunc("/j/teams/[0-9]+/ranking/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Ranking)))
	r.HandleFunc("/j/teams/[0-9]+/accuracies/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AccuracyByTournament)))
	r.HandleFunc("/j/teams/[0-9]+/accuracies/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Accuracies)))
	r.HandleFunc("/j/teams/[0-9]+/prices/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Prices)))
	r.HandleFunc("/j/teams/[0-9]+/prices/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.PriceByTournament)))
	r.HandleFunc("/j/teams/[0-9]+/prices/update/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.UpdatePrice)))
	r.HandleFunc("/j/teams/[0-9]+/admin/add/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AddAdmin)))
	r.HandleFunc("/j/teams/[0-9]+/admin/remove/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.RemoveAdmin)))

	// tournament
	r.HandleFunc("/j/tournaments/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Index)))
	r.HandleFunc("/j/tournaments/new/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.New)))
	r.HandleFunc("/j/tournaments/show/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Show)))
	r.HandleFunc("/j/tournaments/update/[0-9]+/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Update)))
	r.HandleFunc("/j/tournaments/destroy/[0-9]+/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Destroy)))
	r.HandleFunc("/j/tournaments/search/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Search)))
	r.HandleFunc("/j/tournaments/candidates/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.CandidateTeams)))
	r.HandleFunc("/j/tournaments/[0-9]+/participants/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Participants)))

	// relationships
	r.HandleFunc("/j/teams/join/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Join)))
	r.HandleFunc("/j/teams/leave/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Leave)))
	r.HandleFunc("/j/tournaments/join/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Join)))
	r.HandleFunc("/j/tournaments/leave/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Leave)))
	r.HandleFunc("/j/tournaments/joinasteam/[0-9]+/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.JoinAsTeam)))
	r.HandleFunc("/j/tournaments/leaveasteam/[0-9]+/[0-9]+/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.LeaveAsTeam)))

	// invite
	r.HandleFunc("/j/invite/?", handlers.ErrorHandler(handlers.Authorized(invitectrl.Invite)))

	// tournament world cup
	r.HandleFunc("/j/tournaments/newwc/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.NewWorldCup)))
	r.HandleFunc("/j/tournaments/getwc/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.GetWorldCup)))
	r.HandleFunc("/j/tournaments/[0-9]+/groups/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Groups)))
	r.HandleFunc("/j/tournaments/[0-9]+/calendar/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Calendar)))
	r.HandleFunc("/j/tournaments/[0-9]+/[0-9]+/calendarwithprediction/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.CalendarWithPrediction)))
	r.HandleFunc("/j/tournaments/[0-9]+/matches/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Matches)))
	r.HandleFunc("/j/tournaments/[0-9]+/matches/[0-9]+/update/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.UpdateMatchResult)))
	r.HandleFunc("/j/tournaments/[0-9]+/matches/[0-9]+/predict/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Predict)))
	r.HandleFunc("/j/tournaments/[0-9]+/matches/[0-9]+/blockprediction/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.BlockMatchPrediction)))
	r.HandleFunc("/j/tournaments/[0-9]+/ranking/?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Ranking)))
	r.HandleFunc("/j/tournaments/[0-9]+/teams?", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Teams)))
	r.HandleFunc("/j/tournaments/[0-9]+/admin/reset/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Reset)))
	r.HandleFunc("/j/tournaments/[0-9]+/matches/simulate/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.SimulateMatches)))
	r.HandleFunc("/j/tournaments/[0-9]+/admin/updateteam/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.UpdateTeam)))
	r.HandleFunc("/j/tournaments/[0-9]+/admin/add/[0-9]+/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.AddAdmin)))
	r.HandleFunc("/j/tournaments/[0-9]+/admin/remove/[0-9]+/?", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.RemoveAdmin)))
	// activities
	r.HandleFunc("/j/activities/?", handlers.ErrorHandler(handlers.Authorized(activitiesctrl.Index)))

	// admin handlers
	r.HandleFunc("/a/update/scores/", handlers.ErrorHandler(tasksctrl.UpdateScores))
	r.HandleFunc("/a/update/users/scores/", handlers.ErrorHandler(tasksctrl.UpdateUsersScores))
	r.HandleFunc("/a/publish/users/scoreactivities/", handlers.ErrorHandler(tasksctrl.PublishUsersScoreActivities))
	r.HandleFunc("/a/publish/users/deleteactivities/", handlers.ErrorHandler(tasksctrl.DeleteUserActivities))
	r.HandleFunc("/a/create/scoreentities/", handlers.ErrorHandler(tasksctrl.CreateScoreEntities))
	r.HandleFunc("/a/add/scoreentities/score/", handlers.ErrorHandler(tasksctrl.AddScoreToScoreEntities))
	r.HandleFunc("/a/invite/", handlers.ErrorHandler(tasksctrl.Invite))

	http.Handle("/", r)
}
