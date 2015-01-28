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

	"github.com/taironas/route"

	"github.com/santiaago/gonawin/helpers/handlers"

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
	r.HandleFunc("/j/auth", handlers.ErrorHandler(sessionsctrl.Authenticate))
	r.HandleFunc("/j/auth/twitter", handlers.ErrorHandler(sessionsctrl.TwitterAuth))
	r.HandleFunc("/j/auth/twitter/callback", handlers.ErrorHandler(sessionsctrl.TwitterAuthCallback))
	r.HandleFunc("/j/auth/twitter/user", handlers.ErrorHandler(sessionsctrl.TwitterUser))
	r.HandleFunc("/j/auth/googleloginurl", handlers.ErrorHandler(sessionsctrl.GoogleAccountsLoginURL))
	r.HandleFunc("/j/auth/google/callback", handlers.ErrorHandler(sessionsctrl.GoogleAuthCallback))
	r.HandleFunc("/j/auth/google/user", handlers.ErrorHandler(sessionsctrl.GoogleUser))
	r.HandleFunc("/j/auth/google/deletecookie", handlers.ErrorHandler(sessionsctrl.GoogleDeleteCookie))
	r.HandleFunc("/j/auth/serviceids", handlers.ErrorHandler(sessionsctrl.AuthServiceIds))

	// user
	r.HandleFunc("/j/users", handlers.ErrorHandler(handlers.AdminAuthorized(usersctrl.Index)))
	r.HandleFunc("/j/users/show/:userId", handlers.ErrorHandler(handlers.Authorized(usersctrl.Show)))
	r.HandleFunc("/j/users/update/:userId", handlers.ErrorHandler(handlers.Authorized(usersctrl.Update)))
	r.HandleFunc("/j/users/destroy/:userId", handlers.ErrorHandler(handlers.Authorized(usersctrl.Destroy)))
	r.HandleFunc("/j/users/:userId/scores", handlers.ErrorHandler(handlers.Authorized(usersctrl.Score)))
	r.HandleFunc("/j/users/search", handlers.ErrorHandler(handlers.Authorized(usersctrl.Search)))
	r.HandleFunc("/j/users/:userId/teams", handlers.ErrorHandler(handlers.Authorized(usersctrl.Teams)))
	r.HandleFunc("/j/users/:userId/tournaments", handlers.ErrorHandler(handlers.Authorized(usersctrl.Tournaments)))
	r.HandleFunc("/j/users/allow/:teamId", handlers.ErrorHandler(handlers.Authorized(usersctrl.AllowInvitation)))
	r.HandleFunc("/j/users/deny/:teamId", handlers.ErrorHandler(handlers.Authorized(usersctrl.DenyInvitation)))

	// team
	r.HandleFunc("/j/teams", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Index)))
	r.HandleFunc("/j/teams/new", handlers.ErrorHandler(handlers.Authorized(teamsctrl.New)))
	r.HandleFunc("/j/teams/show/:teamId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Show)))
	r.HandleFunc("/j/teams/update/:teamId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Update)))
	r.HandleFunc("/j/teams/destroy/:teamId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Destroy)))
	r.HandleFunc("/j/teams/requestinvite/:teamId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.RequestInvite)))
	r.HandleFunc("/j/teams/sendinvite/:teamId/:userId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.SendInvite)))
	r.HandleFunc("/j/teams/invited/:teamId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Invited)))
	r.HandleFunc("/j/teams/allow/:requestId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AllowRequest)))
	r.HandleFunc("/j/teams/deny/:requestId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.DenyRequest)))
	r.HandleFunc("/j/teams/search", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Search)))
	r.HandleFunc("/j/teams/:teamId/members", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Members)))
	r.HandleFunc("/j/teams/:teamId/ranking", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Ranking)))
	r.HandleFunc("/j/teams/:teamId/accuracies/:tournamentId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AccuracyByTournament)))
	r.HandleFunc("/j/teams/:teamId/accuracies", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Accuracies)))
	r.HandleFunc("/j/teams/:teamId/prices", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Prices)))
	r.HandleFunc("/j/teams/:teamId/prices/:tournamentId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.PriceByTournament)))
	r.HandleFunc("/j/teams/:teamId/prices/update/:tournamentId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.UpdatePrice)))
	r.HandleFunc("/j/teams/:teamId/admin/add/:userId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.AddAdmin)))
	r.HandleFunc("/j/teams/:teamId/admin/remove/:userId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.RemoveAdmin)))

	// tournament
	r.HandleFunc("/j/tournaments", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Index)))
	r.HandleFunc("/j/tournaments/new", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.New)))
	r.HandleFunc("/j/tournaments/show/:tournamentId", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Show)))
	r.HandleFunc("/j/tournaments/update/:tournamentId", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Update)))
	r.HandleFunc("/j/tournaments/destroy/:tournamentId", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Destroy)))
	r.HandleFunc("/j/tournaments/search", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Search)))
	r.HandleFunc("/j/tournaments/:tournamentId/candidates", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.CandidateTeams)))
	r.HandleFunc("/j/tournaments/:tournamentId/participants", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Participants)))

	// relationships
	r.HandleFunc("/j/teams/join/:teamId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Join)))
	r.HandleFunc("/j/teams/leave/:teamId", handlers.ErrorHandler(handlers.Authorized(teamsctrl.Leave)))
	r.HandleFunc("/j/tournaments/join/:tournamentId", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Join)))
	r.HandleFunc("/j/tournaments/leave/:tournamentId", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Leave)))
	r.HandleFunc("/j/tournaments/joinasteam/:tournamentId/:teamId", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.JoinAsTeam)))
	r.HandleFunc("/j/tournaments/leaveasteam/:tournamentId/:teamId", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.LeaveAsTeam)))

	// invite
	r.HandleFunc("/j/invite", handlers.ErrorHandler(handlers.Authorized(invitectrl.Invite)))

	// tournament world cup
	r.HandleFunc("/j/tournaments/newwc", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.NewWorldCup)))
	r.HandleFunc("/j/tournaments/getwc", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.GetWorldCup)))
	r.HandleFunc("/j/tournaments/:tournamentId/groups", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Groups)))
	r.HandleFunc("/j/tournaments/:tournamentId/calendar", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Calendar)))
	r.HandleFunc("/j/tournaments/:tournamentId/:teamId/calendarwithprediction", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.CalendarWithPrediction)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Matches)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches/:matchId/update", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.UpdateMatchResult)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches/:matchId/predict", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Predict)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches/:matchId/blockprediction", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.BlockMatchPrediction)))
	r.HandleFunc("/j/tournaments/:tournamentId/ranking", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Ranking)))
	r.HandleFunc("/j/tournaments/:tournamentId/teams", handlers.ErrorHandler(handlers.Authorized(tournamentsctrl.Teams)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/reset", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.Reset)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches/simulate", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.SimulateMatches)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/updateteam", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.UpdateTeam)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/add/:userId", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.AddAdmin)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/remove/:userId", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.RemoveAdmin)))
	h.HandleFunc("/j/tournaments/:tournamentId/admin/syncscores", handlers.ErrorHandler(handlers.AdminAuthorized(tournamentsctrl.SyncScores)))

	// activities
	r.HandleFunc("/j/activities", handlers.ErrorHandler(handlers.Authorized(activitiesctrl.Index)))

	// admin handlers
	r.HandleFunc("/a/update/scores", handlers.ErrorHandler(tasksctrl.UpdateScores))
	r.HandleFunc("/a/update/users/scores", handlers.ErrorHandler(tasksctrl.UpdateUsersScores))
	r.HandleFunc("/a/publish/users/scoreactivities", handlers.ErrorHandler(tasksctrl.PublishUsersScoreActivities))
	r.HandleFunc("/a/publish/users/deleteactivities", handlers.ErrorHandler(tasksctrl.DeleteUserActivities))
	r.HandleFunc("/a/create/scoreentities", handlers.ErrorHandler(tasksctrl.CreateScoreEntities))
	r.HandleFunc("/a/add/scoreentities/score", handlers.ErrorHandler(tasksctrl.AddScoreToScoreEntities))
	r.HandleFunc("/a/invite", handlers.ErrorHandler(tasksctrl.Invite))
	h.HandleFunc("/a/sync/scores/", handlers.ErrorHandler(tasksctrl.SyncScores))

	http.Handle("/", r)
}
