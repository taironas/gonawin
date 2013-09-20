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
 
package teams

import (
	"net/http"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
	teammdl "github.com/santiaago/purple-wing/models/team"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
)

func Players(r *http.Request, teamId int64) []*usermdl.User {
	var users []*usermdl.User
	
	teamRels := teamrelmdl.Find(r, "TeamId", teamId)

	for _, teamRel := range teamRels {
		user, _ := usermdl.ById(r, teamRel.UserId)

		users = append(users, user)
	}

	return users
}

func Teams(r *http.Request, userId int64) []*teammdl.Team {
	var teams []*teammdl.Team
	
	teamRels := teamrelmdl.Find(r, "UserId", userId)

	for _, teamRel := range teamRels {
		team, _ := teammdl.ById(r, teamRel.TeamId)

		teams = append(teams, team)
	}

	return teams
}