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

package tournamentrels

import (
	"appengine"

	"github.com/santiaago/purple-wing/helpers/log"
	teammdl "github.com/santiaago/purple-wing/models/team"
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	tournamentrelmdl "github.com/santiaago/purple-wing/models/tournamentrel"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

// from a tournamentId returns an array of users the participate in it.
func Participants(c appengine.Context, tournamentId int64) []*usermdl.User {
	var users []*usermdl.User

	tournamentRels := tournamentrelmdl.Find(c, "TournamentId", tournamentId)

	for _, tournamentRel := range tournamentRels {
		user, err := usermdl.ById(c, tournamentRel.UserId)
		if err != nil {
			log.Errorf(c, " Participants, cannot find user with ID=%", tournamentRel.UserId)
		} else {
			users = append(users, user)
		}
	}

	return users
}

// from a tournamentid returns an array of teams involved in tournament
func Teams(c appengine.Context, tournamentId int64) []*teammdl.Team {

	var teams []*teammdl.Team

	tournamentteamRels := tournamentteamrelmdl.Find(c, "TournamentId", tournamentId)

	for _, tournamentteamRel := range tournamentteamRels {
		team, err := teammdl.ById(c, tournamentteamRel.TeamId)
		if err != nil {
			log.Errorf(c, " Teams, cannot find team with ID=%", tournamentteamRel.TeamId)
		} else {
			teams = append(teams, team)
		}
	}

	return teams
}

// from a user id return an array of tournament the user is involved in.
func Tournaments(c appengine.Context, userId int64) []*tournamentmdl.Tournament {

	var tournaments []*tournamentmdl.Tournament

	tournamentRels := tournamentrelmdl.Find(c, "UserId", userId)

	for _, tournamentRel := range tournamentRels {
		tournament, err := tournamentmdl.ById(c, tournamentRel.TournamentId)
		if err != nil {
			log.Errorf(c, " Tournaments, cannot find team with ID=%", tournamentRel.TournamentId)
		} else {
			tournaments = append(tournaments, tournament)
		}
	}

	return tournaments
}
