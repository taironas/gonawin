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
	//teammdl "github.com/santiaago/purple-wing/models/team"
	//tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	tournamentrelmdl "github.com/santiaago/purple-wing/models/tournamentrel"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
	// mdl "github.com/santiaago/purple-wing/models/user"
	mdl "github.com/santiaago/purple-wing/models"
)

// from a tournamentId returns an array of users that participate in it.
func Participants(c appengine.Context, tournamentId int64) []*mdl.User {
	var users []*mdl.User

	tournamentRels := tournamentrelmdl.Find(c, "TournamentId", tournamentId)

	for _, tournamentRel := range tournamentRels {
		user, err := mdl.UserById(c, tournamentRel.UserId)
		if err != nil {
			log.Errorf(c, " Participants, cannot find user with ID=%", tournamentRel.UserId)
		} else {
			users = append(users, user)
		}
	}

	return users
}

// from a tournamentid returns an array of teams involved in tournament
func Teams(c appengine.Context, tournamentId int64) []*mdl.Team {

	var teams []*mdl.Team

	tournamentteamRels := tournamentteamrelmdl.Find(c, "TournamentId", tournamentId)

	for _, tournamentteamRel := range tournamentteamRels {
		team, err := mdl.TeamById(c, tournamentteamRel.TeamId)
		if err != nil {
			log.Errorf(c, " Teams, cannot find team with ID=%", tournamentteamRel.TeamId)
		} else {
			teams = append(teams, team)
		}
	}

	return teams
}

// from a user id return an array of tournament the user is involved in.
func Tournaments(c appengine.Context, userId int64) []*mdl.Tournament {

	var tournaments []*mdl.Tournament

	tournamentRels := tournamentrelmdl.Find(c, "UserId", userId)

	for _, tournamentRel := range tournamentRels {
		tournament, err := mdl.TournamentById(c, tournamentRel.TournamentId)
		if err != nil {
			log.Errorf(c, " Tournaments, cannot find team with ID=%", tournamentRel.TournamentId)
		} else {
			tournaments = append(tournaments, tournament)
		}
	}

	return tournaments
}
