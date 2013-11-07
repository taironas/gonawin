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
 
package teamrels

import (
	"net/http"
	
	"appengine"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
	teammdl "github.com/santiaago/purple-wing/models/team"
	teamrequestmdl "github.com/santiaago/purple-wing/models/teamrequest"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
)

func Players(r *http.Request, teamId int64) []*usermdl.User {
	c := appengine.NewContext(r)
	
	var users []*usermdl.User
	
	teamRels := teamrelmdl.Find(r, "TeamId", teamId)

	for _, teamRel := range teamRels {
		user, err := usermdl.ById(r, teamRel.UserId)
		if err != nil {
			c.Errorf("pw: Players, cannot find user with ID=%", teamRel.UserId)
		} else {
			users = append(users, user)
		}
	}

	return users
}

func Teams(r *http.Request, userId int64) []*teammdl.Team {
	c := appengine.NewContext(r)
	
	var teams []*teammdl.Team
	
	teamRels := teamrelmdl.Find(r, "UserId", userId)

	for _, teamRel := range teamRels {
		team, err := teammdl.ById(r, teamRel.TeamId)
		if err != nil {
			c.Errorf("pw: Teams, cannot find team with ID=%", teamRel.TeamId)
		} else {
			teams = append(teams, team)
		}
	}

	return teams
}

func TeamsRequests(r *http.Request, teams []*teammdl.Team) []*teamrequestmdl.TeamRequest {
	var teamRequests []*teamrequestmdl.TeamRequest
	
	for _, team := range teams {
		teamRequests = append(teamRequests, teamrequestmdl.Find(r, "TeamId", team.Id)...)
	}

	return teamRequests
}