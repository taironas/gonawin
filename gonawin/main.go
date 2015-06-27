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

// Package gonawin provides starts a the gonawin web application.
//
package gonawin

import (
	"net/http"

	// use to do backup/restore of the production database
	_ "appengine/remote_api"

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

	checkErrors := handlers.ErrorHandler
	authorized := handlers.Authorized
	adminAuthorized := handlers.AdminAuthorized
	// ------------- Json Server -----------------

	// session
	r.HandleFunc("/j/auth", checkErrors(sessionsctrl.Authenticate))
	r.HandleFunc("/j/auth/twitter", checkErrors(sessionsctrl.TwitterAuth))
	r.HandleFunc("/j/auth/twitter/callback", checkErrors(sessionsctrl.TwitterAuthCallback))
	r.HandleFunc("/j/auth/twitter/user", checkErrors(sessionsctrl.TwitterUser))
	r.HandleFunc("/j/auth/googleloginurl", checkErrors(sessionsctrl.GoogleAccountsLoginURL))
	r.HandleFunc("/j/auth/google/callback", checkErrors(sessionsctrl.GoogleAuthCallback))
	r.HandleFunc("/j/auth/google/user", checkErrors(sessionsctrl.GoogleUser))
	r.HandleFunc("/j/auth/google/deletecookie", checkErrors(sessionsctrl.GoogleDeleteCookie))
	r.HandleFunc("/j/auth/serviceids", checkErrors(sessionsctrl.AuthServiceIds))

	// user
	r.HandleFunc("/j/users", checkErrors(adminAuthorized(usersctrl.Index)))
	r.HandleFunc("/j/users/show/:userId", checkErrors(authorized(usersctrl.Show)))
	r.HandleFunc("/j/users/update/:userId", checkErrors(authorized(usersctrl.Update)))
	r.HandleFunc("/j/users/destroy/:userId", checkErrors(authorized(usersctrl.Destroy)))
	r.HandleFunc("/j/users/:userId/scores", checkErrors(authorized(usersctrl.Score)))
	r.HandleFunc("/j/users/search", checkErrors(authorized(usersctrl.Search)))
	r.HandleFunc("/j/users/:userId/teams", checkErrors(authorized(usersctrl.Teams)))
	r.HandleFunc("/j/users/:userId/tournaments", checkErrors(authorized(usersctrl.Tournaments)))
	r.HandleFunc("/j/users/allow/:teamId", checkErrors(authorized(usersctrl.AllowInvitation)))
	r.HandleFunc("/j/users/deny/:teamId", checkErrors(authorized(usersctrl.DenyInvitation)))

	// team
	r.HandleFunc("/j/teams", checkErrors(authorized(teamsctrl.Index)))
	r.HandleFunc("/j/teams/new", checkErrors(authorized(teamsctrl.New)))
	r.HandleFunc("/j/teams/show/:teamId", checkErrors(authorized(teamsctrl.Show)))
	r.HandleFunc("/j/teams/update/:teamId", checkErrors(authorized(teamsctrl.Update)))
	r.HandleFunc("/j/teams/destroy/:teamId", checkErrors(authorized(teamsctrl.Destroy)))
	r.HandleFunc("/j/teams/requestinvite/:teamId", checkErrors(authorized(teamsctrl.RequestInvite)))
	r.HandleFunc("/j/teams/sendinvite/:teamId/:userId", checkErrors(authorized(teamsctrl.SendInvite)))
	r.HandleFunc("/j/teams/invited/:teamId", checkErrors(authorized(teamsctrl.Invited)))
	r.HandleFunc("/j/teams/allow/:requestId", checkErrors(authorized(teamsctrl.AllowRequest)))
	r.HandleFunc("/j/teams/deny/:requestId", checkErrors(authorized(teamsctrl.DenyRequest)))
	r.HandleFunc("/j/teams/search", checkErrors(authorized(teamsctrl.Search)))
	r.HandleFunc("/j/teams/:teamId/members", checkErrors(authorized(teamsctrl.Members)))
	r.HandleFunc("/j/teams/:teamId/ranking", checkErrors(authorized(teamsctrl.Ranking)))
	r.HandleFunc("/j/teams/:teamId/accuracies/:tournamentId", checkErrors(authorized(teamsctrl.AccuracyByTournament)))
	r.HandleFunc("/j/teams/:teamId/accuracies", checkErrors(authorized(teamsctrl.Accuracies)))
	r.HandleFunc("/j/teams/:teamId/prices", checkErrors(authorized(teamsctrl.Prices)))
	r.HandleFunc("/j/teams/:teamId/prices/:tournamentId", checkErrors(authorized(teamsctrl.PriceByTournament)))
	r.HandleFunc("/j/teams/:teamId/prices/update/:tournamentId", checkErrors(authorized(teamsctrl.UpdatePrice)))
	r.HandleFunc("/j/teams/:teamId/admin/add/:userId", checkErrors(authorized(teamsctrl.AddAdmin)))
	r.HandleFunc("/j/teams/:teamId/admin/remove/:userId", checkErrors(authorized(teamsctrl.RemoveAdmin)))

	// tournament
	r.HandleFunc("/j/tournaments", checkErrors(authorized(tournamentsctrl.Index)))
	r.HandleFunc("/j/tournaments/new", checkErrors(adminAuthorized(tournamentsctrl.New)))
	r.HandleFunc("/j/tournaments/show/:tournamentId", checkErrors(authorized(tournamentsctrl.Show)))
	r.HandleFunc("/j/tournaments/update/:tournamentId", checkErrors(adminAuthorized(tournamentsctrl.Update)))
	r.HandleFunc("/j/tournaments/destroy/:tournamentId", checkErrors(adminAuthorized(tournamentsctrl.Destroy)))
	r.HandleFunc("/j/tournaments/search", checkErrors(authorized(tournamentsctrl.Search)))
	r.HandleFunc("/j/tournaments/:tournamentId/candidates", checkErrors(authorized(tournamentsctrl.CandidateTeams)))
	r.HandleFunc("/j/tournaments/:tournamentId/participants", checkErrors(authorized(tournamentsctrl.Participants)))

	// relationships
	r.HandleFunc("/j/teams/join/:teamId", checkErrors(authorized(teamsctrl.Join)))
	r.HandleFunc("/j/teams/leave/:teamId", checkErrors(authorized(teamsctrl.Leave)))
	r.HandleFunc("/j/tournaments/join/:tournamentId", checkErrors(authorized(tournamentsctrl.Join)))
	r.HandleFunc("/j/tournaments/leave/:tournamentId", checkErrors(authorized(tournamentsctrl.Leave)))
	r.HandleFunc("/j/tournaments/joinasteam/:tournamentId/:teamId", checkErrors(authorized(tournamentsctrl.JoinAsTeam)))
	r.HandleFunc("/j/tournaments/leaveasteam/:tournamentId/:teamId", checkErrors(authorized(tournamentsctrl.LeaveAsTeam)))

	// invite
	r.HandleFunc("/j/invite", checkErrors(authorized(invitectrl.Invite)))

	// tournament world cup
	r.HandleFunc("/j/tournaments/newwc", checkErrors(adminAuthorized(tournamentsctrl.NewWorldCup)))
	r.HandleFunc("/j/tournaments/getwc", checkErrors(authorized(tournamentsctrl.GetWorldCup)))

	// tournament champions league
	r.HandleFunc("/j/tournaments/newcl", checkErrors(adminAuthorized(tournamentsctrl.NewChampionsLeague)))
	r.HandleFunc("/j/tournaments/getcl", checkErrors(authorized(tournamentsctrl.GetChampionsLeague)))

	// tournament champions league
	r.HandleFunc("/j/tournaments/newca", checkErrors(adminAuthorized(tournamentsctrl.NewCopaAmerica)))
	r.HandleFunc("/j/tournaments/getca", checkErrors(authorized(tournamentsctrl.GetCopaAmerica)))

	// tournament
	r.HandleFunc("/j/tournaments/:tournamentId/groups", checkErrors(authorized(tournamentsctrl.Groups)))
	r.HandleFunc("/j/tournaments/:tournamentId/calendar", checkErrors(authorized(tournamentsctrl.Calendar)))
	r.HandleFunc("/j/tournaments/:tournamentId/:teamId/calendarwithprediction", checkErrors(authorized(tournamentsctrl.CalendarWithPrediction)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches", checkErrors(authorized(tournamentsctrl.Matches)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches/:matchId/update", checkErrors(adminAuthorized(tournamentsctrl.UpdateMatchResult)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches/:matchId/predict", checkErrors(authorized(tournamentsctrl.Predict)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches/:matchId/blockprediction", checkErrors(adminAuthorized(tournamentsctrl.BlockMatchPrediction)))
	r.HandleFunc("/j/tournaments/:tournamentId/ranking", checkErrors(authorized(tournamentsctrl.Ranking)))
	r.HandleFunc("/j/tournaments/:tournamentId/teams", checkErrors(authorized(tournamentsctrl.Teams)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/reset", checkErrors(adminAuthorized(tournamentsctrl.Reset)))
	r.HandleFunc("/j/tournaments/:tournamentId/matches/simulate", checkErrors(adminAuthorized(tournamentsctrl.SimulateMatches)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/updateteam", checkErrors(adminAuthorized(tournamentsctrl.UpdateTeam)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/add/:userId", checkErrors(adminAuthorized(tournamentsctrl.AddAdmin)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/remove/:userId", checkErrors(adminAuthorized(tournamentsctrl.RemoveAdmin)))
	r.HandleFunc("/j/tournaments/:tournamentId/admin/activatephase", checkErrors(adminAuthorized(tournamentsctrl.ActivatePhase)))

	// activities
	r.HandleFunc("/j/activities", checkErrors(authorized(activitiesctrl.Index)))

	// admin handlers
	r.HandleFunc("/a/update/scores", checkErrors(tasksctrl.UpdateScores))
	r.HandleFunc("/a/update/users/scores", checkErrors(tasksctrl.UpdateUsersScores))
	r.HandleFunc("/a/publish/users/scoreactivities", checkErrors(tasksctrl.PublishUsersScoreActivities))
	r.HandleFunc("/a/create/scoreentities", checkErrors(tasksctrl.CreateScoreEntities))
	r.HandleFunc("/a/add/scoreentities/score", checkErrors(tasksctrl.AddScoreToScoreEntities))
	r.HandleFunc("/a/invite", checkErrors(tasksctrl.Invite))
	r.HandleFunc("/a/publish/users/deletepredicts", checkErrors(tasksctrl.DeleteUserPredicts))

	http.Handle("/", r)
}
