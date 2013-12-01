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
	"net/http"
	
	"appengine"

	"github.com/santiaago/purple-wing/helpers/log"	
	usermdl "github.com/santiaago/purple-wing/models/user"
	teammdl "github.com/santiaago/purple-wing/models/team"
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	tournamentrelmdl "github.com/santiaago/purple-wing/models/tournamentrel"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
)

func Participants(r *http.Request, tournamentId int64) []*usermdl.User {
	c := appengine.NewContext(r)
	
	var users []*usermdl.User
	
	tournamentRels := tournamentrelmdl.Find(r, "TournamentId", tournamentId)

	for _, tournamentRel := range tournamentRels {
		user, err := usermdl.ById(r, tournamentRel.UserId)
		if err != nil {
			log.Errorf(c, " Participants, cannot find user with ID=%", tournamentRel.UserId)
		} else {
			users = append(users, user)
		}
	}

	return users
}

func Teams(r *http.Request, tournamentId int64) []*teammdl.Team {
	c := appengine.NewContext(r)
	
	var teams []*teammdl.Team
	
	tournamentteamRels := tournamentteamrelmdl.Find(r, "TournamentId", tournamentId)

	for _, tournamentteamRel := range tournamentteamRels {
		team, err := teammdl.ById(r, tournamentteamRel.TeamId)
		if err != nil {
			log.Errorf(c, " Teams, cannot find team with ID=%", tournamentteamRel.TeamId)
		} else {
			teams = append(teams, team)
		}
	}

	return teams
}

func Tournaments(r *http.Request, userId int64) []*tournamentmdl.Tournament {
	c := appengine.NewContext(r)
	
	var tournaments []*tournamentmdl.Tournament
	
	tournamentRels := tournamentrelmdl.Find(r, "UserId", userId)

	for _, tournamentRel := range tournamentRels {
		tournament, err := tournamentmdl.ById(r, tournamentRel.TournamentId)
		if err != nil {
			log.Errorf(c, " Tournaments, cannot find team with ID=%", tournamentRel.TournamentId)
		} else {
			tournaments = append(tournaments, tournament)
		}
	}

	return tournaments
}
